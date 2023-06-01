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

	"github.com/rs/zerolog"
	"github.com/slack-go/slack/slackevents"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/source/chatai/chat"
)

const (
	name        = "Slack Source"
	defaultPort = 8080
)

var _ cdkgo.HTTPSource = &slackSource{}

func NewSource() cdkgo.HTTPSource {
	return &slackSource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

type slackSource struct {
	config      *slackConfig
	ch          chan *cdkgo.Tuple
	verifyToken slackevents.TokenComparator
	chatService *chat.ChatService
	logger      zerolog.Logger
}

func (s *slackSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not support request method"))
		return
	}
	err := s.event(w, req)
	if err == nil {
		return
	}
	s.logger.Info().Err(err).Msg("event error")
	var code int
	switch err {
	case errVerificationFailed, errVerificationTokenFailed:
		code = http.StatusForbidden
	default:
		code = http.StatusBadRequest
	}
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func (s *slackSource) Chan() <-chan *cdkgo.Tuple {
	return s.ch
}

func (s *slackSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*slackConfig)
	if s.config.Port == 0 {
		s.config.Port = defaultPort
	}
	s.verifyToken = slackevents.TokenComparator{
		VerificationToken: s.config.VerifyToken,
	}
	if s.config.EnableChatAi {
		s.chatService = chat.NewChatService(*s.config.ChatConfig, s.logger)
	}
	return nil
}

func (s *slackSource) Name() string {
	return name
}

func (s *slackSource) Destroy() error {
	if s.chatService != nil {
		s.chatService.Close()
	}
	return nil
}
