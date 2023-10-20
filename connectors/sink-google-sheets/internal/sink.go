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
	"sort"
	"sync"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"

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
	buffer           *BufferWriter
	defaultSheetName string
	summary          []*Summary
	lock             sync.Mutex
	logger           zerolog.Logger
}

func (s *GoogleSheetSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	config := cfg.(*GoogleSheetConfig)
	spreadsheetID := config.SheetID
	sheetName := config.SheetName
	s.defaultSheetName = sheetName
	if config.FlushInterval <= 0 {
		config.FlushInterval = 5
	}
	if config.FlushSize <= 0 {
		config.FlushSize = 500
	}
	service, err := newGoogleSheetService(spreadsheetID, config.Credentials, config.OAuth, s.logger)
	if err != nil {
		return err
	}
	s.service = service
	s.buffer = newBufferWriter(service, time.Duration(config.FlushInterval)*time.Second, config.FlushSize)
	s.buffer.Start(ctx)
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
	if s.buffer != nil {
		s.buffer.Stop()
	}
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
	s.lock.Lock()
	defer s.lock.Unlock()
	sheetName, result := s.checkSheetName(ctx, event)
	if result != cdkgo.SuccessResult {
		return result
	}
	// sheet row
	sheetRow := make(map[string]interface{})
	err := event.DataAs(&sheetRow)
	if err != nil {
		s.logger.Error().Err(err).Msg("event data decode error")
		return cdkgo.NewResult(http.StatusBadRequest, "event data decode error")
	}
	// get sheet headers
	headers, result := s.checkHeader(ctx, sheetName, sheetRow)
	if result != cdkgo.SuccessResult {
		return result
	}
	values := make([]interface{}, len(headers))
	for key, index := range headers {
		v, err := sheetValue(sheetRow[key])
		if err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, "sheet value invalid:"+err.Error())
		}
		values[index] = v
	}
	// append data
	err = s.buffer.AppendData(sheetName, values)
	if err != nil {
		s.logger.Error().Err(err).
			Str("sheet_name", sheetName).
			Msg("append sheet data error")
		return cdkgo.NewResult(http.StatusBadRequest, "append data error")
	}
	if len(s.summary) > 0 {
		for i := range s.summary {
			err = s.summary[i].appendData(ctx, event.Time(), sheetRow)
			if err != nil {
				s.logger.Error().Err(err).
					Int("index", i).
					Msg("append summary data error")
			}
		}
	}
	return cdkgo.SuccessResult
}

func (s *GoogleSheetSink) checkSheetName(ctx context.Context, event *ce.Event) (string, cdkgo.Result) {
	var sheetName string

	v, exist := event.Extensions()[xvSheetName]
	if exist {
		str, ok := v.(string)
		if !ok {
			return "", errInvalidSheetName
		}
		sheetName = str
	} else {
		sheetName = s.defaultSheetName
	}
	err := s.service.createSheetIfNotExist(ctx, sheetName)
	if err != nil {
		s.logger.Error().Err(err).
			Str("sheet_name", sheetName).
			Msg("create sheet error")
		return "", cdkgo.NewResult(http.StatusBadRequest, "create sheet error")
	}
	return sheetName, cdkgo.SuccessResult
}

func (s *GoogleSheetSink) checkHeader(ctx context.Context, sheetName string, sheetRow map[string]interface{}) (map[string]int, cdkgo.Result) {
	// get sheet headers
	headers, err := s.service.getHeader(sheetName)
	if err != nil {
		if err == headerNotExistErr {
			headerArr := mapKeys(sheetRow)
			sort.Strings(headerArr)
			headers = make(map[string]int, len(headerArr))
			for index, key := range headerArr {
				headers[key] = index
			}
			err = s.service.insertHeader(ctx, sheetName, headers)
			if err != nil {
				s.logger.Error().Err(err).
					Str("sheet_name", sheetName).
					Msg("insert header error")
				return nil, cdkgo.NewResult(http.StatusBadRequest, err.Error())
			}
		} else {
			s.logger.Error().Err(err).
				Str("sheet_name", sheetName).
				Msg("get header error")
			return nil, cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
	}
	total := len(headers)
	var headerChange bool
	for key := range sheetRow {
		_, exist := headers[key]
		if exist {
			continue
		}
		headers[key] = total
		headerChange = true
		total++
	}
	if headerChange {
		err = s.buffer.FlushSheet(sheetName)
		if err != nil {
			s.logger.Error().Err(err).
				Str("sheet_name", sheetName).
				Msg("flush data error")
			return nil, cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		err = s.service.updateHeader(ctx, sheetName, headers)
		if err != nil {
			s.logger.Error().Err(err).
				Str("sheet_name", sheetName).
				Msg("update header error")
			return nil, cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
	}
	return headers, cdkgo.SuccessResult
}
