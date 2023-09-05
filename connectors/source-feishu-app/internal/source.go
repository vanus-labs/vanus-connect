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
	"io"
	"net/http"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/rs/zerolog"

	cdk "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/connector"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdk.HTTPSource = &FeishuSource{}

type FeishuSource struct {
	logger  zerolog.Logger
	cfg     *Config
	handler *dispatcher.EventDispatcher
	events  chan *cdk.Tuple
}

func NewSource() cdk.HTTPSource {
	return &FeishuSource{
		events: make(chan *cdk.Tuple, 100),
	}
}

func (s *FeishuSource) Initialize(ctx context.Context, cfg cdk.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.cfg = cfg.(*Config)

	return s.start(ctx)
}

func (s *FeishuSource) Name() string {
	return "FeishuAppSource"
}

func (s *FeishuSource) Destroy() error {
	return nil
}

func (s *FeishuSource) Chan() <-chan *connector.Tuple {
	return s.events
}

func (s *FeishuSource) start(_ context.Context) error {
	feishu := NewFeishu(s.logger, s.cfg, s.events)
	s.handler = dispatcher.NewEventDispatcher(s.cfg.VerificationToken, s.cfg.EventEncryptKey).
		OnP2MessageReceiveV1(feishu.OnChatBotMessageReceived)
	s.handler.InitConfig(larkevent.WithLogLevel(larkcore.LogLevelInfo))
	return nil
}

func (s *FeishuSource) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	event, err := s.body2EventReq(ctx, request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		s.logger.Warn().Err(err).Msg("body to event failed")
		return
	}
	eventResp := s.handler.Handle(ctx, event)
	err = s.write(ctx, writer, eventResp)
	if err != nil {
		s.logger.Warn().Err(err).Msg("write resp failed")
	}
}

func (s *FeishuSource) body2EventReq(ctx context.Context, req *http.Request) (*larkevent.EventReq, error) {
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	eventReq := &larkevent.EventReq{
		Header:     req.Header,
		Body:       rawBody,
		RequestURI: req.RequestURI,
	}
	return eventReq, nil
}

func (s *FeishuSource) write(ctx context.Context, writer http.ResponseWriter, eventResp *larkevent.EventResp) error {
	writer.WriteHeader(eventResp.StatusCode)
	for k, vs := range eventResp.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}

	if len(eventResp.Body) > 0 {
		_, err := writer.Write(eventResp.Body)
		return err
	}
	return nil
}
