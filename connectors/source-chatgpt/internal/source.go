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
	"io"
	"net/http"
	"sync"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	defaultEventType   = "vanus-chatGPT-type"
	defaultEventSource = "vanus-chatGPT-source"
	headerSource       = "vanus-source"
	headerType         = "vanus-type"
)

var _ cdkgo.HTTPSource = &chatGPTSource{}

func NewChatGPTSource() cdkgo.HTTPSource {
	return &chatGPTSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type chatGPTSource struct {
	config  *chatGPTConfig
	events  chan *cdkgo.Tuple
	number  int
	day     string
	lock    sync.Mutex
	service *chatGPTService
	logger  zerolog.Logger
}

func (s *chatGPTSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*chatGPTConfig)
	s.config.Init()
	s.service = newChatGPTService(s.config, s.logger)
	return nil
}

func (s *chatGPTSource) Name() string {
	return "ChatGPTSource"
}

func (s *chatGPTSource) Destroy() error {
	return nil
}

func (s *chatGPTSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *chatGPTSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	eventSource := req.Header.Get(headerSource)
	if eventSource == "" {
		eventSource = defaultEventSource
	}
	eventType := req.Header.Get(headerType)
	if eventType == "" {
		eventType = defaultEventType
	}
	event := ce.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetType(eventType)
	event.SetSource(eventSource)
	go func(event ce.Event) {
		content, err := s.service.CreateChatCompletion(string(body))
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to get content from ChatGPT")
		}
		event.SetData(ce.ApplicationJSON, map[string]string{
			"content": content,
		})
		s.events <- &cdkgo.Tuple{
			Event: &event,
			Success: func() {
				s.logger.Info().Msg("send event to target success")
			},
			Failed: func(err2 error) {
				s.logger.Warn().Err(err2).Msg("failed to send event to target")
			},
		}
	}(event)
}
