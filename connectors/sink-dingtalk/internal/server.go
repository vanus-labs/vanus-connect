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
	"net/http"
	"sync/atomic"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/sink/dingtalk/internal/bot"
)

const (
	name = "Dingtalk Sink"
)

var (
	errSinkWrongEventNumber = cdkgo.NewResult(http.StatusBadRequest,
		"the event number must be 1")
)

func NewSink() cdkgo.Sink {
	return &dingtalkSink{}
}

var _ cdkgo.Sink = &dingtalkSink{}

type dingtalkSink struct {
	cfg    *dingtalkConfig
	count  int64
	b      *bot.Bot
	logger zerolog.Logger
}

func (f *dingtalkSink) Arrived(_ context.Context, events ...*v2.Event) cdkgo.Result {
	if len(events) != 1 {
		return errSinkWrongEventNumber
	}
	e := events[0]
	atomic.AddInt64(&f.count, int64(len(events)))
	f.logger.Info().Int64("count", atomic.LoadInt64(&f.count)).Str("event_id", e.ID()).Msg("receive a new event")
	if err := f.b.SendMessage(e); err != nil {
		return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
	}
	return cdkgo.SuccessResult
}

func (f *dingtalkSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	config, _ := cfg.(*dingtalkConfig)
	f.logger = log.FromContext(ctx)
	f.cfg = config
	f.b = bot.NewBot(f.logger)
	return f.b.Init(config.Bot)
}

func (f *dingtalkSink) Name() string {
	return "Dingtalk Sink"
}

func (f *dingtalkSink) Destroy() error {
	return nil
}
