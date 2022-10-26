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
	"strings"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/tidwall/gjson"
)

type PrimaryKeyType uint8

const (
	None PrimaryKeyType = iota
	EventAttribute
	EventData
)

type PrimaryKey interface {
	Name() string
	Type() PrimaryKeyType
	Value(event ce.Event) string
}

type none struct {
}

func (p none) Type() PrimaryKeyType {
	return None
}

func (p none) Name() string {
	return ""
}

func (p none) Value(event ce.Event) string {
	return ""
}

type eventAttribute struct {
	attr string
}

func (p eventAttribute) Type() PrimaryKeyType {
	return EventAttribute
}

func (p eventAttribute) Name() string {
	return p.attr
}

func (p eventAttribute) Value(event ce.Event) string {
	switch p.attr {
	case "id":
		return event.ID()
	default:
		extensions := event.Context.AsV1().Extensions
		if len(extensions) == 0 {
			return ""
		}
		v, exist := extensions[p.attr]
		if !exist {
			return ""
		}
		value, ok := v.(string)
		if ok {
			return value
		}
		return fmt.Sprintf("%v", v)
	}
}

type eventData struct {
	path string
}

func (p eventData) Type() PrimaryKeyType {
	return EventData
}

func (p eventData) Name() string {
	return p.path
}

func (p eventData) Value(event ce.Event) string {
	result := gjson.GetBytes(event.Data(), p.path)
	return result.String()
}

func GetPrimaryKey(primaryKey string) PrimaryKey {
	primaryKey = strings.TrimSpace(primaryKey)
	if primaryKey == "" {
		return none{}
	}
	index := strings.Index(primaryKey, ".")
	if index == -1 {
		return eventAttribute{
			attr: strings.ToLower(primaryKey),
		}
	}
	// format is data.key remove data. only has key
	return eventData{
		path: primaryKey[5:],
	}
}
