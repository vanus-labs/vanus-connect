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
	"sync"
	"time"

	bssopenapi "github.com/alibabacloud-go/bssopenapi-20171214/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	tutil "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	EventSource = "cloud.alibaba.billing"
	EventType   = "alibaba.billing.daily"
)

type alicloudBillingSource struct {
	client   *bssopenapi.Client
	config   *billingConfig
	events   chan *cdkgo.Tuple
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	timeZone *time.Location
}

func Source() cdkgo.Source {
	return &alicloudBillingSource{
		events: make(chan *cdkgo.Tuple, 10),
		cancel: func() {},
	}
}

func (s *alicloudBillingSource) Initialize(ctx context.Context, config cdkgo.ConfigAccessor) error {
	cfg := config.(*billingConfig)
	s.config = cfg
	if cfg.PullHour <= 0 || cfg.PullHour >= 24 {
		cfg.PullHour = 2
	}
	if cfg.PullZone == "" {
		s.timeZone = time.UTC
	} else {
		loc, err := time.LoadLocation(cfg.PullZone)
		if err != nil {
			return err
		}
		s.timeZone = loc
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "business.aliyuncs.com"
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	conf := &openapi.Config{
		AccessKeyId:     tea.String(s.config.Secret.AccessKeyID),
		AccessKeySecret: tea.String(s.config.Secret.SecretAccessKey),
		Endpoint:        tea.String(s.config.Endpoint),
	}
	client, err := bssopenapi.NewClient(conf)
	if err != nil {
		return err
	}
	s.client = client
	// check
	now := time.Now()
	lastDay := now.Add(time.Hour * 24 * -1)
	dayFmt := FormatTimeDay(lastDay)
	monthFmt := FormatTimeMonth(lastDay)
	_, err = s.queryAccountBill(monthFmt, dayFmt, 1)
	if err != nil {
		return err
	}
	s.start()
	return nil
}

func (s *alicloudBillingSource) Name() string {
	return "AliCloudBillingSource"
}

func (s *alicloudBillingSource) Destroy() error {
	s.cancel()
	s.wg.Wait()
	close(s.events)
	return nil
}

func (s *alicloudBillingSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *alicloudBillingSource) start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.getBilling()
		now := time.Now().In(s.timeZone)
		next := now.Add(time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		t := time.NewTicker(next.Sub(now))
		select {
		case <-s.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			s.getBilling()
		}
		t.Stop()
		tk := time.NewTicker(time.Hour)
		defer tk.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-tk.C:
				s.getBilling()
			}
		}
	}()
}

func (s *alicloudBillingSource) getBilling() {
	log.Info("query account bill begin", nil)
	now := time.Now().In(s.timeZone)
	if now.Hour() != s.config.PullHour {
		log.Info("query account bill hour is not match", nil)
	}
	lastDay := now.Add(time.Hour * 24 * -1)
	dayFmt := FormatTimeDay(lastDay)
	monthFmt := FormatTimeMonth(lastDay)

	var pageNum, totalSize int32
	for {
		pageNum++
		resp, err := s.queryAccountBill(monthFmt, dayFmt, pageNum)
		if err != nil {
			log.Error("query account bill error", map[string]interface{}{
				log.KeyError: err,
			})
			return
		}
		if resp.Body == nil {
			log.Info("resp body is nil", nil)
			return
		}
		if resp.Body.Success == nil {
			log.Info("resp body success is nil", nil)
			break
		}
		if !*resp.Body.Success {
			log.Info("resp body success is false", nil)
			break
		}
		data := resp.Body.Data
		if data == nil {
			log.Info("resp body data is nil", nil)
			break
		}

		for _, item := range data.Items.Item {
			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventType)
			event.SetTime(now)
			event.SetID(uuid.New().String())
			_ = event.SetData(ce.ApplicationJSON, BillingData{
				QueryAccountBillResponseBodyDataItemsItem: *item,
			})
			s.events <- &cdkgo.Tuple{
				Event: &event,
			}
		}
		totalSize += *data.PageSize
		if totalSize >= *data.TotalCount {
			break
		}
	}
	log.Info("get account bill end", nil)
}

func (s *alicloudBillingSource) queryAccountBill(monthFmt, dayFmt string, pageNum int32) (*bssopenapi.QueryAccountBillResponse, error) {
	request := &bssopenapi.QueryAccountBillRequest{
		BillingCycle:     tea.String(monthFmt),
		Granularity:      tea.String("DAILY"),
		BillingDate:      tea.String(dayFmt),
		IsGroupByProduct: tea.Bool(true),
		PageNum:          tea.Int32(pageNum),
		PageSize:         tea.Int32(100),
	}
	opt := &tutil.RuntimeOptions{}
	opt.SetAutoretry(true)
	return s.client.QueryAccountBillWithOptions(request, opt)
}

func (s *alicloudBillingSource) getInstanceBill() {
	log.Info("get instance bill begin", nil)
	opt := &tutil.RuntimeOptions{}
	opt.SetAutoretry(true)
	now := time.Now()
	lastDay := now.Add(time.Hour * 24 * -1)
	dayFmt := FormatTimeDay(lastDay)
	monthFmt := FormatTimeMonth(lastDay)
	var nextToken *string
	for {
		request := &bssopenapi.DescribeInstanceBillRequest{
			BillingCycle: tea.String(monthFmt),
			Granularity:  tea.String("DAILY"),
			BillingDate:  tea.String(dayFmt),
			NextToken:    nextToken,
			MaxResults:   tea.Int32(300),
		}
		resp, err := s.client.DescribeInstanceBillWithOptions(request, opt)
		if err != nil {
			log.Error("get instance bill error", map[string]interface{}{
				log.KeyError: err,
			})
			return
		}
		if resp.Body == nil {
			log.Info("resp body is nil", nil)
			return
		}
		if resp.Body.Success == nil {
			log.Info("resp body success is nil", nil)
			break
		}
		if !*resp.Body.Success {
			log.Info("resp body success is false", nil)
			break
		}
		if resp.Body.Data == nil {
			log.Info("resp body data is nil", nil)
			break
		}
		for _, item := range resp.Body.Data.Items {
			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventType)
			event.SetTime(now)
			event.SetID(fmt.Sprintf("%s-%s", dayFmt, tea.StringValue(item.InstanceID)))
			_ = event.SetData(ce.ApplicationJSON, item)
			s.events <- &cdkgo.Tuple{
				Event: &event,
			}
		}
		nextToken = resp.Body.Data.NextToken
		if nextToken == nil || *nextToken == "" {
			break
		}
	}
	log.Info("get instance bill end", nil)
}
