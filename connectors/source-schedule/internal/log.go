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

	"github.com/vanus-labs/cdk-go/log"
)

type cronLog struct{}

func (l cronLog) Info(msg string, keysAndValues ...interface{}) {
	log.Info(msg, l.convertKeysAndValues(keysAndValues))
}

func (l cronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	m := l.convertKeysAndValues(keysAndValues)
	m[log.KeyError] = err
	log.Error(msg, m)
}

func (l cronLog) convertKeysAndValues(keysAndValues ...interface{}) map[string]interface{} {
	len := len(keysAndValues)
	m := make(map[string]interface{}, len)
	for i := 0; i < len; i += 2 {
		var v interface{}
		if i+1 < len {
			v = keysAndValues[i+1]
		}
		m[fmt.Sprintf("%v", keysAndValues[i])] = v
	}
	return m
}
