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

package alicloud_billing

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"

	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"

	bssopenapi "github.com/alibabacloud-go/bssopenapi-20171214/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	tutil "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/go-logr/logr"
)

const (
	EventType   = "alicloud.account_billing.daily"
	EventSource = "cloud.billing.alicloud"
)

type AlicloudBillingSource struct {
	client   *bssopenapi.Client
	config   Config
	ceClient ce.Client
	events   chan ce.Event
	ctx      context.Context
	logger   logr.Logger
}

func NewAlicloudBillingSource(ctx context.Context, ceClient ce.Client) connector.Source {
	config := getConfig(ctx)
	return &AlicloudBillingSource{
		ctx:      ctx,
		logger:   log.FromContext(ctx),
		config:   config,
		ceClient: ceClient,
		events:   make(chan ce.Event, 10),
	}
}

func (s *AlicloudBillingSource) Adapt(args ...interface{}) ce.Event {
	return ce.Event{}
}

func (s *AlicloudBillingSource) Start() error {
	conf := &openapi.Config{
		AccessKeyId:     tea.String(s.config.AccessKeyID),
		AccessKeySecret: tea.String(s.config.SecretAccessKey),
		Endpoint:        tea.String(s.config.Endpoint),
	}
	client, err := bssopenapi.NewClient(conf)
	if err != nil {
		return err
	}
	s.client = client
	ctx := s.ctx
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			//先执行一次，之后每天2点执行一次
			s.queryAccountBill()
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), s.config.PullHour, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			select {
			case <-ctx.Done():
				return
			case <-t.C:

			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.sendEvent()
	}()
	<-ctx.Done()
	close(s.events)
	wg.Wait()
	return nil
}

func (s *AlicloudBillingSource) queryAccountBill() {
	s.logger.Info("query account bill begin")
	now := time.Now()
	opt := &tutil.RuntimeOptions{}
	opt.SetAutoretry(true)
	lastDay := now.Add(time.Hour * 24 * -1)
	dayFmt := FormatTimeDay(lastDay)
	monthFmt := FormatTimeMonth(lastDay)

	var pageNum, totalSize int32
	for {
		pageNum++
		request := &bssopenapi.QueryAccountBillRequest{
			BillingCycle:     tea.String(monthFmt),
			Granularity:      tea.String("DAILY"),
			BillingDate:      tea.String(dayFmt),
			IsGroupByProduct: tea.Bool(true),
			PageNum:          tea.Int32(pageNum),
			PageSize:         tea.Int32(100),
		}
		resp, err := s.client.QueryAccountBillWithOptions(request, opt)
		if err != nil {
			s.logger.Error(err, "query account bill error")
			return
		}
		if resp.Body == nil {
			s.logger.Info("resp body is nil")
			return
		}
		if resp.Body.Success == nil {
			s.logger.Info("resp body success is nil")
			break
		}
		if !*resp.Body.Success {
			s.logger.Info("resp body success is false")
			break
		}
		data := resp.Body.Data
		if data == nil {
			s.logger.Info("resp body data is nil")
			break
		}

		for _, item := range data.Items.Item {
			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventType)
			event.SetTime(now)
			event.SetID(uuid.New().String())
			err = event.SetData(ce.ApplicationJSON, BillingData{
				VanceSource: EventSource,
				VanceType:   EventType,
				QueryAccountBillResponseBodyDataItemsItem: *item,
			})
			if err != nil {
				s.logger.Error(err, "set event data error")
				continue
			}
			s.events <- event
		}
		totalSize += *data.PageSize
		if totalSize >= *data.TotalCount {
			break
		}
	}
	s.logger.Info("get account bill end")
}

func (s *AlicloudBillingSource) getInstanceBill() {
	s.logger.Info("get instance bill begin")
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
			s.logger.Error(err, "get instance bill error")
			return
		}
		if resp.Body == nil {
			s.logger.Info("resp body is nil")
			return
		}
		if resp.Body.Success == nil {
			s.logger.Info("resp body success is nil")
			break
		}
		if !*resp.Body.Success {
			s.logger.Info("resp body success is false")
			break
		}
		if resp.Body.Data == nil {
			s.logger.Info("resp body data is nil")
			break
		}
		for _, item := range resp.Body.Data.Items {
			event := ce.NewEvent()
			event.SetSource(EventSource)
			event.SetType(EventType)
			event.SetTime(now)
			event.SetID(fmt.Sprintf("%s-%s", dayFmt, tea.StringValue(item.InstanceID)))
			err = event.SetData(ce.ApplicationJSON, item)
			if err != nil {
				s.logger.Error(err, "set event data error")
				continue
			}
			s.events <- event
		}
		nextToken = resp.Body.Data.NextToken
		if nextToken == nil || *nextToken == "" {
			break
		}
	}
	s.logger.Info("get instance bill end")
}

func (s *AlicloudBillingSource) sendEvent() {
	for event := range s.events {
		result := s.ceClient.Send(s.ctx, event)
		if !ce.IsACK(result) {
			s.logger.Error(result, "send event fail")
		} else {
			s.logger.Info("send event success", "event", event)
		}
	}
}
