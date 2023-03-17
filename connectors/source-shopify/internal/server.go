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
<<<<<<< HEAD
<<<<<<< HEAD
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
<<<<<<< HEAD
	"hash"
=======
	"bufio"
	"bytes"
=======
>>>>>>> 2f93b62 (feat: add shopify source)
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
<<<<<<< HEAD
>>>>>>> d269259 (feat: add shopify source)
=======
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
=======
>>>>>>> 7a9eb27 (update readme)
	"hash"
>>>>>>> 2f93b62 (feat: add shopify source)
	"net"
	"net/http"
	"strings"
	"sync"
<<<<<<< HEAD
<<<<<<< HEAD
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
<<<<<<< HEAD
=======

	"github.com/google/uuid"

	v2 "github.com/cloudevents/sdk-go/v2"
>>>>>>> d269259 (feat: add shopify source)
=======
	"time"

>>>>>>> 2f93b62 (feat: add shopify source)
=======
>>>>>>> 7a9eb27 (update readme)
	"github.com/valyala/fasthttp"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
<<<<<<< HEAD
<<<<<<< HEAD
	name        = "Shopify Source"
	defaultPort = 8080

	defaultSource              = "vanus-shopify-source"
=======
	name        = "Shopify Source"
	defaultPort = 8080

<<<<<<< HEAD
	defaultSource              = "vanus-shopify-source" // TODO webhook id
>>>>>>> 2f93b62 (feat: add shopify source)
=======
	defaultSource              = "vanus-shopify-source"
>>>>>>> 7a9eb27 (update readme)
	extendAttributesOrderID    = "xvshopifyorderid"
	extendAttributesTopic      = "xvshopifytopic"
	extendAttributesWebhookID  = "xvshopifywebhookid"
	extendAttributesShopDomain = "xvshopifydomain"
	extendAttributesAPIVersion = "xvshopifyapiversion"

	shopifyXHeaderAPIVersion = "X-Shopify-Api-Version"
	shopifyXHeaderHmac       = "X-Shopify-Hmac-Sha256"
	shopifyXHeaderOrderID    = "X-Shopify-Order-Id"
	shopifyXHeaderShopDomain = "X-Shopify-Shop-Domain"
	shopifyXHeaderTopic      = "X-Shopify-Topic"
	shopifyXHeaderWebhookID  = "X-Shopify-Webhook-Id"
<<<<<<< HEAD
)

var _ cdkgo.SourceConfigAccessor = &shopifySourceConfig{}

type shopifySourceConfig struct {
	cdkgo.SourceConfig `json:"_,inline" yaml:",inline"`
	Port               int    `json:"port" yaml:"port"`
	ClientSecret       string `json:"client_secret" yaml:"client_secret"`
}

func (c *shopifySourceConfig) GetSecret() cdkgo.SecretAccessor {
=======
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
=======
>>>>>>> 2f93b62 (feat: add shopify source)
)

var _ cdkgo.SourceConfigAccessor = &shopifySourceConfig{}

type shopifySourceConfig struct {
	cdkgo.SourceConfig `json:"_,inline" yaml:",inline"`
	Port               int    `json:"port" yaml:"port"`
	ClientSecret       string `json:"client_secret" yaml:"client_secret"`
}

<<<<<<< HEAD
func (c *httpSourceConfig) GetSecret() cdkgo.SecretAccessor {
>>>>>>> d269259 (feat: add shopify source)
=======
func (c *shopifySourceConfig) GetSecret() cdkgo.SecretAccessor {
>>>>>>> 2f93b62 (feat: add shopify source)
	return nil
}

func NewConfig() cdkgo.SourceConfigAccessor {
<<<<<<< HEAD
<<<<<<< HEAD
	return &shopifySourceConfig{}
}

var _ cdkgo.Source = &shopifySource{}

func NewShopifySource() cdkgo.Source {
	return &shopifySource{
=======
	return &httpSourceConfig{}
=======
	return &shopifySourceConfig{}
>>>>>>> 2f93b62 (feat: add shopify source)
}

var _ cdkgo.Source = &shopifySource{}

<<<<<<< HEAD
func NewHTTPSource() cdkgo.Source {
	return &httpSource{
>>>>>>> d269259 (feat: add shopify source)
=======
func NewShopifySource() cdkgo.Source {
	return &shopifySource{
>>>>>>> 2f93b62 (feat: add shopify source)
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

<<<<<<< HEAD
<<<<<<< HEAD
type shopifySource struct {
	cfg   *shopifySourceConfig
	mutex sync.Mutex
	ch    chan *cdkgo.Tuple
	ln    net.Listener
	hm    hash.Hash
}

func (c *shopifySource) Chan() <-chan *cdkgo.Tuple {
	return c.ch
}

func (c *shopifySource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*shopifySourceConfig)
=======
type httpSource struct {
	cfg   *httpSourceConfig
=======
type shopifySource struct {
	cfg   *shopifySourceConfig
>>>>>>> 2f93b62 (feat: add shopify source)
	mutex sync.Mutex
	ch    chan *cdkgo.Tuple
	ln    net.Listener
	hm    hash.Hash
}

func (c *shopifySource) Chan() <-chan *cdkgo.Tuple {
	return c.ch
}

<<<<<<< HEAD
func (c *httpSource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*httpSourceConfig)
>>>>>>> d269259 (feat: add shopify source)
=======
func (c *shopifySource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*shopifySourceConfig)
>>>>>>> 2f93b62 (feat: add shopify source)
	if !ok {
		return errors.New("invalid config")
	}

	if _cfg.Port == 0 {
		_cfg.Port = defaultPort
	}
	c.cfg = _cfg

<<<<<<< HEAD
<<<<<<< HEAD
	c.hm = hmac.New(sha256.New, []byte(c.cfg.ClientSecret))

=======
>>>>>>> d269259 (feat: add shopify source)
=======
	c.hm = hmac.New(sha256.New, []byte(c.cfg.ClientSecret))

>>>>>>> 2f93b62 (feat: add shopify source)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.cfg.Port))
	if err != nil {
		return err
	}
	c.ln = ln
	go func() {
<<<<<<< HEAD
<<<<<<< HEAD
		log.Info("Shopify source is ready to serving", map[string]interface{}{
=======
		log.Info("HTTP source is ready to serving", map[string]interface{}{
>>>>>>> d269259 (feat: add shopify source)
=======
		log.Info("Shopify source is ready to serving", map[string]interface{}{
>>>>>>> 7a9eb27 (update readme)
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

<<<<<<< HEAD
<<<<<<< HEAD
func (c *shopifySource) Name() string {
	return name
}

func (c *shopifySource) Destroy() error {
=======
func (c *httpSource) Name() string {
	return name
}

func (c *httpSource) Destroy() error {
>>>>>>> d269259 (feat: add shopify source)
=======
func (c *shopifySource) Name() string {
	return name
}

func (c *shopifySource) Destroy() error {
>>>>>>> 2f93b62 (feat: add shopify source)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err := c.ln.Close(); err != nil {
		log.Warning("close listener error", map[string]interface{}{
			log.KeyError: err,
		})
	}
	return nil
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 2f93b62 (feat: add shopify source)
func (c *shopifySource) handleFastHTTP(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	h := hmac.New(sha256.New, []byte(c.cfg.ClientSecret))
	h.Write(body)
	hmacCalculated := base64.StdEncoding.EncodeToString(h.Sum(nil))
	hmacHeader := string(ctx.Request.Header.Peek(shopifyXHeaderHmac))
<<<<<<< HEAD

	// validate signature
	if !strings.EqualFold(hmacCalculated, hmacHeader) {
		ctx.Response.SetStatusCode(http.StatusUnauthorized)
		return
	}

	topic := string(ctx.Request.Header.Peek(shopifyXHeaderTopic))
	e := v2.NewEvent()
	e.SetID(uuid.NewString())
	e.SetSource(defaultSource)
	e.SetType(topic)
	e.SetTime(time.Now())
	m := map[string]interface{}{}
	err := json.Unmarshal(ctx.Request.Body(), &m)

	if err = e.SetData(v2.ApplicationJSON, m); err != nil {
=======
func (c *httpSource) handleFastHTTP(ctx *fasthttp.RequestCtx) {
=======
>>>>>>> 2f93b62 (feat: add shopify source)

	// validate signature
	if !strings.EqualFold(hmacCalculated, hmacHeader) {
		ctx.Response.SetStatusCode(http.StatusUnauthorized)
		return
	}

	topic := string(ctx.Request.Header.Peek(shopifyXHeaderTopic))
	e := v2.NewEvent()
	e.SetID(uuid.NewString())
	e.SetSource(defaultSource)
	e.SetType(topic)
	e.SetTime(time.Now())
	m := map[string]interface{}{}
	err := json.Unmarshal(ctx.Request.Body(), &m)

<<<<<<< HEAD
	if err = e.SetData(v2.ApplicationJSON, he.toMap()); err != nil {
>>>>>>> d269259 (feat: add shopify source)
=======
	if err = e.SetData(v2.ApplicationJSON, m); err != nil {
>>>>>>> 2f93b62 (feat: add shopify source)
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBody([]byte(fmt.Sprintf("failed to set data: %s", err.Error())))
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 2f93b62 (feat: add shopify source)
	e.SetExtension(extendAttributesOrderID, string(ctx.Request.Header.Peek(shopifyXHeaderOrderID)))
	e.SetExtension(extendAttributesTopic, topic)
	e.SetExtension(extendAttributesWebhookID, string(ctx.Request.Header.Peek(shopifyXHeaderWebhookID)))
	e.SetExtension(extendAttributesShopDomain, string(ctx.Request.Header.Peek(shopifyXHeaderShopDomain)))
	e.SetExtension(extendAttributesAPIVersion, string(ctx.Request.Header.Peek(shopifyXHeaderAPIVersion)))

<<<<<<< HEAD
=======
>>>>>>> d269259 (feat: add shopify source)
=======
>>>>>>> 2f93b62 (feat: add shopify source)
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
		Failed: func(err2 error) {
			ctx.Response.SetStatusCode(http.StatusInternalServerError)
			ctx.Response.SetBody([]byte(
				fmt.Sprintf("failed to send event to remote server: %s", err2.Error())))
			wg.Done()
		},
	}
	wg.Wait()
}
<<<<<<< HEAD
<<<<<<< HEAD
=======

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
>>>>>>> d269259 (feat: add shopify source)
=======
>>>>>>> 2f93b62 (feat: add shopify source)
