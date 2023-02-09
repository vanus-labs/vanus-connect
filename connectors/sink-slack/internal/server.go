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
	"strings"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/nikoksr/notify/service/slack"
	"github.com/pkg/errors"
)

const (
	name            = "Slack Sink"
	xvSlackApp      = "xvslackapp"
	xvSlackChannels = "xvslackchannels"
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

type SlackConfig struct {
	AppName        string `yaml:"app_name" json:"app_name" validate:"required"`
	Token          string `yaml:"token" json:"token" validate:"required"`
	DefaultChannel string `yaml:"default_channel" json:"default_channel"`
}

type slackConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	DefaultAppID     string        ` json:"default" yaml:"default"`
	Slacks           []SlackConfig `json:"slack" yaml:"slack" validate:"required,gt=0,dive"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &slackConfig{}
}

func NewSlackSink() cdkgo.Sink {
	return &slackSink{
		apps: map[string]SlackConfig{},
	}
}

var _ cdkgo.Sink = &slackSink{}

type slackSink struct {
	cfg   *slackConfig
	count int64
	apps  map[string]SlackConfig
}

func (e *slackSink) Arrived(ctx context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		event := events[idx]

		var appID string
		v, exist := event.Extensions()[xvSlackApp]
		if exist {
			str, ok := v.(string)
			if !ok {
				return errInvalidAppName
			}
			appID = str
		} else if e.cfg.DefaultAppID != "" {
			appID = e.cfg.DefaultAppID
		}

		c, ok := e.apps[appID]
		if !ok {
			return errInvalidAppName
		}

		var ids string
		v, exist = event.Extensions()[xvSlackChannels]
		if exist {
			str, ok := v.(string)
			if !ok {
				return errInvalidChannels
			}
			ids = str
		} else if c.DefaultChannel != "" {
			ids = c.DefaultChannel
		}

		channels := strings.Split(ids, ",")
		if len(channels) == 0 {
			return errInvalidChannels
		}

		m := &Message{}
		if err := json.Unmarshal(event.Data(), m); err != nil {
			return errInvalidMessage
		}

		start := time.Now()
		if err := e.send(ctx, c.Token, m, channels...); err != nil {
			log.Error("failed to send slack", map[string]interface{}{
				log.KeyError: err,
				"channels":   channels,
				"event_id":   event.ID(),
			})
			return errFailedToSend
		} else if time.Now().Sub(start) > time.Second {
			log.Debug("success to send slack, but takes too long", map[string]interface{}{
				"channels":  channels,
				"event_id":  event.ID(),
				"used_time": time.Now().Sub(start),
			})
		} else {
			log.Debug("success to send slack", map[string]interface{}{
				"channels": channels,
				"event_id": event.ID(),
			})
		}
	}
	return cdkgo.SuccessResult
}

func (e *slackSink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*slackConfig)
	if !ok {
		return errors.New("slack: invalid configuration type")
	}

	e.cfg = _cfg
	for _, m := range _cfg.Slacks {
		e.apps[m.AppName] = m
	}
	if e.cfg.DefaultAppID != "" {
		if _, exist := e.apps[e.cfg.DefaultAppID]; !exist {
			return errors.New("slack: the default slack config isn't exist")
		}
	}
	return nil
}

func (e *slackSink) Name() string {
	return name
}

func (e *slackSink) Destroy() error {
	// nothing to do
	return nil
}

func (e *slackSink) send(ctx context.Context, token string, msg *Message, channels ...string) error {
	slackService := slack.New(token)
	slackService.AddReceivers(channels...)
	return slackService.Send(ctx, msg.Subject, msg.Message)
}

type Message struct {
	Subject string
	Message string
}
