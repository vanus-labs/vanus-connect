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
	"sync/atomic"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/connector"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Sink = &cloudEventsSink{}

func NewSink() cdkgo.Sink {
	return &cloudEventsSink{}
}

type cloudEventsSink struct {
	count    int64
	logger   zerolog.Logger
	cfg      *config
	ceClient ce.Client
}

func (s *cloudEventsSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.cfg = cfg.(*config)
	s.logger = log.FromContext(ctx)
	ceClient, err := ce.NewClientHTTP(ce.WithTarget(s.cfg.Target))
	if err != nil {
		return err
	}
	s.ceClient = ceClient
	return nil
}

func (s *cloudEventsSink) Name() string {
	return "CloudEvents sink"
}

func (s *cloudEventsSink) Destroy() error {
	return nil
}

func (s *cloudEventsSink) Arrived(ctx context.Context, events ...*ce.Event) connector.Result {
	for idx := range events {
		e := events[idx]
		atomic.AddInt64(&s.count, 1)
		s.logger.Info().Int64("total", atomic.LoadInt64(&s.count)).Msg("receive event")
		result := s.ceClient.Send(ctx, *e)
		if ce.IsACK(result) {
			s.logger.Info().Str("id", e.ID()).Msg("send event success")
		} else {
			s.logger.Info().Err(result).Str("id", e.ID()).Msg("send event failed")
		}
	}
	return cdkgo.SuccessResult
}
