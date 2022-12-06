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
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
)

const (
	EventType   = "aws.service.daily"
	EventSource = "cloud.aws.billing"
)

type awsBillingSource struct {
	client *costexplorer.Client
	config *billingConfig
	events chan *cdkgo.Tuple
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func Source() cdkgo.Source {
	return &awsBillingSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

func newCostExplorerClient(config *billingConfig) *costexplorer.Client {
	opt := costexplorer.Options{
		Region:           "us-east-1",
		EndpointResolver: costexplorer.EndpointResolverFromURL(config.Endpoint),
	}
	opt.Credentials = credentials.NewStaticCredentialsProvider(config.Secret.AccessKeyID, config.Secret.SecretAccessKey, "")
	return costexplorer.New(opt)
}

func (s *awsBillingSource) Name() string {
	return "AwsBillingSource"
}

func (s *awsBillingSource) Destroy() error {
	s.wg.Wait()
	close(s.events)
	return nil
}

func (s *awsBillingSource) Initialize(_ context.Context, config cdkgo.ConfigAccessor) error {
	cfg := config.(*billingConfig)
	s.config = cfg
	if cfg.PullHour <= 0 || cfg.PullHour >= 24 {
		cfg.PullHour = 2
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = "https://ce.us-east-1.amazonaws.com"
	}
	s.ctx, s.cancel = context.WithCancel(context.TODO())
	s.client = newCostExplorerClient(cfg)
	s.start()
	return nil
}

func (s *awsBillingSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *awsBillingSource) start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
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
}

func (s *awsBillingSource) getCost() {
	ctx := s.ctx
	log.Info("get cost begin", nil)
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
			log.Error("get cost and usage error", map[string]interface{}{
				log.KeyError: err,
			})
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
				_ = event.SetData(ce.ApplicationJSON, data)
				s.events <- &cdkgo.Tuple{
					Event: &event,
				}
			}
		}
		nextToken = output.NextPageToken
		if nextToken == nil || *nextToken == "" {
			break
		}
	}
	log.Info("get cost end", nil)
}
