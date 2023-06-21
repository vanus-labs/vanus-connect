// Copyright 2023 Linkall Inc.
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
	"log"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
	"golang.org/x/oauth2/google"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

var _ cdkgo.Source = &googleAnalyticsSource{}

func GoogleAnalyticsSource() cdkgo.Source {
	return &googleAnalyticsSource{

		events: make(chan *cdkgo.Tuple, 100),
	}
}

type googleAnalyticsSource struct {
	config  *googleAnalyticsConfig
	events  chan *cdkgo.Tuple
	svc     *analyticsdata.Service
	request *analyticsdata.RunReportRequest
}

func (s *googleAnalyticsSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*googleAnalyticsConfig)

	// Create a new AnalyticsData service client with the service account credentials
	creds, err := google.CredentialsFromJSON(ctx, []byte(s.config.Credentials), analyticsdata.AnalyticsReadonlyScope)
	if err != nil {
		log.Fatalf("Failed to create credentials: %v", err)
	}

	// Create a  new AnalyticsData service client
	service, err := analyticsdata.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to create AnalyticsData service client: %v", err)
	}

	s.svc = service

	req := &analyticsdata.RunReportRequest{
		Dimensions: []*analyticsdata.Dimension{
			{
				Name: "continent",
			},

			{
				Name: "country",
			},

			{
				Name: "city",
			},

			{
				Name: "method",
			},

			{
				Name: "newVsReturning",
			},

			{
				Name: "sessionSource",
			},

			{
				Name: "transactionId",
			},
		},
		Metrics: []*analyticsdata.Metric{
			{
				Name: "eventCount",
			},

			{
				Name: "activeUsers",
			},

			{
				Name: "addToCarts",
			},

			{
				Name: "checkouts",
			},

			{
				Name: "ecommercePurchases",
			},

			{
				Name: "eventValue",
			},

			{
				Name: "firstTimePurchasers",
			},

			{
				Name: "purchaseRevenue",
			},

			{
				Name: "totalAdRevenue",
			},
		},

		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: s.config.Start_date,
				EndDate:   s.config.End_date,
			},
		},
	}

	s.request = req

	go s.loopProduceEvent()
	return nil
}

func (s *googleAnalyticsSource) Name() string {

	return "GoogleAnalyticsSource"
}

func (s *googleAnalyticsSource) Destroy() error {

	return nil
}

func (s *googleAnalyticsSource) Chan() <-chan *cdkgo.Tuple {

	return s.events
}

func (s *googleAnalyticsSource) loopProduceEvent() *ce.Event {

	for {
		event := s.makeEvent()
		b, _ := json.Marshal(event)
		success := func() {
			fmt.Println("send event success: " + string(b))
		}
		failed := func(err error) {
			fmt.Println("send event failed: " + string(b) + ", error: " + err.Error())
		}
		s.events <- cdkgo.NewTuple(event, success, failed)
	}
}

func (s *googleAnalyticsSource) makeEvent() *ce.Event {
	event := ce.NewEvent()
	event.SetSource("GoogleAnalyticsSource")
	event.SetData(ce.ApplicationJSON, map[string]interface{}{})

	//Retrieve Event
	resp, err := s.svc.Properties.RunReport(fmt.Sprintf("properties/%s", s.config.PropertyID), s.request).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}

	if len(resp.Rows) == 0 {
		fmt.Println("No result found. Confirm end date and try again")
	}

	if len(resp.Rows) > 0 {
		row := resp.Rows[0]
		dimensions := make(map[string]string)
		metrics := make(map[string]string)

		for i, dimension := range row.DimensionValues {
			dimensionName := s.request.Dimensions[i].Name
			dimensions[dimensionName] = dimension.Value
		}

		for i, metric := range row.MetricValues {
			metricName := s.request.Metrics[i].Name
			metrics[metricName] = metric.Value
		}

		data, err := json.Marshal(map[string]interface{}{
			"Dimensions": dimensions,
			"Metrics":    metrics,
		})
		if err != nil {
			log.Fatalf("Failed to marshal data to JSON: %v", err)
		}

		fmt.Println(string(data))
	}

	return &event
}
