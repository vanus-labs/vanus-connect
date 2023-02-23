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
	"log"
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
}

func (s *GoogleSheetSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	// TODO
	s.config = cfg.(*GoogleSheetConfig)
	
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

		//Create API Context
	ctx := context.Background()

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON([]byte(s.config.Credentials), "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			log.Fatalf("Failed to authenticate google service accounts key %v", err)
			return
		}

	//Create Client
	client := config.Client(ctx)

	//Create Service using Client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
        log.Fatalf("Failed to Create Service Account %v",err)
        return
	}

	//Initialize Sheet ID & Spreadsheet ID
	
	spreadSheetUrl := s.config.Sheet_url
	
	sheetId, err := strconv.Atoi(spreadSheetUrl[93:94])
	if err != nil {
        log.Fatalf("Failed to Convert String %v",err)
        return
	}

	spreadSheetId := spreadSheetUrl[39:83]

	//Get SheetName from SpreadSheetID
	res1, err := srv.Spreadsheets.Get(spreadSheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil {
        log.Fatalf("Failed to get SheetName %v",err)
        return
	}

	sheetName := ""
	for _, v := range res1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			break
		}
	}


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

	response, err := srv.Spreadsheets.Values.Append(spreadSheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(ctx).Do()
		if err != nil || response.HTTPStatusCode != 200 {
		log.Fatalf("Failed to Append Value to Spreadsheet %v",err)
		return
	}


}


