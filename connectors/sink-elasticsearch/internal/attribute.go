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
	"github.com/pkg/errors"
	"strconv"
)

const (
	attributePrefix = "xv"
	attributeOp     = attributePrefix + "op"
	attributeIndex  = attributePrefix + "indexname"
	attributeId     = attributePrefix + "id"
)

func getAttr(extensions map[string]interface{}, key string) (string, error) {
	val, exist := extensions[key]
	if !exist {
		return "", errors.Errorf("attribute %s not found", key)
	}
	str, ok := val.(string)
	if !ok {
		return "", errors.Errorf("invalid attribute %s=%v", key, val)
	}
	return str, nil
}

type action string

const (
	actionIndex  action = "index"
	actionUpdate action = "update"
	actionDelete action = "delete"
)

func (s *elasticsearchSink) getAction(extensions map[string]interface{}) (action, error) {
	if len(extensions) == 0 {
		return s.action, nil
	}
	val, exist := extensions[attributeOp]
	if !exist {
		return s.action, nil
	}
	op, ok := val.(string)
	if !ok {
		return "", errors.Errorf("invalid attribute %s=%v", attributeOp, val)
	}
	switch op {
	case "c":
		return s.action, nil
	case "u":
		return actionUpdate, nil
	case "d":
		return actionDelete, nil
	default:
		return "", errors.Errorf("invalid attribute value %s=%s", attributeOp, op)
	}
}

func (s *elasticsearchSink) getDocumentId(extensions map[string]interface{}) (string, error) {
	val, exist := extensions[attributeId]
	if !exist {
		return "", nil
	}
	str, ok := val.(string)
	if ok {
		return str, nil
	}
	intV, ok := val.(int32)
	if ok {
		return strconv.Itoa(int(intV)), nil
	}
	return "", errors.Errorf("invalid attribute %s=%v", attributeId, val)
}

func (s *elasticsearchSink) getIndexName(extensions map[string]interface{}) (string, error) {
	if len(extensions) == 0 {
		return s.config.Secret.IndexName, nil
	}
	val, exist := extensions[attributeIndex]
	if !exist {
		return s.config.Secret.IndexName, nil
	}
	str, ok := val.(string)
	if !ok {
		return "", errors.Errorf("invalid attribute %s=%v", attributeIndex, val)
	}
	return str, nil
}
