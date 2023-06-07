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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	v2 "github.com/cloudevents/sdk-go/v2"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name                       = "HTTP Source"
	reqSource                  = "source"
	reqType                    = "type"
	reqID                      = "id"
	reqSubject                 = "subject"
	reqSchema                  = "dataschema"
	defaultSource              = "vanus-http-source"
	defaultType                = "naive-http-request"
	extendAttributesBodyIsJSON = "xvhttpbodyisjson"
	extendAttributesRemoteIP   = "xvhttpremoteip"
	extendAttributesRemoteAddr = "xvhttpremoteaddr"
)

type HTTPEvent struct {
	Path      string            `json:"path"`
	Method    string            `json:"method"`
	QueryArgs map[string]string `json:"query_args"`
	Headers   map[string]string `json:"headers"`
	Body      interface{}       `json:"body"`
}

func (he *HTTPEvent) toMap() map[string]interface{} {
	return map[string]interface{}{
		"path":       he.Path,
		"method":     he.Method,
		"query_args": he.QueryArgs,
		"headers":    he.Headers,
		"body":       he.Body,
	}
}

var _ cdkgo.Source = &httpSource{}

func NewHTTPSource() cdkgo.HTTPSource {
	return &httpSource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

type httpSource struct {
	cfg    *httpSourceConfig
	mutex  sync.Mutex
	ch     chan *cdkgo.Tuple
	logger zerolog.Logger
}

func (c *httpSource) Chan() <-chan *cdkgo.Tuple {
	return c.ch
}

func (c *httpSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	c.logger = log.FromContext(ctx)
	c.cfg = cfg.(*httpSourceConfig)
	return nil
}

func (c *httpSource) Name() string {
	return name
}

func (c *httpSource) Destroy() error {
	return nil
}

func (c *httpSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	he := &HTTPEvent{
		Path:      req.RequestURI,
		Method:    req.Method,
		QueryArgs: getQueryArgs(req),
		Headers:   getHeaders(req),
	}
	e := v2.NewEvent()
	mappingAttributes(req, he, &e)

	// try to convert request.Body to json
	m := map[string]interface{}{}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = json.Unmarshal(body, &m)
	if err == nil {
		he.Body = m
		e.SetExtension(extendAttributesBodyIsJSON, true)
	} else {
		he.Body = string(body)
		e.SetExtension(extendAttributesBodyIsJSON, false)
	}

	if err = e.SetData(v2.ApplicationJSON, he); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to set data: %s", err.Error())))
		return
	}
	c.logger.Debug().Interface("event", e).Msg("received a HTTP Request, ready to send")

	wg := sync.WaitGroup{}
	wg.Add(1)
	c.ch <- &cdkgo.Tuple{
		Event: &e,
		Success: func() {
			defer wg.Done()
			c.logger.Info().Str("event_id", e.ID()).Msg("send event to target success")
			w.WriteHeader(http.StatusOK)
		},
		Failed: func(err2 error) {
			defer wg.Done()
			c.logger.Warn().Interface("event_id", e.ID()).Err(err2).Msg("failed to send event to target")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(
				fmt.Sprintf("failed to send event to remote server: %s", err2.Error())))
		},
	}
	wg.Wait()
}

func getQueryArgs(req *http.Request) map[string]string {
	m := map[string]string{}
	values := req.URL.Query()
	for key, value := range values {
		m[key] = value[0]
	}
	return m
}

func getHeaders(req *http.Request) map[string]string {
	m := map[string]string{}
	headers := req.Header
	for key, header := range headers {
		if key == "Authorization" {
			continue
		}
		m[key] = header[0]
	}
	return m
}

func mappingAttributes(req *http.Request, he *HTTPEvent, e *v2.Event) {
	args := he.QueryArgs
	if v, ok := args[reqID]; ok && v != "" {
		e.SetID(v)
	} else {
		e.SetID(uuid.NewString())
	}
	if v, ok := args[reqSource]; ok && v != "" {
		e.SetSource(v)
	} else {
		e.SetSource(defaultSource)
	}
	if v, ok := args[reqType]; ok && v != "" {
		e.SetType(v)
	} else {
		e.SetType(defaultType)
	}
	if v, ok := args[reqSubject]; ok && v != "" {
		e.SetSubject(v)
	}
	if v, ok := args[reqSchema]; ok && v != "" {
		e.SetDataSchema(reqSchema)
	}
	//e.SetExtension(extendAttributesRemoteIP, ctx.RemoteIP().String())
	e.SetExtension(extendAttributesRemoteAddr, req.RemoteAddr)
}
