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

	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Source = &slackSource{}

func NewSource() cdkgo.Source {
	return &slackSource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

type slackSource struct {
	config *slackConfig
	events chan *cdkgo.Tuple
	logger zerolog.Logger
	slack  *Slack
}

func (s *slackSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *slackSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*slackConfig)
	s.slack = NewSlack(s.config, s.logger, s.events)
	return s.slack.Start()
}

func (s *slackSource) Name() string {
	return "SlackAppSource "
}

func (s *slackSource) Destroy() error {
	if s.slack != nil {
		s.slack.Stop()
	}
	return nil
}
