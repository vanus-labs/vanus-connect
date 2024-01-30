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
	"time"

	goshopify "github.com/bold-commerce/go-shopify/v3"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/cdk-go/store"
)

var _ cdkgo.Source = &shopifySource{}

func NewSource() cdkgo.HTTPSource {
	return &shopifySource{
		events:    make(chan *cdkgo.Tuple, 1024),
		eventType: map[ApiType]struct{}{},
	}
}

var syncApiArr = []ApiType{OrderApi, ProductApi}

type shopifySource struct {
	events        chan *cdkgo.Tuple
	config        *shopifyConfig
	client        *goshopify.Client
	syncBeginTime time.Time
	eventType     map[ApiType]struct{}
	logger        zerolog.Logger
	store         store.KVStore
	shopifyApp    goshopify.App
	count         int64
}

func (s *shopifySource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.store = store.FromContext(ctx)
	s.config = cfg.(*shopifyConfig)
	t, err := time.Parse("2006-01-02", s.config.SyncBeginDate)
	if err != nil {
		return err
	}
	s.syncBeginTime = t
	if s.config.DelaySecond <= 0 {
		s.config.DelaySecond = 5
	}
	for _, et := range s.config.EventTypes {
		s.eventType[et] = struct{}{}
	}
	s.shopifyApp = goshopify.App{ApiSecret: s.config.ClientSecret}
	s.client = goshopify.NewClient(s.shopifyApp, s.config.ShopName, s.config.ApiAccessToken, goshopify.WithVersion("2024-01"))
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
	s.sync(ctx)
}

func (s *shopifySource) isSync(t ApiType) bool {
	if len(s.eventType) == 0 {
		return true
	}
	_, exist := s.eventType[t]
	return exist
}
func (s *shopifySource) sync(ctx context.Context) error {
	time.Sleep(time.Second * time.Duration(s.config.DelaySecond))
	for _, apiType := range syncApiArr {
		if !s.isSync(apiType) {
			s.logger.Info().Interface("api", apiType).Msg("sync ignore")
			continue
		}
		need, err := s.isApiNeedSync(ctx, apiType)
		if err != nil {
			s.logger.Error().Err(err).Interface("api", apiType).Msg("check need sync error")
			continue
		}
		if !need {
			s.logger.Info().Interface("api", apiType).Msg("has sync")
			continue
		}
		begin := s.syncBeginTime
		end := time.Now().UTC()
		s.logger.Info().
			Time("begin", begin).
			Time("end", end).
			Interface("api", apiType).
			Msg("sync data begin")
		var c int
		switch apiType {
		case OrderApi:
			c, err = s.syncOrders(ctx, begin, end)
		case ProductApi:
			c, err = s.syncProducts(ctx, begin, end)
		}
		if err != nil {
			s.logger.Warn().Err(err).
				Interface("api", apiType).
				Msg("sync data error")
			continue
		}
		s.logger.Info().
			Int("count", c).
			Interface("api", apiType).
			Msg("sync data success")
		err = s.setApiSync(ctx, apiType)
		if err != nil {
			s.logger.Warn().Err(err).
				Interface("api", apiType).
				Msg("set sync error")
		}
	}
	return nil
}

func (s *shopifySource) syncOrders(_ context.Context, begin, end time.Time) (int, error) {
	var listOptions interface{}
	listOptions = &goshopify.OrderListOptions{
		ListOptions: goshopify.ListOptions{
			CreatedAtMin: begin,
			CreatedAtMax: end,
			Limit:        250,
		},
		Status: "any",
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

func (s *shopifySource) syncProducts(_ context.Context, begin, end time.Time) (int, error) {
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
