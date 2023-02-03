// Copyright 2022 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
)

type elasticsearchSink struct {
	count    int64
	config   *esConfig
	esClient *es.Client

	timeout time.Duration
	buf     *bytes.Buffer
	action  action
}

func Sink() cdkgo.Sink {
	return &elasticsearchSink{}
}

func (s *elasticsearchSink) Initialize(_ context.Context, config cdkgo.ConfigAccessor) error {
	cfg := config.(*esConfig)
	s.config = cfg
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * 1000
	}
	if cfg.BufferBytes == 0 {
		// default 5MB
		cfg.BufferBytes = 5 * 1024 * 1024
	}
	if cfg.InsertMode == Upsert {
		s.action = actionUpdate
	} else {
		s.action = actionIndex
	}
	s.timeout = time.Duration(cfg.Timeout) * time.Millisecond
	s.buf = bytes.NewBuffer(make([]byte, 0, cfg.BufferBytes))
	// init es client
	esClient, err := es.NewClient(generateEsConfig(cfg))
	if err != nil {
		return errors.Wrap(err, "new es client error")
	}
	resp, err := esClient.Info()
	if err != nil {
		return errors.Wrap(err, "es info api error")
	}
	if !resp.IsError() {
		log.Info("es connect success", map[string]interface{}{
			"esInfo": resp,
		})
	}
	s.esClient = esClient
	return nil
}

func (s *elasticsearchSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	if len(events) == 0 {
		return cdkgo.SuccessResult
	}
	atomic.AddInt64(&s.count, int64(len(events)))
	log.Info("receive event count", map[string]interface{}{
		"total": s.count,
	})
	for _, event := range events {
		err := s.appendEvent(event)
		if err != nil {
			s.cleanBuffer()
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
	}
	defer s.cleanBuffer()
	timeoutCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	req := esapi.BulkRequest{
		Body: s.buf,
	}
	res, err := req.Do(timeoutCtx, s.esClient)
	if err != nil {
		log.Warning("es bulk do error", map[string]interface{}{
			log.KeyError: err,
			"total":      len(events),
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "write to es error")
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.IsError() {
		log.Warning("es bulk response error", map[string]interface{}{
			"total":    len(events),
			"response": res.String(),
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "es response error")
	}
	var blk esutil.BulkIndexerResponse
	if err = json.NewDecoder(res.Body).Decode(&blk); err != nil {
		log.Warning("parse response error", map[string]interface{}{
			log.KeyError: err,
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "parse response error")
	}
	if !blk.HasErrors {
		for i, blkItem := range blk.Items {
			for k, v := range blkItem {
				if v.Error.Type != "" || v.Status > 201 {
					log.Warning("event write to es failed", map[string]interface{}{
						"index":       i + 1,
						"id":          events[i].ID(),
						"action":      k,
						"errorType":   v.Error.Type,
						"errorReason": v.Error.Reason,
					})
				}
			}
		}
	}
	return cdkgo.SuccessResult
}

func (s *elasticsearchSink) cleanBuffer() {
	if s.buf.Cap() > s.config.BufferBytes {
		s.buf = bytes.NewBuffer(make([]byte, 0, s.config.BufferBytes))
	} else {
		s.buf.Reset()
	}
}

func (s *elasticsearchSink) appendEvent(event *ce.Event) error {
	extensions := event.Extensions()
	index, err := s.getIndexName(extensions)
	if err != nil {
		return err
	}
	actionName, err := s.getAction(extensions)
	if err != nil {
		return err
	}
	documentId, err := s.getDocumentId(extensions)
	if err != nil {
		return err
	}
	if documentId == "" {
		if actionName == actionUpdate || actionName == actionDelete {
			return errors.Errorf("action is %s but documentId is empty", actionName)
		}
	}
	// https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html
	buf := s.buf
	buf.WriteRune('{')
	buf.WriteString(strconv.Quote(string(actionName)))
	buf.WriteRune(':')
	buf.WriteRune('{')
	buf.WriteString(`"_index":`)
	buf.WriteString(strconv.Quote(index))
	if documentId != "" {
		buf.WriteRune(',')
		buf.WriteString(`"_id":`)
		buf.WriteString(strconv.Quote(documentId))
	}
	buf.WriteRune('}')
	buf.WriteRune('}')
	buf.WriteRune('\n')
	if actionName == actionIndex {
		json.Compact(buf, event.Data())
	} else if actionName == actionUpdate {
		buf.WriteRune('{')
		buf.WriteString(`"doc":`)
		json.Compact(buf, event.Data())
		buf.WriteRune(',')
		buf.WriteString(`"doc_as_upsert":true`)
		buf.WriteRune('}')
	}
	buf.WriteRune('\n')
	return nil
}

func (s *elasticsearchSink) writeEvent(ctx context.Context, event *ce.Event) cdkgo.Result {
	timeoutCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var (
		req          esapi.Request
		documentID   string
		documentType string
	)
	if s.config.InsertMode == Upsert {
		var body bytes.Buffer
		body.WriteByte('{')
		body.WriteString(`"doc":`)
		body.Write(event.Data())
		body.WriteByte(',')
		body.WriteString(`"doc_as_upsert": true`)
		body.WriteByte('}')
		// https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-update.html
		req = esapi.UpdateRequest{
			Index:        s.config.Secret.IndexName,
			Body:         bytes.NewReader(body.Bytes()),
			DocumentID:   documentID,
			DocumentType: documentType,
		}
	} else {
		// https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-index_.html
		req = esapi.IndexRequest{
			Index:        s.config.Secret.IndexName,
			Body:         bytes.NewReader(event.Data()),
			DocumentID:   documentID,
			DocumentType: documentType,
		}
	}
	resp, err := req.Do(timeoutCtx, s.esClient)
	if err != nil {
		log.Warning("es api do error", map[string]interface{}{
			log.KeyError: err,
			"event":      event,
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "write to es error")
	}
	defer resp.Body.Close()
	if resp.IsError() {
		log.Warning("es api response error", map[string]interface{}{
			"resp": resp.String(),
			"id":   event.ID(),
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "es api response error ")
	}
	var res map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Warning("parse response error", map[string]interface{}{
			log.KeyError: err,
			"id":         event.ID(),
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "parse response error")
	}
	log.Debug("index api result ", map[string]interface{}{
		"result": res["result"],
		"id":     event.ID(),
	})
	return cdkgo.SuccessResult
}

func (s *elasticsearchSink) Name() string {
	return "ElasticsearchSink"
}

func (s *elasticsearchSink) Destroy() error {
	return nil
}

func generateEsConfig(conf *esConfig) es.Config {
	config := es.Config{
		Addresses:     strings.Split(conf.Secret.Address, ","),
		Username:      conf.Secret.Username,
		Password:      conf.Secret.Password,
		RetryOnStatus: []int{429, 502, 503, 504},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return config
}
