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
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/api/sheets/v4"
)

// https://docs.google.com/spreadsheets/d/1tZJPUCOiiR0liRsNtLKhCoQR-Cb8_oPVGMU0kvnabcd/edit#gid=0
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

// https://docs.google.com/spreadsheets/d/1tZJPUCOiiR0liRsNtLKhCoQR-Cb8_oPVGMU0kvnabcd/edit#gid=0
func getSpreadsheetID(sheetURL string) (string, error) {
	arr := strings.Split(sheetURL[39:], "/")
	return arr[0], nil
}

func sheetValue(value interface{}) *sheets.ExtendedValue {
	if value == nil {
		return &sheets.ExtendedValue{
			StringValue: ptrString(""),
		}
	}
	switch v := value.(type) {
	case bool:
		return &sheets.ExtendedValue{
			BoolValue: &v,
		}
	case float64:
		return &sheets.ExtendedValue{
			NumberValue: &v,
		}
	case string:
		return &sheets.ExtendedValue{
			StringValue: &v,
		}
	default:
		return &sheets.ExtendedValue{
			StringValue: ptrString(fmt.Sprintf("%v", v)),
		}
	}
}

func ptrString(str string) *string {
	return &str
}
