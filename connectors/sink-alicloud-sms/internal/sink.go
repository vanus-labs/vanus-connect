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
	"context"
	"encoding/json"
	"errors"
	"net/http"

	v2 "github.com/cloudevents/sdk-go/v2"

	cdkgo "github.com/vanus-labs/cdk-go"
)

const (
	FieldPhones = "phones"
)

// Config
var _ cdkgo.SinkConfigAccessor = &smsConfig{}

type smsConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	aliConfig        `json:",inline" yaml:",inline"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &smsConfig{}
}

func (fc *smsConfig) Validate() error {
	return fc.SinkConfig.Validate()
}

// Sink
var _ cdkgo.Sink = &smsSink{}

type smsSink struct {
	cfg *smsConfig
	sms *aliSMS
}

func NewSink() cdkgo.Sink {
	return &smsSink{
		sms: new(aliSMS),
	}
}

func (s *smsSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*smsConfig)
	if !ok {
		return errors.New("aliyun sms: invalid configuration type")
	}

	s.cfg = _cfg
	return s.sms.init(_cfg.aliConfig)
}

func (s *smsSink) Name() string {
	return "Aliyun SMS"
}

func (s *smsSink) Destroy() error {
	return nil
}

func (s *smsSink) Arrived(_ context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		e := events[idx]
		var data map[string]interface{}
		err := json.Unmarshal(e.Data(), &data)
		if err != nil {
			return cdkgo.NewResult(http.StatusInternalServerError, "event data unmarshal error")
		}

		phones, ok := data[FieldPhones].(string)
		if !ok {
			return cdkgo.NewResult(http.StatusInternalServerError, "event data not contain phones")
		}

		err = s.sms.sendMsg(phones)
		if err != nil {
			return cdkgo.NewResult(http.StatusInternalServerError, "failed send sms")
		}
	}
	return cdkgo.SuccessResult
}
