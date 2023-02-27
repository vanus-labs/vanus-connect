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
	"strconv"
	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ cdkgo.Sink = &GoogleSheetSink{}

func NewGoogleSheetSink() cdkgo.Sink {
	return &GoogleSheetSink{}
}

type GoogleSheetSink struct {
	config *GoogleSheetConfig
	client *sheets.Service
	spreadSheetId string
	sheetName  string
}

func (s *GoogleSheetSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	// TODO
	s.config = cfg.(*GoogleSheetConfig)

	// Authenticate and get configuration
	config, err := google.JWTConfigFromJSON([]byte(s.config.Credentials), "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			return err
		}

	//Create Client
	client := config.Client(context.Background())

	//Create Service using Client
	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
        return err
	}
	s.client = srv

	//Initialize Sheet ID & Spreadsheet ID
	spreadSheetUrl := s.config.Sheet_url
	
	//Get Sheet ID

	sheetId, err := strconv.Atoi(spreadSheetUrl[93:94])
	if err != nil {
        fmt.Errorf("Failed to get sheet ID: %s", err)
        return err

	}

	// Get SpreadSheet ID
	spreadSheetID := spreadSheetUrl[39:83]

	s.spreadSheetId = spreadSheetID

	//Get SheetName from SpreadSheetID
	res1, err := s.client.Spreadsheets.Get(spreadSheetID).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil {
        fmt.Errorf("Failed to get sheet name: %s", err)
        return err
	}

	sheetName := ""
	for _, v := range res1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			break
		}
	}

	s.sheetName = sheetName


	return nil
}

func (s *GoogleSheetSink) Name() string {
	// TODO
	return "GoogleSheetSink"
}

func (s *GoogleSheetSink) Destroy() error {
	// TODO
	return nil
}

func (s *GoogleSheetSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	// TODO
	for _, event := range events {

		s.saveDataToSpreadsheet(event)
	}
	return cdkgo.SuccessResult
}



func (s *GoogleSheetSink) saveDataToSpreadsheet(event *ce.Event) {

	// Receive any kind of Cloud Event
	sheetRow := make(map[string]interface{})
	json.Unmarshal(event.Data(), &sheetRow)
	
	var values []interface{}
	for _, v := range sheetRow {
		values = append(values, v)
	}

	//Insert Row Value
	row := &sheets.ValueRange{
		Values: [][] interface{}{ values },
	}

	response, err := s.client.Spreadsheets.Values.Append(s.spreadSheetId, s.sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
		if err != nil || response.HTTPStatusCode != 200 {
		fmt.Errorf("Failed to Insert new row %s", err)
		return
	}




}


