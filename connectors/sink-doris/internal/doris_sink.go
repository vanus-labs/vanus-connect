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
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/linkall-labs/cdk-go/log"
	cdkutil "github.com/linkall-labs/cdk-go/utils"
)

type DorisSink struct {
	streamLoad *StreamLoad
	timeout    time.Duration
	logger     log.Logger
}

func (s *DorisSink) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}
	err := cfg.Validate()
	if err != nil {
		return err
	}
	// init stream load
	s.streamLoad = NewStreamLoad(cfg, s.logger)
	return s.streamLoad.Start()
}

func (s *DorisSink) Name() string {
	return "DorisSink"
}

func (s *DorisSink) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *DorisSink) Destroy() error {
	if s.streamLoad != nil {
		s.streamLoad.Stop()
	}
	return nil
}

func (s *DorisSink) Receive(ctx context.Context, event ce.Event) protocol.Result {
	err := s.streamLoad.WriteEvent(ctx, event)
	if err != nil {
		return err
	}
	return ce.ResultACK
}
