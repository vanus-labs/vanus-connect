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
	"time"

	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
)

type GitHubSource struct {
	config *GitHubConfig
	events chan *cdkgo.Tuple
	server *http.Server
}

func Source() cdkgo.Source {
	return &GitHubSource{
		events: make(chan *cdkgo.Tuple, 10),
	}
}

func (s *GitHubSource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*GitHubConfig)
	s.start()
	return nil
}

func (s *GitHubSource) Name() string {
	return "GitHubSource"
}

func (s *GitHubSource) Destroy() error {
	if s.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error("shutdown the server error", map[string]interface{}{
			log.KeyError: err,
		})
	}
	close(s.events)
	return nil
}

func (s *GitHubSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *GitHubSource) start() {
	_handler := newHandler(s.events, s.config.GitHub)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: _handler,
	}
	go func() {
		log.Info("server is ready to start", map[string]interface{}{
			"port": s.config.Port,
		})
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("cloud not listen on %d, error:%s", s.config.Port, err.Error()))
		}
		log.Info("server stopped", nil)
	}()
}
