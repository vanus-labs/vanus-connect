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
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	fb "github.com/huandu/facebook/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
)

var (
	errFailedToCreate = cdkgo.NewRequest(http.StatusBadRequest, "facebook: failed to create form with data")
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
		event := events[idx]

		var form map[string]interface{}

		err := json.Unmarshal([]byte(event.Data()), &form)
		if err != nil {
			return errFailedToCreate
		}

		access_token, pageId := s.cfg.AccessToken, s.cfg.PageId
		name, follow_up_action_url, questions, context_card_id, legal_content_id := form["name"], form["follow_up_action_url"], form["questions"], form["context_card_id"], form["legal_content_id"]

		res, err := fb.Post(pageId+"leadgen_forms", fb.Params{
			"name":                 name,
			"follow_up_action_url": follow_up_action_url,
			"questions":            questions,
			"context_card_id":      context_card_id,
			"legal_content_id":     legal_content_id,
			"access_token":         access_token,
		})
		if err != nil {
			return nil
		}

		return res["data"]

	}
	return cdkgo.SuccessResult
}

func (s *FacebookLeadAdsSink) Name() string {
	return "FacebookLeadAdsSink"
}

func (s *FacebookLeadAdsSink) Destroy() error {
	// nothing to do here
	return nil
}
