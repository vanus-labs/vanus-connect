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
)

type SummaryConfig struct {
	PrimaryKey string   `json:"primary_key" yaml:"primary_key" validate:"required"`
	Columns    []Column `json:"columns" yaml:"columns"`
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

type Summary struct {
	service    *GoogleSheetService
	sheetName  string
	primaryKey string
	columnMap  map[string]CalType
	headers    []string
}

func newSummary(config *SummaryConfig, sheetName string, service *GoogleSheetService) (*Summary, error) {
	s := &Summary{
		service:    service,
		sheetName:  sheetName,
		primaryKey: config.PrimaryKey,
	}
	err := s.init(config)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Summary) init(config *SummaryConfig) error {
	s.columnMap = make(map[string]CalType, len(config.Columns)+1)
	s.headers = make([]string, 0, len(config.Columns)+1)
	s.headers = append(s.headers, s.primaryKey)
	s.columnMap[s.primaryKey] = Replace
	for _, column := range config.Columns {
		if column.Name == s.primaryKey {
			continue
		}
		s.headers = append(s.headers, column.Name)
		calType := column.Type
		if calType != Sum {
			calType = Replace
		}
		s.columnMap[column.Name] = calType
	}
	// create sheet
	err := s.service.createSheetIfNotExist(context.Background(), s.sheetName)
	if err != nil {
		return err
	}
	// get header
	headers, err := s.service.getHeader(s.sheetName)
	if err != nil {
		if err == headerNotExistErr {
			// insert header
			err = s.service.insertHeader(context.Background(), s.sheetName, s.headers)
			if err != nil {
				return errors.Wrap(err, "insert summary header error")
			}
		} else {
			return err
		}
	} else {
		// check headers
		for i, v := range headers {
			if _, ok := s.columnMap[v]; !ok {
				if v == s.primaryKey {
					continue
				}
				return fmt.Errorf("sheet column %d name %s exist but config culumn not eixst", i, v)
			}
		}
		s.headers = headers
	}
	return nil
}

func (s *Summary) appendData(ctx context.Context, data map[string]interface{}) error {
	value, ok := data[s.primaryKey]
	if !ok {
		return errors.New("primary key not exist")
	}
	// get data
	rowIndex, rowValues, err := s.service.getData(ctx, s.sheetName, 0, value)
	if err != nil {
		return err
	}
	if rowValues == nil {
		// no exist, need insert data
		var values []interface{}
		for _, key := range s.headers {
			v, _ := data[key]
			calType, _ := s.columnMap[key]
			if calType == Sum {
				vFloat, err := convertFloat(v)
				if err != nil {
					return fmt.Errorf("column %s must be number, but event value is %v", key, v)
				}
				values = append(values, vFloat)
			} else {
				values = append(values, sheetValue(v))
			}
		}
		err = s.service.appendData(ctx, s.sheetName, values)
		if err != nil {
			return err
		}
		return nil
	}
	// update
	for i, key := range s.headers {
		v, ok := data[key]
		if !ok || v == nil {
			continue
		}
		calType, _ := s.columnMap[key]
		if calType == Sum {
			currFloat, err := convertFloat(rowValues[i])
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("column %s must be number,but sheet value is %v", key, rowValues[i]))
			}
			vFloat, err := convertFloat(v)
			if err != nil {
				return fmt.Errorf("column %s must be number, but event value is %v", key, v)
			}
			rowValues[i] = currFloat + vFloat
		} else {
			rowValues[i] = sheetValue(v)
		}
	}
	return s.service.updateData(ctx, s.sheetName, rowIndex+1, rowValues)
}
