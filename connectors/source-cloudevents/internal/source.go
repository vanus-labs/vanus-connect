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
	"fmt"
	"net/http"
	"sync/atomic"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Source = &cloudEventsSource{}

func NewSource() cdkgo.Source {
	return &cloudEventsSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type cloudEventsSource struct {
	config *cloudEventsConfig
	events chan *cdkgo.Tuple
	count  int64
	server ce.Client
	logger zerolog.Logger
}

func (s *cloudEventsSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*cloudEventsConfig)
	if s.config.Port <= 0 {
		s.config.Port = 8080
	}
	options := []cehttp.Option{ce.WithPort(s.config.Port)}
	if s.config.Path != "" {
		options = append(options, ce.WithPath(s.config.Path))
	}
	if !s.config.Auth.IsEmpty() {
		options = append(options, ce.WithMiddleware(s.handleAuthentication))
	}
	server, err := ce.NewClientHTTP(options...)
	if err != nil {
		return err
	}
	s.server = server
	go func() {
		if err = s.server.StartReceiver(ctx, s.handleEvent); err != nil {
			panic(fmt.Sprintf("start CloudEvents receiver failed: %s", err.Error()))
		}
	}()
	return nil
}

func (s *cloudEventsSource) Name() string {
	return "CloudEvents Source"
}

func (s *cloudEventsSource) Destroy() error {
	return nil
}

func (s *cloudEventsSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *cloudEventsSource) handleEvent(_ context.Context, e event.Event) ce.Result {
	atomic.AddInt64(&s.count, 1)
	s.logger.Info().Int64("total", atomic.LoadInt64(&s.count)).Msg("receive a new event")
	s.events <- &cdkgo.Tuple{Event: &e, Success: func() {
		s.logger.Info().Str("event_id", e.ID()).Msg("send an event success")
	}}
	return ce.ResultACK
}

func (s *cloudEventsSource) handleAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok && s.config.Auth.Username == username && s.config.Auth.Password == password {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="connector", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
