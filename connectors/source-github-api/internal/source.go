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
	lodash "github.com/samber/lo"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"go.uber.org/ratelimit"
	"sync"
	"time"
)

var _ cdkgo.Source = &GitHubAPISource{}

type GitHubAPISource struct {
	config     *GitHubAPIConfig
	events     chan *cdkgo.Tuple
	client     *github.Client
	m          sync.Map
	numRepos   int
	numRecords int
	Limiter    ratelimit.Limiter
}

func Source() cdkgo.Source {
	return &GitHubAPISource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

func (s *GitHubAPISource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*GitHubAPIConfig)
	s.client = github.NewTokenClient(ctx, s.config.GitHubAccessToken)
	if s.config.GitHubHourLimit == 0 {
		s.config.GitHubHourLimit = 5000
	}
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
	// Repository
	listOption := &github.RepositoryListByOrgOptions{
		Type: "sources",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 250,
		},
	}
	log.Info("start", map[string]interface{}{
		"time": time.Now(),
	})
	for {
		s.Limiter.Take()
		repos, resp, err := s.client.Repositories.ListByOrg(ctx, s.config.OrgName, listOption)
		if err != nil {
			log.Warning("Repositories.ListByOrg error", map[string]interface{}{
				log.KeyError: err,
			})
		}
		if len(repos) == 0 {
			break
		}

		log.Info("ListByOrg", map[string]interface{}{
			"Current Page": listOption.ListOptions.Page,
			"Next Page":    resp.NextPage,
			"GitHub Rate":  resp.Rate,
		})

		for _, repo := range repos {
			s.listContributors(ctx, repo)
			s.numRepos += 1
		}

		log.Info("stats", map[string]interface{}{
			"numRecords": s.numRecords,
			"numRepos":   s.numRepos,
			"page":       listOption.ListOptions.Page,
		})

		if resp.NextPage <= listOption.ListOptions.Page {
			break
		}
		listOption.ListOptions.Page = resp.NextPage
	}
	log.Info("end", map[string]interface{}{
		"time":       time.Now(),
		"numRecords": s.numRecords,
		"numRepos":   s.numRepos,
	})
}

func (s *GitHubAPISource) listContributors(ctx context.Context, repo *github.Repository) {
	// Contributors
	listOption := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 250,
		},
	}
	for {
		s.Limiter.Take()
		contributors, resp, err := s.client.Repositories.ListContributors(ctx, *repo.Owner.Login, *repo.Name, listOption)
		if err != nil {
			log.Warning("Repositories.ListContributors error", map[string]interface{}{
				log.KeyError: err,
			})
		}
		if len(contributors) == 0 {
			break
		}

		log.Info("ListContributors", map[string]interface{}{
			"Current Page": listOption.ListOptions.Page,
			"Next Page":    resp.NextPage,
			"Project":      *repo.Name,
		})

		s.numRecords += len(contributors)
		for _, contributor := range contributors {
			s.userInfo(ctx, contributor, repo)
		}

		if resp.NextPage <= listOption.ListOptions.Page {
			break
		}
		listOption.ListOptions.Page = resp.NextPage
	}
}

func (s *GitHubAPISource) userInfo(ctx context.Context, contributor *github.Contributor, repo *github.Repository) {
	user := new(github.User)
	v, ok := s.m.Load(*contributor.Login)
	if !ok {
		s.Limiter.Take()
		user0, _, err := s.client.Users.Get(ctx, *contributor.Login)
		if err != nil {
			log.Warning("Users.Get error", map[string]interface{}{
				log.KeyError: err,
			})
			return
		}
		user = user0
		s.m.Store(*contributor.Login, user)
	} else {
		user = v.(*github.User)
	}

	data := make(map[string]interface{})
	data["repo"] = lodash.TernaryF(repo.Name != nil, func() string { return *repo.Name }, func() string { return "" })
	data["org"] = lodash.TernaryF(repo.Owner.Login != nil, func() string { return *repo.Owner.Login }, func() string { return "" })
	data["url"] = lodash.TernaryF(repo.HTMLURL != nil, func() string { return *repo.HTMLURL }, func() string { return "" })
	data["uid"] = lodash.TernaryF(user.Login != nil, func() string { return *user.Login }, func() string { return "" })
	data["username"] = lodash.TernaryF(user.Name != nil, func() string { return *user.Name }, func() string { return "" })
	data["email"] = lodash.TernaryF(user.Email != nil, func() string { return *user.Email }, func() string { return "" })
	data["company"] = lodash.TernaryF(user.Company != nil, func() string { return *user.Company }, func() string { return "" })

	s.sendEvent("contributors", data)
}

func (s *GitHubAPISource) sendEvent(eventType string, data map[string]interface{}) []byte {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetTime(time.Now())
	event.SetType(eventType)
	event.SetSource(fmt.Sprintf("https://github.com/%s", data["org"]))
	event.SetData(ce.ApplicationJSON, data)
	s.events <- &cdkgo.Tuple{
		Event: &event,
	}
	return event.Data()
}

func (s *GitHubAPISource) needWait(rate github.Rate) bool {
	if rate.Remaining > 0 {
		return false
	}
	log.Info("needWait", map[string]interface{}{
		"start-time": time.Now(),
	})
	for time.Now().Before(rate.Reset.Time) {
		time.Sleep(1 * time.Minute)
	}
	log.Info("needWait", map[string]interface{}{
		"end-time": time.Now(),
	})
	return true
}
