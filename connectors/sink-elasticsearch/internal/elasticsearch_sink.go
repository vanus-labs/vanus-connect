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
	"net/http"
	"strings"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
)

type ElasticsearchSink struct {
	config   *Config
	esClient *elasticsearch.Client

	timeout    time.Duration
	primaryKey PrimaryKey
}

func (s *ElasticsearchSink) Initialize(_ context.Context, config cdkgo.ConfigAccessor) error {
	cfg := config.(*Config)
	s.config = cfg
	s.primaryKey = GetPrimaryKey(cfg.PrimaryKey)
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * 1000
	}
	s.timeout = time.Duration(cfg.Timeout) * time.Millisecond
	// init es client
	esClient, err := elasticsearch.NewClient(generateEsConfig(cfg))
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

func (s *ElasticsearchSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		result := s.writeEvent(ctx, event)
		if result == cdkgo.SuccessResult {
			return result
		}
	}
	return cdkgo.SuccessResult
}

func (s *ElasticsearchSink) writeEvent(ctx context.Context, event *ce.Event) cdkgo.Result {
	timeoutCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var (
		req          esapi.Request
		documentID   string
		documentType string
	)
	if s.primaryKey.Type() != None {
		documentID = s.primaryKey.Value(event)
	}
	if s.config.InsertMode == Upsert {
		if documentID == "" {
			log.Warning("documentID is empty", map[string]interface{}{
				"event":      event,
				"primaryKey": s.config.PrimaryKey,
			})
			return cdkgo.NewResult(http.StatusBadRequest, "documentID is empty")
		}
		var body bytes.Buffer
		body.WriteByte('{')
		body.WriteString(`"doc":`)
		body.Write(event.Data())
		body.WriteByte(',')
		body.WriteString(`"upsert":`)
		body.Write(event.Data())
		body.WriteByte('}')
		// https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-update.html
		req = esapi.UpdateRequest{
			Index:        s.config.IndexName,
			Body:         bytes.NewReader(body.Bytes()),
			DocumentID:   documentID,
			DocumentType: documentType,
		}
	} else {
		// https://www.elastic.co/guide/en/elasticsearch/reference/7.17/docs-index_.html
		req = esapi.IndexRequest{
			Index:        s.config.IndexName,
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

func (s *ElasticsearchSink) Name() string {
	return "ElasticsearchSink"
}

func (s *ElasticsearchSink) Destroy() error {
	return nil
}

func generateEsConfig(conf *Config) elasticsearch.Config {
	config := elasticsearch.Config{
		Addresses:     strings.Split(conf.Address, ","),
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
