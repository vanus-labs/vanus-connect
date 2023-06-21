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
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ cdkgo.Source = &googleSheetsSource{}

func NewGoogleSheetsSource() cdkgo.Source {
	return &googleSheetsSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type googleSheetsSource struct {
	config *googleSheetsConfig
	events chan *cdkgo.Tuple
}

func (s *googleSheetsSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*googleSheetsConfig)

	// Set up the Sheets API client
	var err error

	credentials := &google.Credentials{}
	err = json.Unmarshal([]byte(s.config.Credentials), credentials)
	if err != nil {
		log.Fatalf("Unable to parse credentials: %v", err)
	}

	s.config.srv, err = sheets.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	go s.loopProduceEvent()
	return nil
}

func (s *googleSheetsSource) Name() string {
	// TODO
	return "GoogleSheetsSource"
}

func (s *googleSheetsSource) Destroy() error {
	// TODO
	return nil
}

func (s *googleSheetsSource) Chan() <-chan *cdkgo.Tuple {
	// TODO
	return s.events
}

func (s *googleSheetsSource) loopProduceEvent() *ce.Event {

	poolingInterval := 10 * time.Second
	pool := make(chan struct{})

	for {
		select {
		case <-pool:
			// Perform pooling tasks
			data, err := fetchNewData(s.config.srv, s.config)
			if err != nil {
				log.Printf("Unable to fetch data: %v", err)
			} else if len(data) > 0 {
				jsonString, err := json.Marshal(data)
				if err != nil {
					log.Printf("Unable to marshal JSON: %v", err)
				}

				event := s.makeEvent(jsonString)
				b, _ := json.Marshal(event)
				success := func() {
					fmt.Println("send event success: " + string(b))
				}
				failed := func(err error) {
					fmt.Println("send event failed: " + string(b) + ", error: " + err.Error())
				}
				s.events <- cdkgo.NewTuple(event, success, failed)
			} else {
				log.Println("No new data found")
			}
		case <-time.After(poolingInterval):
			// Perform pooling tasks

		}
	}
}

func (s *googleSheetsSource) makeEvent(data []byte) *ce.Event {
	event := ce.NewEvent()
	event.SetSource("googleSheetsSource")
	event.SetData(ce.ApplicationJSON, data)

	return &event
}

var lastRetrievedRow int

func fetchNewData(srv *sheets.Service, s *googleSheetsConfig) ([]map[string]interface{}, error) {

	// Fetch the latest data from the Google Sheet
	resp, err := srv.Spreadsheets.Values.Get(s.SheetID, s.SheetName).ValueRenderOption("UNFORMATTED_VALUE").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	// Retrieve the values
	values := resp.Values

	// Check if there are new rows
	if len(values) <= lastRetrievedRow {
		return nil, nil // No new data
	}

	// Check if the first row contains headers
	if len(values[0]) == 0 {
		return nil, fmt.Errorf("no header row found in the spreadsheet")
	}

	// Convert the new data to JSON
	jsonData := make([]map[string]interface{}, 0)
	for i := lastRetrievedRow; i < len(values); i++ {
		row := values[i]
		data := make(map[string]interface{})
		for j, cell := range row {
			if j < len(values[0]) {
				data[values[0][j].(string)] = cell
			}
		}
		jsonData = append(jsonData, data)
	}

	// Update the last retrieved row
	lastRetrievedRow = len(values)

	return jsonData, nil
}
