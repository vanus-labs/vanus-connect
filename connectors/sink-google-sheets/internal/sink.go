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
	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"net/http"
	"sort"
)

const xvSheetName = "xvsheetname"

var (
	errInvalidSheetName = cdkgo.NewResult(http.StatusBadRequest, fmt.Sprintf("extension %s is invalid", xvSheetName))
)

var _ cdkgo.Sink = &GoogleSheetSink{}

func NewGoogleSheetSink() cdkgo.Sink {
	return &GoogleSheetSink{}
}

type GoogleSheetSink struct {
	service          *GoogleSheetService
	defaultSheetName string
	summary          *Summary
}

func (s *GoogleSheetSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	config := cfg.(*GoogleSheetConfig)
	spreadsheetID := config.SheetID
	sheetName := config.SheetName

	//sheetURL := strings.TrimSpace(config.SheetURL)
	//spreadsheetID, err := getSpreadsheetID(sheetURL)
	//if err != nil {
	//	return errors.Wrap(err, "ge spread sheet id error")
	//}
	//sheetID, err := getSheetID(sheetURL)
	//if err != nil {
	//	return errors.Wrap(err, "ge sheet id error")
	//}
	service, err := newGoogleSheetService(spreadsheetID, config.Credentials)
	if err != nil {
		return err
	}
	s.service = service
	s.defaultSheetName = sheetName
	//s.defaultSheetName = service.getSheetName(sheetID)
	//if s.defaultSheetName == "" {
	//	return fmt.Errorf("sheetURL sheetID %d not exist", sheetID)
	//}
	if config.Summary != nil {
		summary, err := newSummary(config.Summary, s.defaultSheetName+"_summary", service)
		if err != nil {
			return err
		}
		s.summary = summary
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
		r := s.saveDataToSpreadsheet(ctx, event)
		if r != cdkgo.SuccessResult {
			return r
		}
	}
	return cdkgo.SuccessResult
}

func (s *GoogleSheetSink) saveDataToSpreadsheet(ctx context.Context, event *ce.Event) cdkgo.Result {
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
	err := s.service.createSheetIfNotExist(ctx, sheetName)
	if err != nil {
		log.Error("create sheet error", map[string]interface{}{
			log.KeyError: err,
			"sheetName":  sheetName,
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "create sheet error")
	}
	// sheet row
	sheetRow := make(map[string]interface{})
	err = event.DataAs(&sheetRow)
	if err != nil {
		log.Error("event data decode error", map[string]interface{}{
			log.KeyError: err,
			"data":       string(event.Data()),
		})
		return cdkgo.NewResult(http.StatusBadRequest, "event data decode error")
	}
	// get sheet headers
	headers, err := s.service.getHeader(sheetName)
	if err != nil {
		if err == headerNotExistErr {
			for k := range sheetRow {
				headers = append(headers, k)
			}
			sort.Strings(headers)
			err = s.service.insertHeader(ctx, sheetName, headers)
			if err != nil {
				log.Error("insert header error", map[string]interface{}{
					log.KeyError: err,
					"sheetName":  sheetName,
				})
				return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
			}
		} else {
			log.Error("get header error", map[string]interface{}{
				log.KeyError: err,
				"sheetName":  sheetName,
			})
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		}
	}
	var values []interface{}
	for _, key := range headers {
		values = append(values, sheetValue(sheetRow[key]))
	}
	// append data
	err = s.service.appendData(ctx, sheetName, values)
	if err != nil {
		log.Error("append sheet data error", map[string]interface{}{
			log.KeyError: err,
			"sheetName":  sheetName,
		})
		return cdkgo.NewResult(http.StatusInternalServerError, "append data error")
	}
	if s.summary != nil {
		err = s.summary.appendData(ctx, sheetRow)
		if err != nil {
			log.Error("append summary data error", map[string]interface{}{
				log.KeyError: err,
			})
		}
	}
	return cdkgo.SuccessResult
}
