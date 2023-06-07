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
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/valyala/fasthttp"

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
		QueryArgs: getQueryArgsNew(req),
		Headers:   getHeadersNew(req),
	}
	e := v2.NewEvent()
	mappingAttributesNew(req, he, &e)

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
			w.WriteHeader(http.StatusOK)
			wg.Done()
		},
		Failed: func(err2 error) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(
				fmt.Sprintf("failed to send event to remote server: %s", err2.Error())))
			wg.Done()
		},
	}
	wg.Wait()
}

func getQueryArgsNew(req *http.Request) map[string]string {
	m := map[string]string{}
	values := req.URL.Query()
	for key, value := range values {
		m[key] = value[0]
	}
	return m
}

func getHeadersNew(req *http.Request) map[string]string {
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

func mappingAttributesNew(req *http.Request, he *HTTPEvent, e *v2.Event) {
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

func (c *httpSource) handleFastHTTP(ctx *fasthttp.RequestCtx) {

	he := &HTTPEvent{
		Path:      string(ctx.Path()),
		Method:    string(ctx.Method()),
		QueryArgs: getQueryArgs(ctx),
		Headers:   getHeaders(ctx),
	}

	e := v2.NewEvent()
	mappingAttributes(ctx, &e)

	// try to convert request.Body to json
	m := map[string]interface{}{}
	err := json.Unmarshal(ctx.Request.Body(), &m)
	if err == nil {
		he.Body = m
		e.SetExtension(extendAttributesBodyIsJSON, true)
	} else {
		he.Body = string(ctx.Request.Body())
		e.SetExtension(extendAttributesBodyIsJSON, false)
	}

	if err = e.SetData(v2.ApplicationJSON, he.toMap()); err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBody([]byte(fmt.Sprintf("failed to set data: %s", err.Error())))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	c.ch <- &cdkgo.Tuple{
		Event: &e,
		Success: func() {
			ctx.Response.SetStatusCode(http.StatusOK)
			wg.Done()
		},
		Failed: func(err2 error) {
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.Response.SetBody([]byte(
				fmt.Sprintf("failed to send event to remote server: %s", err2.Error())))
			wg.Done()
		},
	}
	wg.Wait()
}

func mappingAttributes(ctx *fasthttp.RequestCtx, e *v2.Event) {
	args := ctx.QueryArgs()
	if args.Has(reqID) && len(args.Peek(reqID)) > 0 {
		e.SetID(string(args.Peek(reqID)))
	} else {
		e.SetID(uuid.NewString())
	}

	if args.Has(reqSource) && len(args.Peek(reqSource)) > 0 {
		e.SetSource(string(args.Peek(reqSource)))
	} else {
		e.SetSource(defaultSource)
	}

	if args.Has(reqType) && len(args.Peek(reqType)) > 0 {
		e.SetType(string(args.Peek(reqType)))
	} else {
		e.SetType(defaultType)
	}

	if args.Has(reqSubject) && len(args.Peek(reqSubject)) > 0 {
		e.SetSubject(string(args.Peek(reqSubject)))
	}

	if args.Has(reqSchema) && len(args.Peek(reqSchema)) > 0 {
		e.SetDataSchema(string(args.Peek(reqSchema)))
	}

	e.SetExtension(extendAttributesRemoteIP, ctx.RemoteIP().String())
	e.SetExtension(extendAttributesRemoteAddr, ctx.RemoteAddr().String())
}

func getQueryArgs(ctx *fasthttp.RequestCtx) map[string]string {
	m := map[string]string{}
	args := strings.Split(ctx.QueryArgs().String(), "&")
	for _, arg := range args {
		kv := strings.Split(arg, "=")
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return m
}

func getHeaders(ctx *fasthttp.RequestCtx) map[string]string {
	m := map[string]string{}
	r := bufio.NewReader(bytes.NewReader(ctx.Request.Header.Header()))
	for {
		l, isPrefix, err := r.ReadLine()
		if err != nil {
			break
		}
		var _l []byte
		for isPrefix {
			_l, isPrefix, err = r.ReadLine()
			if err != nil {
				break
			}
			l = append(l, _l...)
		}
		str := string(l)
		idx := strings.Index(str, ":")
		if idx == -1 {
			// ignore something like POST /webhook?source=123&id=1234sda&type=xxxxx&subject=12eqsd&asdax=asdasd HTTP/1.1
			continue
		}
		if idx+2 < len(str)-1 {
			m[str[0:idx]] = str[idx+2:]
		}
	}

	return m
}
