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

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name        = "Facebook Source"
	defaultPort = 8080
)

var _ cdkgo.Source = &facebookSource{}

func NewSource() cdkgo.Source {
	return &facebookSource{
		ch: make(chan *cdkgo.Tuple, 1024),
	}
}

type facebookSource struct {
	config *facebookConfig
	ch     chan *cdkgo.Tuple
	server *http.Server
}

func (s *facebookSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.verifyToken(w, req)
	case http.MethodPost:
		err := s.event(req)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("received"))
			return
		}
		var code int
		switch err {
		case errInvalidContentTypeHeader, errMissingHubSignatureHeader, errReadPayload, errInvalidPayload:
			code = http.StatusBadRequest
		case errVerificationFailed:
			code = http.StatusForbidden
		default:
			code = http.StatusBadRequest
		}
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("not support request method"))
	}
}

func (s *facebookSource) verifyToken(w http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	mode := query.Get("hub.mode")
	token := query.Get("hub.verify_token")
	challenge := query.Get("hub.challenge")
	if mode == "subscribe" && token == s.config.VerifyToken {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}
	w.WriteHeader(http.StatusForbidden)
}

func (s *facebookSource) Chan() <-chan *cdkgo.Tuple {
	return s.ch
}

func (s *facebookSource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*facebookConfig)
	if s.config.Port == 0 {
		s.config.Port = defaultPort
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

func (s *facebookSource) Name() string {
	return name
}

func (s *facebookSource) Destroy() error {
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}
