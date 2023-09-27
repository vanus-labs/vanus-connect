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

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/hanzoai/gochimp3"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Sink = &mailchimpSink{}

func NewMailchimpSink() cdkgo.Sink {
	return &mailchimpSink{
		listMap: map[string]*gochimp3.ListResponse{},
	}
}

type mailchimpSink struct {
	api     *gochimp3.API
	listMap map[string]*gochimp3.ListResponse
	logger  zerolog.Logger
	config  *sinkConfig
}

func (s *mailchimpSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*sinkConfig)
	s.api = gochimp3.New(s.config.ApiKey)
	s.api.User = "vanus-sink-mailchimp"
	return nil
}

func (s *mailchimpSink) Name() string {
	return "Mailchimp Sink"
}

func (s *mailchimpSink) Destroy() error {
	return nil
}

func (s *mailchimpSink) Arrived(_ context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		list, err := s.getList(event)
		if err != nil {
			s.logger.Warn().Err(err).Str("event_id", event.ID()).Msg("get list error")
			continue
		}
		eventAction := s.getEventAction(event)
		switch eventAction {
		case MemberAddAction:
			err = s.addMember(list, event)
		case MemberPutAction:
			err = s.putMember(list, event)
		case MemberUpdateAction:
			err = s.updateMember(list, event)
		case MemberArchiveAction:
			err = s.archiveMember(list, event)
		case MemberDeleteAction:
			err = s.deleteMember(list, event)
		default:
			s.logger.Info().Str("event_id", event.ID()).Str("type", eventAction).Msg("unknown event action")
			continue
		}
		if err != nil {
			s.logger.Warn().Err(err).Str("event_id", event.ID()).Str("type", eventAction).Msg("failed")
		} else {
			s.logger.Info().Str("event_id", event.ID()).Str("type", eventAction).Msg("success")
		}
	}
	return cdkgo.SuccessResult
}

func (s *mailchimpSink) getEventAction(event *ce.Event) string {
	action, exist := event.Extensions()[AttrEventAction].(string)
	if !exist {
		return ""
	}
	return action
}
