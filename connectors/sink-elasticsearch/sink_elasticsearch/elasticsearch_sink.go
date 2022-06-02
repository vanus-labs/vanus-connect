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

package sink_elasticsearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"io/ioutil"
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/go-logr/logr"
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
)

type ElasticsearchSink struct {
	config    elasticsearch.Config
	client    *elasticsearch.Client
	indexName string

	ceClient ce.Client
	logger   logr.Logger
	ctx      context.Context
}

func NewElasticsearchSink(ctx context.Context, ceClient ce.Client) connector.Sink {
	conf := getConfig()
	config := getEsConfig(conf)
	return &ElasticsearchSink{
		ctx:       ctx,
		config:    config,
		indexName: conf.IndexName,
		ceClient:  ceClient,
		logger:    log.FromContext(ctx),
	}
}

func (es *ElasticsearchSink) Start() error {
	ctx := context.Background()
	es.logger.Info("start es target")
	client, err := elasticsearch.NewClient(es.config)
	if err != nil {
		es.logger.Error(err, "create es client error")
		return errors.Wrap(err, "create es client error")
	}

	resp, err := client.Info()
	if err != nil {
		es.logger.Error(err, "client info api error")
		return errors.Wrap(err, "client info api error")
	}
	es.client = client
	if !resp.IsError() {
		info, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			es.logger.Error(err, "read es info error")
			return errors.Wrap(err, "read es info error")
		}
		es.logger.Info("es info is", "info", string(info))
	}
	return es.ceClient.StartReceiver(ctx, es.dispatch)
}

func (es *ElasticsearchSink) dispatch(ctx context.Context, event ce.Event) ce.Result {
	req := esapi.IndexRequest{
		Index: es.indexName,
		Body:  bytes.NewReader(event.Data()),
	}
	resp, err := req.Do(ctx, es.client)
	if err != nil {
		es.logger.Error(err, "es index request api do error")
		return ce.ResultNACK
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		es.logger.Error(err, "read es body error")
		return ce.ResultNACK
	}
	if resp.IsError() {
		es.logger.Error(err, "es api response is error",
			"statusCode", resp.StatusCode,
			"body", string(body),
		)
		return ce.ResultNACK
	}
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		es.logger.Error(err, "decode resp body error", "body", string(body))
		return ce.ResultNACK
	}
	es.logger.Info("index api result ", "result", res["result"])
	return ce.ResultACK
}

func getEsConfig(conf Config) elasticsearch.Config {
	config := elasticsearch.Config{
		Addresses:     conf.Addresses,
		Username:      conf.Username,
		Password:      conf.Password,
		RetryOnStatus: []int{429, 502, 503, 504},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: conf.SkipVerify,
			},
		},
	}
	return config
}
