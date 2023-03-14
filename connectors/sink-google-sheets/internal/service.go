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
	"fmt"

	"github.com/pkg/errors"
	"github.com/vanus-labs/cdk-go/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	headerNotExistErr = errors.New("header not exist")
)

type GoogleSheetService struct {
	client        *sheets.Service
	spreadsheetID string
	sheetIDs      map[string]int64          // key: sheetName, value: sheetID
	sheetHeaders  map[string]map[string]int // key: sheetName, value: sheet headers
}

func newGoogleSheetService(spreadSheetID, credentialsJSON string) (*GoogleSheetService, error) {
	service := &GoogleSheetService{
		sheetHeaders: map[string]map[string]int{},
		sheetIDs:     map[string]int64{},
	}
	err := service.init(spreadSheetID, credentialsJSON)
	return service, err
}

func (s *GoogleSheetService) init(spreadSheetID, credentialsJSON string) error {
	s.spreadsheetID = spreadSheetID
	// new sheet Service
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsJSON([]byte(credentialsJSON)))
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
	}
	return nil
}

func (s *GoogleSheetService) getSheetName(sheetID int64) string {
	for k, v := range s.sheetIDs {
		if v == sheetID {
			return k
		}
	}
	return ""
}
func (s *GoogleSheetService) getHeader(sheetName string) (map[string]int, error) {
	headers, exist := s.sheetHeaders[sheetName]
	if exist && len(headers) != 0 {
		return headers, nil
	}
	resp, err := s.client.Spreadsheets.Values.Get(s.spreadsheetID, fmt.Sprintf("%s!1:1", sheetName)).Do()
	if err != nil {
		return nil, errors.Wrap(err, "get sheet header error")
	}
	if len(resp.Values) == 0 {
		return nil, headerNotExistErr
	}
	headers = make(map[string]int, len(resp.Values[0]))
	for index, value := range resp.Values[0] {
		columnName := fmt.Sprintf("%v", value)
		headers[columnName] = index
	}
	s.sheetHeaders[sheetName] = headers
	return headers, nil
}

func (s *GoogleSheetService) insertHeader(ctx context.Context, sheetName string, headers map[string]int) error {
	// insert headers
	values := make([]interface{}, len(headers))
	for key, index := range headers {
		values[index] = key
	}
	err := s.appendData(ctx, sheetName, values)
	if err != nil {
		return errors.Wrap(err, "insert sheet header error")
	}
	s.sheetHeaders[sheetName] = headers
	return nil
}

func (s *GoogleSheetService) updateHeader(ctx context.Context, sheetName string, headers map[string]int) error {
	// update headers
	values := make([]interface{}, len(headers))
	for key, index := range headers {
		values[index] = key
	}
	err := s.updateData(ctx, sheetName, 1, values)
	if err != nil {
		return errors.Wrap(err, "update sheet header error")
	}
	s.sheetHeaders[sheetName] = headers
	return nil
}

func (s *GoogleSheetService) createSheetIfNotExist(ctx context.Context, sheetName string) error {
	if _, exist := s.sheetIDs[sheetName]; !exist {
		// sheetName no exist sheetID, create the sheetName
		err := s.createSheet(ctx, sheetName)
		if err != nil {
			log.Error("create sheet error", map[string]interface{}{
				log.KeyError: err,
				"sheetName":  sheetName,
			})
			return err
		}
	}
	return nil
}

func (s *GoogleSheetService) createSheet(ctx context.Context, sheetName string) error {
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

	resp, err := s.client.Spreadsheets.BatchUpdate(s.spreadsheetID, &updateRequests).Context(ctx).Do()
	if err != nil {
		return errors.Wrap(err, "create sheet api error")
	}
	for _, sheet := range resp.UpdatedSpreadsheet.Sheets {
		if sheet.Properties.Title == sheetName {
			s.sheetIDs[sheetName] = sheet.Properties.SheetId
			return nil
		}
	}
	return errors.New("create sheet failed")
}

func (s *GoogleSheetService) appendData(ctx context.Context, sheetName string, values []interface{}) error {
	_, err := s.client.Spreadsheets.Values.Append(s.spreadsheetID, sheetName, &sheets.ValueRange{
		Values: [][]interface{}{values},
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *GoogleSheetService) getData(ctx context.Context, sheetName string, columnIndex int, value interface{}) (int, []interface{}, error) {
	resp, err := s.client.Spreadsheets.Values.Get(s.spreadsheetID, sheetName).Context(ctx).Do()
	if err != nil {
		return 0, nil, err
	}
	if len(resp.Values) == 0 {
		return 0, nil, nil
	}
	for i := range resp.Values {
		if resp.Values[i][columnIndex] == value {
			return i, resp.Values[i], nil
		}
	}
	return 0, nil, nil
}

func (s *GoogleSheetService) updateData(ctx context.Context, sheetName string, rowIndex int, value []interface{}) error {
	_, err := s.client.Spreadsheets.Values.Update(s.spreadsheetID, fmt.Sprintf("%s!%d:%d", sheetName, rowIndex, rowIndex), &sheets.ValueRange{
		Values: [][]interface{}{value},
	}).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}
