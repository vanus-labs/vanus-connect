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
}

func (s *shopifySource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*shopifyConfig)
	t, err := time.Parse("2006-01-02", s.config.SyncBeginDate)
	if err != nil {
		return err
	}
	s.syncBeginTime = t
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
	t := time.NewTicker(time.Hour * 24)
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
	beginTime, err := s.getSyncBeginTime(ctx)
	if err != nil {
		return errors.Wrap(err, "get sync begin timer error")
	}
	now := time.Now()
	listOptions := &goshopify.ListOptions{
		PageInfo:     fmt.Sprintf("%d", time.Now().Unix()),
		Page:         1,
		CreatedAtMin: beginTime,
		CreatedAtMax: now,
	}
	for {
		orders, pageOptions, err := s.client.Order.ListWithPagination(listOptions)
		if err != nil {
			return errors.Wrap(err, "list order failed")
		}
		if len(orders) == 0 {
			break
		}
		s.orderEvent(orders)
		if pageOptions == nil || pageOptions.NextPageOptions == nil {
			break
		}
		listOptions = pageOptions.NextPageOptions
	}
	s.setSyncTime(ctx, now)
	return nil
}
