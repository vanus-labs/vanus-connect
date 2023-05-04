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

var syncApiArr = []apiType{OrderApi, ProductApi}

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
	err = s.initSyncTime(ctx)
	if err != nil {
		return err
	}
	if s.config.SyncIntervalHour <= 0 || s.config.SyncIntervalHour > 24 {
		s.config.SyncIntervalHour = 1
	}
	s.syncInternal = time.Duration(s.config.SyncIntervalHour) * time.Hour
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

func (s *shopifySource) initSyncTime(ctx context.Context) error {
	syncBeginDate, err := getSyncBeginDate(ctx)
	if err != nil {
		return errors.Wrap(err, "get sync begin date error")
	}
	if syncBeginDate == s.config.SyncBeginDate {
		log.Info("sync begin date no change", map[string]interface{}{
			"sync_begin_date": s.config.SyncBeginDate,
		})
		return nil
	}
	for _, t := range syncApiArr {
		err = setSyncTime(ctx, t, s.syncBeginTime)
		if err != nil {
			return errors.Wrapf(err, "api %v set sync time error", t)
		}
	}
	err = setSyncBeginDate(ctx, s.config.SyncBeginDate)
	if err != nil {
		return errors.Wrapf(err, "set sync begin date error")
	}
	log.Info("init sync time success", map[string]interface{}{
		"sync_begin_date": s.config.SyncBeginDate,
	})
	return nil
}

func (s *shopifySource) start(ctx context.Context) {
	s.sync(ctx)
	t := time.NewTicker(s.syncInternal)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.sync(ctx)
		}
	}
}

func (s *shopifySource) sync(ctx context.Context) {
	for _, apiType := range syncApiArr {
		begin, err := getSyncTime(ctx, apiType)
		if err != nil {
			log.Warning(fmt.Sprintf("get %v sync begin timer error", apiType), map[string]interface{}{
				log.KeyError: err,
			})
			continue
		}
		end := time.Now().UTC()
		log.Info(fmt.Sprintf("sync %v data begin", apiType), map[string]interface{}{
			"begin_time": begin,
			"end_time":   end,
		})
		var c int
		switch apiType {
		case OrderApi:
			c, err = s.syncOrders(ctx, begin, end)
		case ProductApi:
			c, err = s.syncProducts(ctx, begin, end)
		}
		if err != nil {
			log.Warning(fmt.Sprintf("sync %v data error", apiType), map[string]interface{}{
				log.KeyError: err,
				"count":      c,
			})
			continue
		}
		log.Info(fmt.Sprintf("sync %v data success", apiType), map[string]interface{}{
			"count": c,
		})
		err = setSyncTime(ctx, apiType, end)
		if err != nil {
			log.Warning(fmt.Sprintf("%v set sync time error", apiType), map[string]interface{}{
				log.KeyError: err,
			})
		}
	}
}

func (s *shopifySource) syncOrders(ctx context.Context, begin, end time.Time) (int, error) {
	listOptions := &goshopify.ListOptions{
		CreatedAtMin: begin,
		CreatedAtMax: end,
	}
	c := 0
	for {
		list, pageOptions, err := s.client.Order.ListWithPagination(listOptions)
		if err != nil {
			return c, err
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
	return c, nil
}

func (s *shopifySource) syncProducts(ctx context.Context, begin, end time.Time) (int, error) {
	listOptions := &goshopify.ListOptions{
		CreatedAtMin: begin,
		CreatedAtMax: end,
	}
	c := 0
	for {
		list, pageOptions, err := s.client.Product.ListWithPagination(listOptions)
		if err != nil {
			return c, err
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
	return c, nil
}
