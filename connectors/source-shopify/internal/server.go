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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"hash"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name        = "Shopify Source"
	defaultPort = 8080

	defaultSource              = "vanus-shopify-source" // TODO webhook id
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
)

var _ cdkgo.SourceConfigAccessor = &shopifySourceConfig{}

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

type shopifySourceConfig struct {
	cdkgo.SourceConfig `json:"_,inline" yaml:",inline"`
	Port               int    `json:"port" yaml:"port"`
	ClientSecret       string `json:"client_secret" yaml:"client_secret"`
}

func (c *shopifySourceConfig) GetSecret() cdkgo.SecretAccessor {
	return nil
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &shopifySourceConfig{}
}

var _ cdkgo.Source = &shopifySource{}

func NewShopifySource() cdkgo.Source {
	return &shopifySource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

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
	if !ok {
		return errors.New("invalid config")
	}

	if _cfg.Port == 0 {
		_cfg.Port = defaultPort
	}
	c.cfg = _cfg

	c.hm = hmac.New(sha256.New, []byte(c.cfg.ClientSecret))

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

func (c *shopifySource) Name() string {
	return name
}

func (c *shopifySource) Destroy() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err := c.ln.Close(); err != nil {
		log.Warning("close listener error", map[string]interface{}{
			log.KeyError: err,
		})
	}
	return nil
}

func (c *shopifySource) handleFastHTTP(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	h := hmac.New(sha256.New, []byte(c.cfg.ClientSecret))
	h.Write(body)
	hmacCalculated := base64.StdEncoding.EncodeToString(h.Sum(nil))
	hmacHeader := string(ctx.Request.Header.Peek(shopifyXHeaderHmac))

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
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Response.SetBody([]byte(fmt.Sprintf("failed to set data: %s", err.Error())))
		return
	}

	e.SetExtension(extendAttributesOrderID, string(ctx.Request.Header.Peek(shopifyXHeaderOrderID)))
	e.SetExtension(extendAttributesTopic, topic)
	e.SetExtension(extendAttributesWebhookID, string(ctx.Request.Header.Peek(shopifyXHeaderWebhookID)))
	e.SetExtension(extendAttributesShopDomain, string(ctx.Request.Header.Peek(shopifyXHeaderShopDomain)))
	e.SetExtension(extendAttributesAPIVersion, string(ctx.Request.Header.Peek(shopifyXHeaderAPIVersion)))

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
