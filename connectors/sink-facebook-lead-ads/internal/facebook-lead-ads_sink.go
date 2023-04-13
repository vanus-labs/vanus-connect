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
	"encoding/json"
	"errors"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.Sink = &FacebookLeadAdsSink{}

func NewFacebookLeadAdsSink() cdkgo.Sink {
	return &FacebookLeadAdsSink{}
}

func (s *FacebookLeadAdsSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*FacebookLeadAdsConfig)
	if !ok {
		return errors.New("invalid configuration type")
	}
	return nil
}

func (s *FacebookLeadAdsSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		r := s.createLeadAdsForm(ctx, event)
		if r != cdkgo.SuccessResult {
			return r
		}
	}
	return cdkgo.SuccessResult
}

func (f *FacebookLeadAdsSink) createLeadAdsForm(ctx context.Context, event ...*ce.Event) cdkgo.Result {
	val, exist := event.Extensions()[xvFormName]
	if exist {
		str, ok := val.(string)
		if !ok {
			return errInvalidFormName
		}
		formName := str
	}
	val, exist := event.Extensions()[xvFollowUpUrl]
	if exist {
		str, ok := val.(string)
		if !ok {
			return errInvalidFollowUpUrl
		}
		followUpUrl := str
	}
	val, exist := event.Extensions()[xvQuestions]
	var allQuestions map[string]interface{}
	if exist {
		err := json.Unmarshal([]byte(val, &allQuestions))
		if err != nil {
			return nil, err
		}
	}
	// TODO
}

func (s *FacebookLeadAdsSink) Name() string {
	return "FacebookLeadAdsSink"
}

func (s *FacebookLeadAdsSink) Destroy() error {
	// nothing to do here
	return nil
}
