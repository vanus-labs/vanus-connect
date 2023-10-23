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
	"sync/atomic"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name = "Shopify Source"

	defaultSource               = "vanus-shopify-source"
	extendAttributesTopic       = "xvshopifytopic"
	extendAttributesWebhookID   = "xvshopifywebhookid"
	extendAttributesShopDomain  = "xvshopifydomain"
	extendAttributesAPIVersion  = "xvshopifyapiversion"
	extendAttributesTriggeredAt = "xvshopifytriggeredat"

	headerAPIVersion  = "X-Shopify-Api-Version"
	headerHmac        = "X-Shopify-Hmac-Sha256"
	headerShopDomain  = "X-Shopify-Shop-Domain"
	headerTopic       = "X-Shopify-Topic"
	headerWebhookID   = "X-Shopify-Webhook-Id"
	headerTriggeredAt = "X-Shopify-Triggered-At"
)

var _ cdkgo.Source = &shopifySource{}

func NewShopifySource() cdkgo.HTTPSource {
	return &shopifySource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

type shopifySource struct {
	ch     chan *cdkgo.Tuple
	logger zerolog.Logger
	app    goshopify.App
	count  int64
}

func (c *shopifySource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !c.app.VerifyWebhookRequest(r) {
		c.logger.Info().Msg("hmac invalid")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("hmac invalid"))
		return
	}
	atomic.AddInt64(&c.count, 1)
	log.Info().Int64("total", atomic.LoadInt64(&c.count)).Msg("receive a new event")
	topic := r.Header.Get(headerTopic)
	e := ce.NewEvent()
	e.SetID(uuid.NewString())
	e.SetSource(defaultSource)
	e.SetType(topic)
	e.SetTime(time.Now())
	body, _ := io.ReadAll(r.Body)
	var m map[string]interface{}
	json.Unmarshal(body, &m)
	e.SetData(ce.ApplicationJSON, m)

	e.SetExtension(extendAttributesWebhookID, r.Header.Get(headerWebhookID))
	e.SetExtension(extendAttributesShopDomain, r.Header.Get(headerShopDomain))
	e.SetExtension(extendAttributesAPIVersion, r.Header.Get(headerAPIVersion))
	e.SetExtension(extendAttributesTriggeredAt, r.Header.Get(headerTriggeredAt))
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

func (c *shopifySource) Chan() <-chan *cdkgo.Tuple {
	return c.ch
}

func (c *shopifySource) Initialize(ctx context.Context, config cdkgo.ConfigAccessor) error {
	c.logger = log.FromContext(ctx)
	cfg := config.(*shopifySourceConfig)
	c.app = goshopify.App{ApiSecret: cfg.ClientSecret}
	return nil
}

func (c *shopifySource) Name() string {
	return name
}

func (c *shopifySource) Destroy() error {
	return nil
}
