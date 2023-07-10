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

//
//import (
//	"encoding/json"
//	"fmt"
//
//	ce "github.com/cloudevents/sdk-go/v2"
//	"github.com/ohler55/ojg/jp"
//	"github.com/ohler55/ojg/oj"
//	"github.com/pkg/errors"
//	"github.com/vanus-labs/cdk-go/log"
//)
//
//type ConvertConfig struct {
//	Database   string   `json:"database" yaml:"database"`
//	Collection string   `json:"collection" yaml:"collection"`
//	UniqueKey  []string `json:"unique_key" yaml:"unique_key"`
//	UniquePath []string `json:"unique_path" yaml:"unique_path"`
//}
//
//func (c *ConvertConfig) Validate() error {
//	if len(c.UniqueKey) != len(c.UniquePath) {
//		return errors.Errorf("key and path length not same")
//	}
//	if c.Database == "" {
//		return errors.New("database can't be empty")
//	}
//	if c.Collection == "" {
//		return errors.New("collection can't be empty")
//	}
//	return nil
//}
//
//const debeziumOp = "iodebeziumop"
//
//type convertStruct struct {
//	config *ConvertConfig
//}
//
//func newConvert(config *ConvertConfig) *convertStruct {
//	return &convertStruct{
//		config: config,
//	}
//}
//
//func getOp(event *ce.Event) string {
//	return event.Extensions()[debeziumOp].(string)
//}
//
//func (c *convertStruct) convert(event *ce.Event) error {
//	op, ok := event.Extensions()[debeziumOp].(string)
//	if !ok {
//		return errors.Errorf("attribute %s must be string", debeziumOp)
//	}
//	var data map[string]interface{}
//	err := json.Unmarshal(event.Data(), &data)
//	if err != nil {
//		return errors.Wrap(err, "event data unmarshal error")
//	}
//
//	event.SetExtension(mongoSinkDatabase, c.config.Database)
//	event.SetExtension(mongoSinkCollection, c.config.Collection)
//	uniqueValue := make([]interface{}, len(c.config.UniquePath))
//	for i, path := range c.config.UniquePath {
//		v, err := getValue(event.Data(), "$."+path)
//		if err != nil {
//			return err
//		}
//		uniqueValue[i] = v
//	}
//	uniqueMap := make(map[string]interface{})
//	for i, key := range c.config.UniqueKey {
//		uniqueMap[key] = uniqueValue[i]
//	}
//	result := map[string]interface{}{}
//	switch op {
//	case "r", "c":
//		// insert.
//		result["inserts"] = []interface{}{data}
//	case "u":
//		// update.
//		if !ok {
//			return fmt.Errorf("data only support map")
//		}
//		for k := range uniqueMap {
//			// remove unique key form update.
//			delete(data, k)
//		}
//		result["updates"] = []interface{}{
//			map[string]interface{}{
//				"filter": uniqueMap,
//				"update": map[string]interface{}{
//					"$set": data,
//				},
//			},
//		}
//	case "d":
//		// delete.
//		result["deletes"] = []interface{}{
//			map[string]interface{}{
//				"filter": uniqueMap,
//			},
//		}
//	default:
//		return fmt.Errorf("unknown op %s", op)
//	}
//	event.SetData(ce.ApplicationJSON, result)
//	return nil
//}
//
//func getValue(d []byte, path string) (interface{}, error) {
//	p, err := jp.ParseString(path)
//	if err != nil {
//		return nil, err
//	}
//	obj, err := oj.Parse(d)
//	if err != nil {
//		return nil, err
//	}
//	res := p.Get(obj)
//	if len(res) == 0 {
//		return nil, errors.New("not found")
//	} else if len(res) == 1 {
//		return res[0], nil
//	}
//	return res, nil
//}
//
//func (s *mongoSink) convertEvents(events ...*ce.Event) ([]*ce.Event, error) {
//	if s.convertStruct == nil {
//		return events, nil
//	}
//	es := make([]*ce.Event, len(events))
//	for idx := range events {
//		err := s.convertStruct.convert(events[idx])
//		if err != nil {
//			log.Warning("convert event failed", map[string]interface{}{
//				log.KeyError: err,
//				"event":      events[idx].ID(),
//			})
//			return nil, err
//		}
//		es[idx] = events[idx]
//	}
//	return es, nil
//}
