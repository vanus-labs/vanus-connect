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
	"fmt"
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/pkg/errors"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Source = &shopifySource{}

func NewSource() cdkgo.Source {
	return &shopifySource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

type shopifySource struct {
	events        chan *cdkgo.Tuple
	config        *shopifyConfig
	client        *goshopify.Client
	syncBeginTime time.Time
	syncInternal  time.Duration
}

func (s *shopifySource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*shopifyConfig)
	t, err := time.Parse("2006-01-02", s.config.SyncBeginDate)
	if err != nil {
		return err
	}
	s.syncBeginTime = t
	if s.config.SyncInternalHour <= 0 || s.config.SyncInternalHour > 24 {
		s.config.SyncInternalHour = 1
	}
	s.syncInternal = time.Duration(s.config.SyncInternalHour) * time.Hour
	s.client = goshopify.NewClient(goshopify.App{}, s.config.ShopName, s.config.ApiAccessToken, goshopify.WithVersion("2023-04"))
	go s.start(ctx)
	return nil
}

func (s *shopifySource) Name() string {
	return "ShopifyAppSource"
}

func (s *shopifySource) Destroy() error {
	return nil
}

func (s *shopifySource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *shopifySource) start(ctx context.Context) {
	err := s.sync(ctx)
	if err != nil {
		log.Info("sync failed", map[string]interface{}{
			log.KeyError: err,
		})
	}
	t := time.NewTicker(s.syncInternal)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return
	case <-t.C:
		err = s.sync(ctx)
		if err != nil {
			log.Info("sync failed", map[string]interface{}{
				log.KeyError: err,
			})
		}
	}
}

func (s *shopifySource) sync(ctx context.Context) error {
	var err error
	err = s.syncOrders(ctx)
	if err != nil {
		log.Error("sync order failed", map[string]interface{}{
			log.KeyError: err,
		})
	}
	err = s.syncProducts(ctx)
	if err != nil {
		log.Error("sync product failed", map[string]interface{}{
			log.KeyError: err,
		})
	}
	return nil
}

func (s *shopifySource) syncOrders(ctx context.Context) error {
	apiType := OrderApi
	begin, err := s.getSyncBeginTime(ctx, apiType)
	if err != nil {
		return errors.Wrapf(err, "get %v sync begin timer error", apiType)
	}
	end := time.Now().UTC()
	listOptions := &goshopify.ListOptions{
		CreatedAtMin: begin,
		CreatedAtMax: end,
	}
	c := 0
	for {
		list, pageOptions, err := s.client.Order.ListWithPagination(listOptions)
		if err != nil {
			return errors.Wrapf(err, "list %v failed", apiType)
		}
		if len(list) == 0 {
			break
		}
		c += len(list)
		s.orderEvent(list)
		if pageOptions == nil || pageOptions.NextPageOptions == nil {
			break
		}
		listOptions = pageOptions.NextPageOptions
	}
	log.Info(fmt.Sprintf("sync %v success", apiType), map[string]interface{}{
		"count": c,
	})
	return s.setSyncTime(ctx, apiType, end)
}

func (s *shopifySource) syncProducts(ctx context.Context) error {
	apiType := ProductApi
	begin, err := s.getSyncBeginTime(ctx, apiType)
	if err != nil {
		return errors.Wrapf(err, "get %v sync begin timer error", apiType)
	}
	end := time.Now().UTC()
	listOptions := &goshopify.ListOptions{
		CreatedAtMin: begin,
		CreatedAtMax: end,
	}
	c := 0
	for {
		list, pageOptions, err := s.client.Product.ListWithPagination(listOptions)
		if err != nil {
			return errors.Wrapf(err, "list %v failed", apiType)
		}
		if len(list) == 0 {
			break
		}
		c += len(list)
		s.productEvent(list)
		if pageOptions == nil || pageOptions.NextPageOptions == nil {
			break
		}
		listOptions = pageOptions.NextPageOptions
	}
	log.Info(fmt.Sprintf("sync %v success", apiType), map[string]interface{}{
		"count": c,
	})
	return s.setSyncTime(ctx, apiType, end)
}
