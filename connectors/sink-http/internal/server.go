// Copyright 2023 Linkall Inc.
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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync/atomic"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
)

var _ cdkgo.Sink = &httpSink{}

func NewHTTPSink() cdkgo.Sink {
	return &httpSink{}
}

type httpSink struct {
	count   int64
	client  *http.Client
	url     *url.URL
	method  string
	headers map[string]string
	Auth    Auth
}

func (s *httpSink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	config := cfg.(*httpConfig)
	s.url, _ = url.Parse(config.Target)
	s.headers = config.Headers
	s.method = config.Method
	if s.method == "" {
		s.method = "POST"
	}
	s.Auth = config.Auth
	s.client = http.DefaultClient
	return nil
}

func (s *httpSink) Name() string {
	return "HTTP Sink"
}

func (s *httpSink) Destroy() error {
	return nil
}

func (s *httpSink) Arrived(_ context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		atomic.AddInt64(&s.count, 1)
		log.Info("receive a new event", map[string]interface{}{
			"in_total": atomic.LoadInt64(&s.count),
		})
		r := s.sendEvent(event)
		if r != cdkgo.SuccessResult {
			return r
		}
		log.Info("send event success", map[string]interface{}{
			"id": event.ID(),
		})
	}
	return cdkgo.SuccessResult
}

func (s *httpSink) sendEvent(event *ce.Event) cdkgo.Result {
	m := &Request{}
	if err := event.DataAs(m); err != nil {
		return cdkgo.NewResult(http.StatusBadRequest, fmt.Sprintf("event data invalid %s", err.Error()))
	}
	u := *s.url
	if m.Query != "" {
		kv, err := url.ParseQuery(m.Query)
		if err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, fmt.Sprintf("event data query invalid %s", err.Error()))
		}
		values := u.Query()
		for k := range values {
			kv.Set(k, values.Get(k))
		}
		u.RawQuery = kv.Encode()
	}
	if m.Path != "" {
		u.Path = path.Join(u.Path, m.Path)
	}
	method := m.Method
	if method == "" {
		method = s.method
	}

	req, err := http.NewRequest(method, u.String(), strings.NewReader(m.Body))
	if err != nil {
		return cdkgo.NewResult(http.StatusInternalServerError, fmt.Sprintf("new http request error %s", err.Error()))
	}

	if s.Auth.Username != "" || s.Auth.Password != "" {
		req.SetBasicAuth(s.Auth.Username, s.Auth.Password)
	}
	// common default header
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}
	// request headers
	for k, v := range m.Headers {
		req.Header.Set(k, v)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return cdkgo.NewResult(http.StatusInternalServerError, fmt.Sprintf("sned http request error %s", err.Error()))
	}
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return cdkgo.NewResult(http.StatusInternalServerError, fmt.Sprintf("read http response error %s", err.Error()))
	}
	if res.StatusCode >= 400 {
		return cdkgo.NewResult(connector.Code(res.StatusCode), fmt.Sprintf("http response code %d resp %s", res.StatusCode, string(resp)))
	}
	return cdkgo.SuccessResult
}
