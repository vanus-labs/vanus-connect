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

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

type GitHubSource struct {
	config  *GitHubConfig
	events  chan *cdkgo.Tuple
	handler *handler
}

func Source() cdkgo.HTTPSource {
	return &GitHubSource{
		events: make(chan *cdkgo.Tuple, 10),
	}
}

func (s *GitHubSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	logger := log.FromContext(ctx)
	s.config = cfg.(*GitHubConfig)
	s.handler = newHandler(s.events, s.config.GitHub, logger)
	return nil
}

func (s *GitHubSource) Name() string {
	return "GitHubSource"
}

func (s *GitHubSource) Destroy() error {
	return nil
}

func (s *GitHubSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *GitHubSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := s.handler.handle(req)
	if err == nil {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("accepted"))
		return
	}
	var code int
	switch err {
	case errPingEvent:
		code = http.StatusOK
	case errInvalidHTTPMethod, errInvalidContentTypeHeader, errMissingGithubEventHeader, errMissingHubDeliveryHeader, errMissingHubSignatureHeader, errReadPayload, errInvalidBody:
		code = http.StatusBadRequest
	case errVerificationFailed:
		code = http.StatusForbidden
	default:
		code = http.StatusInternalServerError
	}
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}
