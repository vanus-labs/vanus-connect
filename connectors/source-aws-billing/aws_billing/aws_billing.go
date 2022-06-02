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

package aws_billing

import (
	"context"
	"sync"
	"time"

	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

const (
	EventType   = "aws.service.daily"
	EventSource = "cloud.billing.aws"
)

type AwsBillingSource struct {
	client   *costexplorer.Client
	config   Config
	ceClient ce.Client
	events   chan ce.Event
	ctx      context.Context
	logger   logr.Logger
}

func NewAwsBillingSource(ctx context.Context, ceClient ce.Client) connector.Source {
	config := getConfig(ctx)
	return &AwsBillingSource{
		client:   newCostExplorerClient(config),
		ceClient: ceClient,
		events:   make(chan ce.Event, 10),
		ctx:      ctx,
		logger:   log.FromContext(ctx),
	}
}

func newCostExplorerClient(config Config) *costexplorer.Client {
	opt := costexplorer.Options{
		Region:           "us-east-1",
		EndpointResolver: costexplorer.EndpointResolverFromURL(config.Endpoint),
	}
	opt.Credentials = credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, "")
	return costexplorer.New(opt)
}

func (s *AwsBillingSource) Adapt(args ...interface{}) ce.Event {
	return ce.Event{}
}

func (s *AwsBillingSource) Start() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			s.getCost()
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), s.config.PullHour, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			select {
			case <-s.ctx.Done():
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
	<-s.ctx.Done()
	close(s.events)
	wg.Wait()
	return nil
}

func (s *AwsBillingSource) getCost() {
	ctx := s.ctx
	s.logger.Info("get cost begin")
	now := time.Now()
	endDayFmt := FormatTimeDay(now)
	dayFmt := FormatTimeDay(now.Add(time.Hour * 24 * -1))
	var nextToken *string
	for {
		input := &costexplorer.GetCostAndUsageInput{
			Granularity: types.GranularityDaily,
			TimePeriod: &types.DateInterval{
				Start: aws.String(dayFmt),
				End:   aws.String(endDayFmt),
			},
			//GroupDefinitionTypeDimension type
			//AZ, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, PURCHASE_TYPE, SERVICE, USAGE_TYPE, PLATFORM,
			//TENANCY, RECORD_TYPE, LEGAL_ENTITY_NAME, INVOICING_ENTITY, DEPLOYMENT_OPTION,
			//DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, REGION, BILLING_ENTITY,
			//RESERVATION_ID, SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM
			GroupBy: []types.GroupDefinition{
				{Key: aws.String("SERVICE"), Type: types.GroupDefinitionTypeDimension},
			},
			Metrics:       []string{"BlendedCost"},
			NextPageToken: nextToken,
		}
		output, err := s.client.GetCostAndUsage(ctx, input)
		if err != nil {
			s.logger.Error(err, "get cost and usage error")
			return
		}
		for _, result := range output.ResultsByTime {
			for _, item := range result.Groups {
				data := BillingData{
					VanceSource: EventSource,
					VanceType:   EventType,
					Date:        dayFmt,
					Service:     item.Keys[0],
				}
				if cost, exist := item.Metrics["BlendedCost"]; exist {
					data.Amount = cost.Amount
					data.Unit = cost.Unit
				}
				event := ce.NewEvent()
				event.SetID(uuid.New().String())
				event.SetType(EventType)
				event.SetSource(EventSource)
				event.SetTime(now)
				err = event.SetData(ce.ApplicationJSON, data)
				if err != nil {
					s.logger.Error(err, "set event data error")
					continue
				}
				s.events <- event
			}
		}
		nextToken = output.NextPageToken
		if nextToken == nil || *nextToken == "" {
			break
		}
	}
	s.logger.Info("get cost end")
}

func (s *AwsBillingSource) sendEvent() {
	for event := range s.events {
		result := s.ceClient.Send(s.ctx, event)
		if !ce.IsACK(result) {
			s.logger.Error(result, "send event fail")
		} else {
			s.logger.Info("send event success", "event", event)
		}
	}
}
