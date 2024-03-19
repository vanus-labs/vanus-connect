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
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/keighl/mandrill"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Sink = &mandrillSink{}

func NewMandrillSink() cdkgo.Sink {
	return &mandrillSink{}
}

type mandrillSink struct {
	cli    *mandrill.Client
	logger zerolog.Logger
	config *sinkConfig
}

func (s *mandrillSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*sinkConfig)
	s.cli = mandrill.ClientWithKey(s.config.ApiKey)
	return nil
}

func (s *mandrillSink) Name() string {
	return "Mandrill Sink"
}

func (s *mandrillSink) Destroy() error {
	return nil
}

func (s *mandrillSink) Arrived(_ context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		eventAction := s.getEventAction(event)
		message := &mandrill.Message{}
		err := json.Unmarshal(event.Data(), message)
		if err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, "data unmarshal error:%s"+err.Error())
		}
		switch eventAction {
		case MessageSendTempAction:
			_, err = s.cli.MessagesSendTemplate(message, s.getEventTemplate(event), nil)
		case MessageSendAction:
			_, err = s.cli.MessagesSend(message)
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

func (s *mandrillSink) getEventAction(event *ce.Event) string {
	action, exist := event.Extensions()[AttrEventAction].(string)
	if !exist {
		return MessageSendTempAction
	}
	return action
}

func (s *mandrillSink) getEventTemplate(event *ce.Event) string {
	tName, exist := event.Extensions()[AttrTemplateName].(string)
	if !exist {
		return s.config.TemplateName
	}
	return tName
}
