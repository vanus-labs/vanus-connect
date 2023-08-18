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
	"encoding/json"
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/go-github/v52/github"
	"github.com/rs/zerolog"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/sink/github/internal/vanusai"
)

const (
	name = "Github Sink"
)

func NewSlackSink() cdkgo.Sink {
	return &githubSink{}
}

var _ cdkgo.Sink = &githubSink{}

type githubSink struct {
	count  int64
	client *github.Client
	aiCli  *vanusai.Client
	logger zerolog.Logger
}

func (e *githubSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for idx := range events {
		event := events[idx]
		var issue IssueData
		err := json.Unmarshal(event.Data(), &issue)
		if err != nil {
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		if issue.Body == "" {
			e.logger.Info().Str("event_id", event.ID()).Msg("body is empty")
			continue
		}
		i, err := issue.GetIssue()
		if err != nil {
			e.logger.Warn().Err(err).
				Str("event_id", event.ID()).
				Str("url", issue.URL).Msg("url is invalid")
			return cdkgo.NewResult(http.StatusBadRequest, err.Error())
		}
		prompt, err := e.aiCli.Chat(ctx, vanusai.NewChatRequest(issue.Body, issue.URL))
		if err != nil {
			e.logger.Warn().Err(err).
				Str("prompt", issue.Body).Msg("call ai failed")
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		}
		_, _, err = e.client.Issues.CreateComment(ctx, i.Owner, i.Repo, i.Number, &github.IssueComment{
			Body: &prompt,
		})
		if err != nil {
			e.logger.Warn().Err(err).
				Str("event_id", event.ID()).
				Str("url", issue.URL).Msg("create comment failed")
		} else {
			e.logger.Info().
				Str("event_id", event.ID()).
				Str("url", issue.URL).Msg("create comment success")
		}
	}
	return cdkgo.SuccessResult
}

func (e *githubSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	e.logger = log.FromContext(ctx)
	config := cfg.(*githubConfig)
	e.client = github.NewTokenClient(context.TODO(), config.Github.AccessToken)
	e.aiCli = vanusai.NewClient(config.VanusAIURL)
	return nil
}

func (e *githubSink) Name() string {
	return name
}

func (e *githubSink) Destroy() error {
	// nothing to do
	return nil
}
