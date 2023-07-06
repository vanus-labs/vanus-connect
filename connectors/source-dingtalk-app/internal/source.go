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
	"github.com/google/uuid"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Source = &DingtalkSource{}

type DingtalkSource struct {
	cfg    *Config
	logger zerolog.Logger
	events chan *cdkgo.Tuple
}

func Source() cdkgo.Source {
	return &DingtalkSource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

func (s *DingtalkSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.cfg = cfg.(*Config)
	s.cfg.Init()

	go s.start(ctx)
	return nil
}

func (s *DingtalkSource) Name() string {
	return "Dingtalk Source"
}

func (s *DingtalkSource) Destroy() error {
	return nil
}

func (s *DingtalkSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *DingtalkSource) start(ctx context.Context) {
	ding := NewDingtalk(s)
	cli := client.NewStreamClient(client.WithAppCredential(client.NewAppCredentialConfig(s.cfg.DingtalkAppKey, s.cfg.DingtalkAppSecret)))
	cli.RegisterChatBotCallbackRouter(ding.OnChatBotMessageReceived)

	err := cli.Start(ctx)
	if err != nil {
		s.logger.Fatal().Err(err)
	}
	defer cli.Close()
}

func (s *DingtalkSource) sendEvent(data map[string]interface{}) []byte {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetTime(time.Now())
	event.SetType("Conversion")
	event.SetSource(s.Name())
	event.SetData(ce.ApplicationJSON, data)
	s.events <- &cdkgo.Tuple{
		Event: &event,
	}
	return event.Data()
}
