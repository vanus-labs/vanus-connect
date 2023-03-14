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
	"time"

	"github.com/pkg/errors"
	"github.com/vanus-labs/cdk-go/log"
)

type SummaryConfig struct {
	SheetName  string   `json:"sheet_name" yaml:"sheet_name"`
	PrimaryKey string   `json:"primary_key" yaml:"primary_key" validate:"required"`
	Columns    []Column `json:"columns" yaml:"columns"`
	GroupBy    GroupBy  `json:"group_by" yaml:"group_by"`
}

type Column struct {
	Name string  `json:"name" yaml:"name"`
	Type CalType `json:"type" yaml:"type"`
}

type CalType string

const (
	Sum     CalType = "sum"
	Replace CalType = "replace"
)

type GroupBy string

const (
	Yearly  GroupBy = "yearly"
	Monthly GroupBy = "monthly"
	Weekly  GroupBy = "weekly"
)

type Summary struct {
	service    *GoogleSheetService
	sheetName  string
	groupBy    GroupBy
	primaryKey string
	columnMap  map[string]CalType
	headers    map[string]int
}

func newSummary(config *SummaryConfig, sheetName string, service *GoogleSheetService) (*Summary, error) {
	s := &Summary{
		service:    service,
		sheetName:  config.SheetName,
		groupBy:    config.GroupBy,
		primaryKey: config.PrimaryKey,
	}
	if s.sheetName == "" {
		s.sheetName = sheetName
	}
	err := s.init(config)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Summary) init(config *SummaryConfig) error {
	s.columnMap = make(map[string]CalType, len(config.Columns)+1)
	s.headers = make(map[string]int, len(config.Columns))
	var headerIndex int
	s.headers[s.primaryKey] = headerIndex
	s.columnMap[s.primaryKey] = Replace
	for _, column := range config.Columns {
		if column.Name == s.primaryKey {
			continue
		}
		headerIndex++
		s.headers[column.Name] = headerIndex
		calType := column.Type
		if calType != Sum {
			calType = Replace
		}
		s.columnMap[column.Name] = calType
	}
	return nil
}

func (s *Summary) appendData(ctx context.Context, eventTime time.Time, data map[string]interface{}) error {
	value, ok := data[s.primaryKey]
	if !ok {
		return errors.New("primary key not exist")
	}
	sheetName := s.getSheetName(eventTime)
	// check sheet
	err := s.service.createSheetIfNotExist(ctx, sheetName)
	if err != nil {
		return err
	}
	// get header
	headers, err := s.getHeader(ctx, sheetName)
	if err != nil {
		return err
	}
	// get data
	rowIndex, rowValues, err := s.service.getData(ctx, sheetName, 0, value)
	if err != nil {
		return err
	}
	if rowValues == nil {
		// no exist insert data
		return s.insertData(ctx, sheetName, headers, data)
	}
	// update data
	for key, index := range headers {
		v, exist := data[key]
		if !exist || v == nil {
			continue
		}
		calType, _ := s.columnMap[key]
		if calType == Sum {
			currFloat, err := convertFloat(rowValues[index])
			if err != nil {
				log.Warning("number sheet value is invalid", map[string]interface{}{
					s.primaryKey: value,
					"column":     key,
					"value":      rowValues[index],
				})
			}
			vFloat, err := convertFloat(v)
			if err != nil {
				log.Warning("number event value is invalid", map[string]interface{}{
					s.primaryKey: value,
					"column":     key,
					"value":      v,
				})
				continue
			}
			rowValues[index] = currFloat + vFloat
		} else {
			rowValues[index] = sheetValue(v)
		}
	}
	return s.service.updateData(ctx, sheetName, rowIndex+1, rowValues)
}

func (s *Summary) insertData(ctx context.Context, sheetName string, headers map[string]int, data map[string]interface{}) error {
	values := make([]interface{}, len(headers))
	for key, index := range headers {
		v, _ := data[key]
		calType, _ := s.columnMap[key]
		if calType == Sum {
			vFloat, err := convertFloat(v)
			if err != nil {
				log.Warning("number event value is invalid", map[string]interface{}{
					s.primaryKey: data[s.primaryKey],
					"column":     key,
					"value":      v,
				})
			}
			values[index] = vFloat
		} else {
			values[index] = sheetValue(v)
		}
	}
	err := s.service.appendData(ctx, sheetName, values)
	if err != nil {
		return err
	}
	return nil
}

func (s *Summary) getSheetName(eventTime time.Time) string {
	if s.groupBy == "" {
		return s.sheetName
	}
	if eventTime.IsZero() {
		eventTime = time.Now()
	}
	switch s.groupBy {
	case Yearly:
		return fmt.Sprintf("%s_%d", s.sheetName, eventTime.Year())
	case Monthly:
		return fmt.Sprintf("%s_%s", s.sheetName, eventTime.Format("2006_01"))
	case Weekly:
		year, week := eventTime.ISOWeek()
		return fmt.Sprintf("%s_%d_%d", s.sheetName, year, week)
	default:
		return s.sheetName
	}
}

func (s *Summary) getHeader(ctx context.Context, sheetName string) (map[string]int, error) {
	headers, err := s.service.getHeader(sheetName)
	if err == nil {
		return headers, nil
	}
	if err != nil {
		if err == headerNotExistErr {
			// insert header
			err = s.service.insertHeader(ctx, sheetName, s.headers)
			if err != nil {
				return nil, err
			}
			return s.headers, nil
		}
		return nil, err
	}
	return headers, nil
}
