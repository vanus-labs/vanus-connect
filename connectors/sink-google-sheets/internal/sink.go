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
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
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
	summary          []*Summary
}

func (s *GoogleSheetSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	config := cfg.(*GoogleSheetConfig)
	spreadsheetID := config.SheetID
	sheetName := config.SheetName

	service, err := newGoogleSheetService(spreadsheetID, config.Credentials)
	if err != nil {
		return err
	}
	s.service = service
	s.defaultSheetName = sheetName
	if len(config.Summary) > 0 {
		for i := range config.Summary {
			summary, err := newSummary(config.Summary[i], service)
			if err != nil {
				return err
			}
			s.summary = append(s.summary, summary)
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
			var index int
			headers = make(map[string]int, len(sheetRow))
			for k := range sheetRow {
				headers[k] = index
				index++
			}
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
	total := len(headers)
	var headerChange bool
	for key := range sheetRow {
		_, exist = headers[key]
		if exist {
			continue
		}
		headers[key] = total
		headerChange = true
		total++
	}
	if headerChange {
		err = s.service.updateHeader(ctx, sheetName, headers)
		if err != nil {
			log.Error("insert header error", map[string]interface{}{
				log.KeyError: err,
				"sheetName":  sheetName,
			})
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		}
	}
	values := make([]interface{}, total)
	for key, index := range headers {
		values[index] = sheetValue(sheetRow[key])
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
	if len(s.summary) > 0 {
		for i := range s.summary {
			err = s.summary[i].appendData(ctx, event.Time(), sheetRow)
			if err != nil {
				log.Error("append summary data error", map[string]interface{}{
					log.KeyError: err,
					"index":      i,
				})
			}
		}
	}
	return cdkgo.SuccessResult
}
