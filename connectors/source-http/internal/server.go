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
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net"
	"net/http"
	"strings"
	"sync"

	v2 "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/valyala/fasthttp"
)

const (
	name                       = "HTTP Source"
	defaultPort                = 8080
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

var _ cdkgo.SourceConfigAccessor = &httpSourceConfig{}

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

type httpSourceConfig struct {
	cdkgo.SourceConfig `json:"_,inline" yaml:",inline"`
	Port               int `json:"port" yaml:"port"`
}

func (c *httpSourceConfig) GetSecret() cdkgo.SecretAccessor {
	return nil
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &httpSourceConfig{}
}

var _ cdkgo.Source = &httpSource{}

func NewHTTPSource() cdkgo.Source {
	return &httpSource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

type httpSource struct {
	cfg   *httpSourceConfig
	mutex sync.Mutex
	ch    chan *cdkgo.Tuple
	ln    net.Listener
}

func (c *httpSource) Chan() <-chan *cdkgo.Tuple {
	return c.ch
}

func (c *httpSource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*httpSourceConfig)
	if !ok {
		return errors.New("invalid config")
	}

	if _cfg.Port == 0 {
		_cfg.Port = defaultPort
	}
	c.cfg = _cfg

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.cfg.Port))
	if err != nil {
		return err
	}
	c.ln = ln
	go func() {
		log.Info("HTTP source is ready to serving", map[string]interface{}{
			"listen": c.cfg.Port,
		})
		if err := fasthttp.Serve(ln, c.handleFastHTTP); err != nil {
			log.Error("failed to start http server", map[string]interface{}{
				log.KeyError: err,
			})
			panic(err)
		}
	}()
	return nil
}

func (c *httpSource) Name() string {
	return name
}

func (c *httpSource) Destroy() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err := c.ln.Close(); err != nil {
		log.Warning("close listener error", map[string]interface{}{
			log.KeyError: err,
		})
	}
	return nil
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
	d, _ := e.MarshalJSON()
	println(string(d))
	log.Debug("received a HTTP Request, ready to send", map[string]interface{}{
		"event": e.String(),
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	c.ch <- &cdkgo.Tuple{
		Event: &e,
		Success: func() {
			ctx.Response.SetStatusCode(http.StatusOK)
			wg.Done()
		},
		Failed: func() {
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.Response.SetBody([]byte("failed to send event to remote server"))
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
