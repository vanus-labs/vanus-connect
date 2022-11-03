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
	"github.com/cloudevents/sdk-go/v2/protocol"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
	cdkutil "github.com/linkall-labs/cdk-go/utils"
	"github.com/pkg/errors"
)

type ElasticsearchSink struct {
	config   *Config
	esClient *elasticsearch.Client

	timeout    time.Duration
	primaryKey PrimaryKey
	logger     log.Logger
}

func NewElasticsearchSink() connector.Sink {
	return &ElasticsearchSink{}
}

func (s *ElasticsearchSink) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}
	err := cfg.Validate()
	if err != nil {
		return err
	}
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
		s.logger.Info(context.TODO(), "es connect success", map[string]interface{}{
			"esInfo": resp,
		})
	}
	s.esClient = esClient
	return nil
}

func (s *ElasticsearchSink) Name() string {
	return "ElasticsearchSink"
}

func (s *ElasticsearchSink) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *ElasticsearchSink) Destroy() error {
	return nil
}

func (s *ElasticsearchSink) Receive(ctx context.Context, event ce.Event) protocol.Result {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
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
			s.logger.Warning(ctx, "documentID is empty", map[string]interface{}{
				"event":      event,
				"primaryKey": s.config.PrimaryKey,
			})
			return ce.NewHTTPResult(http.StatusBadRequest, "documentID is empty")
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
	resp, err := req.Do(ctx, s.esClient)
	if err != nil {
		s.logger.Warning(ctx, "es api do error", map[string]interface{}{
			log.KeyError: err,
			"event":      event,
		})
		return ce.NewHTTPResult(http.StatusInternalServerError, "write to es error %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.IsError() {
		respStr := resp.String()
		s.logger.Warning(ctx, "es api response error", map[string]interface{}{
			"resp": resp,
			"id":   event.ID(),
		})
		return ce.NewHTTPResult(http.StatusInternalServerError, "es api response error %s", respStr)
	}
	var res map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		s.logger.Warning(ctx, "parse response error", map[string]interface{}{
			log.KeyError: err,
			"id":         event.ID(),
		})
		return ce.NewHTTPResult(http.StatusInternalServerError, "parse response error %s", err.Error())
	}
	s.logger.Debug(ctx, "index api result ", map[string]interface{}{
		"result": res["result"],
		"id":     event.ID(),
	})
	return ce.ResultACK
}

func generateEsConfig(conf *Config) elasticsearch.Config {
	config := elasticsearch.Config{
		Addresses:     strings.Split(conf.Address, ","),
		Username:      conf.Username,
		Password:      conf.Password,
		RetryOnStatus: []int{429, 502, 503, 504},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return config
}
