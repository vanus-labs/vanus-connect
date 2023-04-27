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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit#gid=0
func getSheetID(sheetURL string) (int64, error) {
	// get sheet id
	arr := strings.Split(sheetURL, "#")
	if len(arr) < 2 {
		return 0, fmt.Errorf("sheet URL is invalid")
	}
	sheetIDStr := arr[1]
	if len(sheetIDStr) < 5 {
		return 0, fmt.Errorf("sheet URL gid %s is invalid", sheetIDStr)
	}
	sheetID, err := strconv.ParseInt(sheetIDStr[4:], 10, 64)
	if err != nil {
		return 0, err
	}
	return sheetID, nil
}

// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit#gid=0
func getSpreadsheetID(sheetURL string) (string, error) {
	arr := strings.Split(sheetURL[39:], "/")
	return arr[0], nil
}

func sheetValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	switch value.(type) {
	case string, bool, float64:
		return value
	default:
		v, err := json.Marshal(value)
		if err == nil {
			return string(v)
		}
	}
	return nil
}

func convertFloat(value interface{}) (v float64, err error) {
	if value == nil {
		return 0, nil
	}
	switch _value := value.(type) {
	case float64:
		v = _value
	case string:
		if _value == "" {
			v = 0
		} else {
			v, err = strconv.ParseFloat(_value, 64)
			if err != nil {
				return 0, err
			}
		}
	}
	return v, nil
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
