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
	"net/http"
	"strings"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const xvSheetName = "xvsheetname"

var (
	errInvalidSheetName = cdkgo.NewResult(http.StatusBadRequest, fmt.Sprintf("extension %s is invalid", xvSheetName))
)

var _ cdkgo.Sink = &GoogleSheetSink{}

func NewGoogleSheetSink() cdkgo.Sink {
	return &GoogleSheetSink{
		sheetHeaders: map[string][]string{},
		sheetIDs:     map[string]int64{},
	}
}

type GoogleSheetSink struct {
	client           *sheets.Service
	spreadsheetID    string
	defaultSheetName string
	sheetIDs         map[string]int64    // key: sheetName, value: sheetID
	sheetHeaders     map[string][]string // key: sheetName, value: sheet headers which is on the first row
}

func (s *GoogleSheetSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {

	config := cfg.(*GoogleSheetConfig)
	sheetURL := strings.TrimSpace(config.SheetURL)
	sheetID, err := getSheetID(sheetURL)
	if err != nil {
		return errors.Wrap(err, "ge sheet id error")
	}
	spreadsheetID, err := getSpreadsheetID(sheetURL)
	if err != nil {
		return errors.Wrap(err, "ge spread sheet id error")
	}
	s.spreadsheetID = spreadsheetID
	// new sheet Service
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsJSON([]byte(config.Credentials)))
	if err != nil {
		return errors.Wrap(err, "new sheet service with credential error")
	}
	s.client = srv

	// get SheetName from SpreadSheetID
	spreadSheet, err := s.client.Spreadsheets.Get(s.spreadsheetID).Do()
	if err != nil {
		return errors.Wrap(err, "spreadsheets get error")
	}
	for _, sheet := range spreadSheet.Sheets {
		s.sheetIDs[sheet.Properties.Title] = sheet.Properties.SheetId
		if sheet.Properties.SheetId == sheetID {
			s.defaultSheetName = sheet.Properties.Title
		}
	}
	if s.defaultSheetName == "" {
		return fmt.Errorf("sheetURL sheetID %d no exist", sheetID)
	}
	return nil
}

func (s *GoogleSheetSink) createSheet(sheetName string) error {
	sheetAdd := sheets.AddSheetRequest{
		Properties: &sheets.SheetProperties{
			Hidden:    false,
			SheetType: "GRID",
			Title:     sheetName,
		},
	}

	updateRequests := sheets.BatchUpdateSpreadsheetRequest{
		IncludeSpreadsheetInResponse: true,
		Requests:                     []*sheets.Request{{AddSheet: &sheetAdd}},
		ResponseIncludeGridData:      false,
	}

	resp, err := s.client.Spreadsheets.BatchUpdate(s.spreadsheetID, &updateRequests).Do()
	if err != nil {
		return errors.Wrap(err, "create sheet error")
	}
	for _, sheet := range resp.UpdatedSpreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			s.sheetIDs[sheetName] = sheet.Properties.SheetId
			break
		}
	}
	return nil
}

func (s *GoogleSheetSink) Name() string {
	return "GoogleSheetSink"
}

func (s *GoogleSheetSink) Destroy() error {
	return nil
}

func (s *GoogleSheetSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		r := s.saveDataToSpreadsheet(event)
		if r != cdkgo.SuccessResult {
			return r
		}
	}
	return cdkgo.SuccessResult
}

func (s *GoogleSheetSink) saveDataToSpreadsheet(event *ce.Event) cdkgo.Result {
	var sheetName string
	// get sheetName
	v, exist := event.Extensions()[xvSheetName]
	if exist {
		str, ok := v.(string)
		if !ok {
			return errInvalidSheetName
		}
		sheetName = str
	} else {
		sheetName = s.defaultSheetName
	}
	// check sheetName's sheetID is existed
	if _, exist = s.sheetIDs[sheetName]; !exist {
		// sheetName no exist sheetID, create the sheetName
		err := s.createSheet(sheetName)
		if err != nil {
			log.Error("create sheet error", map[string]interface{}{
				log.KeyError: err,
				"sheetName":  sheetName,
			})
			return cdkgo.NewResult(http.StatusInternalServerError, "create sheet error")
		}
	}
	// sheet row
	sheetRow := make(map[string]interface{})
	err := json.Unmarshal(event.Data(), &sheetRow)
	if err != nil {
		log.Error("data json unmarshal error", map[string]interface{}{
			log.KeyError: err,
			"data":       string(event.Data()),
		})
		return cdkgo.NewResult(http.StatusBadRequest, "data json unmarshal error")
	}
	// get sheet headers
	headers, err := s.getHeader(sheetName, sheetRow)
	if err != nil {
		log.Error("get header error", map[string]interface{}{
			log.KeyError: err,
			"sheetName":  sheetName,
		})
		return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
	}
	// make cell data
	var values []*sheets.CellData
	for _, key := range headers {
		values = append(values, &sheets.CellData{
			UserEnteredValue: sheetValue(sheetRow[key]),
		})
	}
	// append data
	err = s.appendData(sheetName, values)
	if err != nil {
		log.Error("append sheet data error", map[string]interface{}{
			log.KeyError: err,
			"sheetName":  sheetName,
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "append data error")
	}
	return cdkgo.SuccessResult
}

func (s *GoogleSheetSink) getHeader(sheetName string, sheetRow map[string]interface{}) ([]string, error) {
	headers, exist := s.sheetHeaders[sheetName]
	if exist && len(headers) != 0 {
		return headers, nil
	}
	// get header in first row
	resp, err := s.client.Spreadsheets.Values.Get(s.spreadsheetID, fmt.Sprintf("%s!A1:Z1", sheetName)).Do()
	if err != nil {
		return nil, errors.Wrap(err, "get sheet header error")
	}
	if len(resp.Values) == 0 {
		// make header
		for k := range sheetRow {
			headers = append(headers, k)
		}
		// insert headers
		var values []*sheets.CellData
		for _, key := range headers {
			values = append(values, &sheets.CellData{
				UserEnteredValue: sheetValue(key),
			})
		}
		err = s.appendData(sheetName, values)
		if err != nil {
			return nil, errors.Wrap(err, "insert sheet header error")
		}
	} else {
		for _, value := range resp.Values[0] {
			headers = append(headers, fmt.Sprintf("%v", value))
		}
	}
	s.sheetHeaders[sheetName] = headers
	return headers, nil
}

func (s *GoogleSheetSink) appendData(sheetName string, values []*sheets.CellData) error {
	addValues := sheets.AppendCellsRequest{
		Fields:  "*",
		Rows:    []*sheets.RowData{{Values: values}},
		SheetId: s.sheetIDs[sheetName],
	}

	updateRequests := sheets.BatchUpdateSpreadsheetRequest{
		IncludeSpreadsheetInResponse: false,
		Requests:                     []*sheets.Request{{AppendCells: &addValues}},
		ResponseIncludeGridData:      false,
	}
	_, err := s.client.Spreadsheets.BatchUpdate(s.spreadsheetID, &updateRequests).Do()
	if err != nil {
		return err
	}
	return nil
}
