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
	"fmt"
	"net/http"

	"github.com/slack-go/slack/slackevents"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/source/chatai/chat"
)

const (
	name        = "Slack Source"
	defaultPort = 8080
)

var _ cdkgo.Source = &slackSource{}

func NewSource() cdkgo.Source {
	return &slackSource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

type slackSource struct {
	config      *slackConfig
	ch          chan *cdkgo.Tuple
	server      *http.Server
	verifyToken slackevents.TokenComparator
	chatService *chat.ChatService
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
	log.Info("event error", map[string]interface{}{
		log.KeyError: err,
	})
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

func (s *slackSource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*slackConfig)
	if s.config.Port == 0 {
		s.config.Port = defaultPort
	}
	s.verifyToken = slackevents.TokenComparator{
		VerificationToken: s.config.VerifyToken,
	}
	if s.config.EnableChatAi {
		s.chatService = chat.NewChatService(*s.config.ChatConfig)
	}
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s,
	}
	go func() {
		log.Info("http server is ready to start", map[string]interface{}{
			"port": s.config.Port,
		})
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("cloud not listen on %d, error:%s", s.config.Port, err.Error()))
		}
		log.Info("http server stopped", nil)
	}()
	return nil
}

func (s *slackSource) Name() string {
	return name
}

func (s *slackSource) Destroy() error {
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	if s.chatService != nil {
		s.chatService.Close()
	}
	return nil
}
