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
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Sink = &outlookSink{}

func NewSink() cdkgo.Sink {
	return &outlookSink{}
}

type outlookSink struct {
	service *gmailService
	logger  zerolog.Logger
}

func (s *outlookSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	config := cfg.(*gmailConfig)
	service, err := NewOutlookService(config.OAuth)
	if err != nil {
		return err
	}
	s.service = service
	return nil
}

func (s *outlookSink) Name() string {
	return "outlookSink"
}

func (s *outlookSink) Destroy() error {
	return nil
}

func (s *outlookSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		em := &EmailMessage{}
		err := json.Unmarshal(event.Data(), em)
		if err != nil {
			s.logger.Warn().Err(err).Str("event_id", event.ID()).Msg("event data unmarshal error")
			return cdkgo.NewResult(http.StatusBadRequest, "event data unmarshal error")
		}
		err = em.Validate()
		if err != nil {
			s.logger.Warn().Err(err).Str("event_id", event.ID()).Msg("event data invalid")
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		err = s.service.Send(context.Background(), em)
		if err != nil {
			s.logger.Warn().Err(err).Str("event_id", event.ID()).
				Str("receiver", em.To).
				Msg("failed to send email")
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		} else {
			s.logger.Info().Str("event_id", event.ID()).
				Str("receiver", em.To).
				Msg("success to send email")
		}
	}
	return cdkgo.SuccessResult
}
