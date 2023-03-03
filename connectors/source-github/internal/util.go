// Copyright 2022 Linkall Inc.
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
	"time"
)

const (
	GHHeaderEvent        = "X-GitHub-Event"
	GHHeaderDelivery     = "X-GitHub-Delivery"
	GHHeaderSignature256 = "X-Hub-Signature-256"
	HeaderContentType    = "Content-Type"
)

func getTime(v interface{}) time.Time {
	str, ok := v.(string)
	if !ok {
		return time.Now()
	}
	value, err := time.Parse(time.RFC3339, str)
	if err == nil {
		return value
	}
	return time.Now()
}

func getTimeByTimestamp(v interface{}) time.Time {
	value, ok := v.(float64)
	if ok {
		return time.Unix(int64(value), 0)
	}
	return time.Now()
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	str, ok := v.(string)
	if ok {
		return str
	}
	return fmt.Sprintf("%v", v)
}
