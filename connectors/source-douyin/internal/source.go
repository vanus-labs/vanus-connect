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

	"github.com/amorist/douyin"
	"github.com/amorist/douyin/open"
	"github.com/amorist/douyin/open/config"
	"github.com/amorist/douyin/open/oauth"
	"github.com/amorist/douyin/util"
	"go.uber.org/ratelimit"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Source = &DouyinSource{}

type DouyinSource struct {
	config *DouyinConfig
	events chan *cdkgo.Tuple

	douyin *open.API

	Limiter  ratelimit.Limiter
	numVideo int
}

func Source() cdkgo.Source {
	return &DouyinSource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

func (s *DouyinSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*DouyinConfig)
	s.config.Init()

	dy := douyin.New()
	dyScope := oauth.GetAllScope()
	dyCfg := &config.Config{
		ClientKey:    s.config.ClientKey,
		ClientSecret: s.config.ClientSecret,
		Scopes:       dyScope,
		Cache:        util.NewMemCache(),
	}
	s.douyin = dy.GetOpenAPI(dyCfg)
	err := s.douyin.GetOauth().SetAccessToken(&s.config.DouyinToken)
	if err != nil {
		log.Error("Douyin SetAccessToken failed", map[string]interface{}{
			"error": err,
		})
		return err
	}

	s.Limiter = ratelimit.New(s.config.RateHourLimit / 3600)

	s.start(ctx)
	return nil
}

func (s *DouyinSource) Name() string {
	return "DouyinSource"
}

func (s *DouyinSource) Destroy() error {
	return nil
}

func (s *DouyinSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *DouyinSource) start(ctx context.Context) {
	log.Info("!!! starting !!!", map[string]interface{}{
		"starting time": time.Now(),
	})

	go s.syncVideo(ctx)
	go s.syncUser(ctx)
}
