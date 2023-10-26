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
	"fmt"
	"net/http"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name      = "Slack Sink"
	xvChannel = "xvchannel"
	xvMsgType = "xvmsgtype"
)

var (
	errInvalidMessage = cdkgo.NewResult(http.StatusBadRequest,
		"slack: invalid message body")
	errInvalidChannels = cdkgo.NewResult(http.StatusBadRequest,
		"slack: invalid or empty channels")
	errInvalidAppName = cdkgo.NewResult(http.StatusBadRequest,
		"slack: invalid or empty AppName")
	errFailedToSend = cdkgo.NewResult(http.StatusInternalServerError,
		"slack: failed to sent message, please view logs")
)

var _ cdkgo.SinkConfigAccessor = &slackConfig{}

type slackConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	DefaultChannel   string `json:"default_channel" yaml:"default_channel" validate:"required"`
	Token            string `yaml:"token" json:"token" validate:"required"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &slackConfig{}
}

func NewSlackSink() cdkgo.Sink {
	return &slackSink{}
}

var _ cdkgo.Sink = &slackSink{}

type slackSink struct {
	count          int64
	client         *slack.Client
	defaultChannel string
	defaultMsgType string
	logger         zerolog.Logger
}

func (e *slackSink) Arrived(ctx context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		event := events[idx]

		channelID, ok := event.Extensions()[xvChannel].(string)
		if !ok {
			channelID = e.defaultChannel
		}
		m := &Message{}
		if err := json.Unmarshal(event.Data(), m); err != nil {
			e.logger.Error().Err(err).Str("channel", channelID).Str("event_id", event.ID()).Msg("json unmarshal failed")
			return errInvalidMessage
		}
		if err := m.validate(); err != nil {
			e.logger.Error().Err(err).Str("channel", channelID).Str("event_id", event.ID()).Msg("message validate failed")
			return errInvalidMessage
		}
		start := time.Now()
		if err := e.send(ctx, channelID, m); err != nil {
			e.logger.Error().Err(err).Str("channel", channelID).Str("event_id", event.ID()).Msg("failed to send slack")
			return errFailedToSend
		} else if time.Now().Sub(start) > time.Second {
			e.logger.Info().Str("channel", channelID).Str("event_id", event.ID()).
				Interface("used_time", time.Now().Sub(start)).Msg("success to send slack, but takes too long")
		} else {
			e.logger.Info().Str("channel", channelID).Str("event_id", event.ID()).Msg("success to send slack")
		}
	}
	return cdkgo.SuccessResult
}

func (e *slackSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	e.logger = log.FromContext(ctx)
	config := cfg.(*slackConfig)
	e.defaultChannel = config.DefaultChannel
	e.defaultMsgType = "plain_text"
	e.client = slack.New(config.Token)
	_, err := e.client.AuthTest()
	return err
}

func (e *slackSink) Name() string {
	return name
}

func (e *slackSink) Destroy() error {
	// nothing to do
	return nil
}

func (e *slackSink) send(ctx context.Context, channelID string, m *Message) (err error) {
	_, _, err = e.client.PostMessageContext(ctx, channelID, slack.MsgOptionBlocks(m.Blocks.BlockSet...))
	return err
}

type Message struct {
	Blocks *slack.Blocks `json:"blocks"`
}

func (m *Message) validate() error {
	if m.Blocks != nil && len(m.Blocks.BlockSet) > 0 {
		return nil
	}
	return fmt.Errorf("message is invalid")
}
