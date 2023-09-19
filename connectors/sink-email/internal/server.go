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
	"net/http"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var (
	errFailedToSend = cdkgo.NewResult(http.StatusInternalServerError,
		"email: failed to sent email, please view logs")
)

func NewEmailSink() cdkgo.Sink {
	return &emailSink{}
}

var _ cdkgo.Sink = &emailSink{}

type emailSink struct {
	count    int64
	logger   zerolog.Logger
	emailCfg EmailConfig
}

func (e *emailSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	e.logger = log.FromContext(ctx)
	config := cfg.(*emailConfig)
	e.emailCfg = config.Email
	return nil
}

func (e *emailSink) Arrived(ctx context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		event := events[idx]
		m := &EmailMessage{}
		err := json.Unmarshal(event.Data(), m)
		if err != nil {
			e.logger.Warn().Err(err).Str("event_id", event.ID()).Msg("event data unmarshal error")
			return cdkgo.NewResult(http.StatusBadRequest, "event data unmarshal error")
		}
		err = m.Validate()
		if err != nil {
			e.logger.Warn().Err(err).Str("event_id", event.ID()).Msg("event data invalid")
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		if err := e.send(ctx, m); err != nil {
			e.logger.Warn().Err(err).Str("event_id", event.ID()).
				Str("receiver", m.Receiver).
				Msg("failed to send email")
			return errFailedToSend
		} else {
			e.logger.Info().Str("event_id", event.ID()).
				Str("receiver", m.Receiver).
				Msg("success to send email")
		}
	}
	return cdkgo.SuccessResult
}

func (e *emailSink) Name() string {
	return "emailSink"
}

func (e *emailSink) Destroy() error {
	return nil
}
