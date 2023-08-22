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
	"sync/atomic"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Source = &scheduleSource{}

func NewSource() cdkgo.Source {
	return &scheduleSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type scheduleSource struct {
	events chan *cdkgo.Tuple
	number uint64
	cron   *cron.Cron
	logger zerolog.Logger
}

func (s *scheduleSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	config := cfg.(*scheduleConfig)
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	spec := config.Cron
	timeZone := config.TimeZone
	if timeZone != "" {
		spec = "TZ=" + timeZone + " " + spec
	}
	schedule, err := parser.Parse(spec)
	if err != nil {
		return err
	}
	c := cron.New(cron.WithParser(parser),
		cron.WithLocation(time.UTC),
		cron.WithChain(cron.Recover(cronLog{logger: s.logger})))
	c.Schedule(schedule, cron.FuncJob(s.makeEvent))
	c.Start()
	return nil
}

func (s *scheduleSource) Name() string {
	return "ScheduleSource"
}

func (s *scheduleSource) Destroy() error {
	if s.cron == nil {
		return nil
	}
	s.cron.Stop()
	return nil
}

func (s *scheduleSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *scheduleSource) makeEvent() {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetSource("vanus.ai/schedule")
	event.SetType("schedule")
	event.SetTime(time.Now())
	_ = event.SetData(ce.ApplicationJSON, map[string]interface{}{})
	number := atomic.AddUint64(&s.number, 1)
	s.events <- &cdkgo.Tuple{
		Event: &event,
		Success: func() {
			s.logger.Info().Uint64("number", number).Msg("send new event success")
		}, Failed: func(err error) {
			s.logger.Warn().Uint64("number", number).Err(err).Msg("send new event failed")
		},
	}
}
