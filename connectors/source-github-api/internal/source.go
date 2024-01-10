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
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/go-github/v52/github"
	"github.com/google/uuid"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"go.uber.org/ratelimit"
	"sync"
	"time"
)

var _ cdkgo.Source = &GitHubAPISource{}

type GitHubAPISource struct {
	config  *GitHubAPIConfig
	events  chan *cdkgo.Tuple
	client  *github.Client
	m       sync.Map
	Limiter ratelimit.Limiter

	numRepos   int
	numRecords int
	numPRs     int
}

func Source() cdkgo.Source {
	return &GitHubAPISource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

func (s *GitHubAPISource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*GitHubAPIConfig)
	s.config.Init()
	s.client = github.NewTokenClient(ctx, s.config.GitHubAccessToken)
	s.Limiter = ratelimit.New(s.config.GitHubHourLimit / 3600)

	go s.start(ctx)
	return nil
}

func (s *GitHubAPISource) Name() string {
	return "GitHubAPISource"
}

func (s *GitHubAPISource) Destroy() error {
	return nil
}

func (s *GitHubAPISource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *GitHubAPISource) start(ctx context.Context) {
	log.Info("!!! starting !!!", map[string]interface{}{
		"starting time": time.Now(),
	})

	switch s.config.APIType {
	case PR:
		s.startPR(ctx)
	case Contributor:
		s.startContr(ctx)
	}

	log.Info("!!! ending !!!", map[string]interface{}{
		"ending time": time.Now(),
	})
}

func (s *GitHubAPISource) sendEvent(eventType, org string, data map[string]interface{}) []byte {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetTime(time.Now())
	event.SetType(eventType)
	event.SetSource(fmt.Sprintf("https://github.com/%s", org))
	event.SetData(ce.ApplicationJSON, data)
	s.events <- &cdkgo.Tuple{
		Event: &event,
	}
	return event.Data()
}
