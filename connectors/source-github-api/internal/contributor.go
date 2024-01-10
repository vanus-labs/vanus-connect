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

	"github.com/google/go-github/v58/github"
	lodash "github.com/samber/lo"
)

func (s *GitHubAPISource) startContr(ctx context.Context) {
	switch s.config.ListType {
	case ListByOrg:
		for i := range s.config.Organizations {
			orgName := s.config.Organizations[i]
			s.listOrgRepo(ctx, orgName)
		}
	case ListByUser:
		for i := range s.config.UserList {
			user := s.config.UserList[i]
			s.listUserRepo(ctx, user)
		}
	}

}

func (s *GitHubAPISource) listOrgRepo(ctx context.Context, orgName string) {
	// Repository
	listOption := &github.RepositoryListByOrgOptions{
		Type: "sources",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 250,
		},
	}

	for {
		s.Limiter.Take()
		repos, resp, err := s.client.Repositories.ListByOrg(ctx, orgName, listOption)
		if err != nil {
			s.logger.Warn().Err(err).Msg("repo list by org error")
		}
		if len(repos) == 0 {
			break
		}
		s.logger.Info().
			Str("org", orgName).
			Int("page", listOption.ListOptions.Page).
			Int("next_page", resp.NextPage).
			Interface("rate", resp.Rate).
			Msg("list by org")

		for _, repo := range repos {
			if *repo.StargazersCount < 1000 {
				continue
			}
			s.numRepos += 1
			s.listContributors(ctx, repo)
		}

		if resp.NextPage <= listOption.ListOptions.Page {
			break
		}
		listOption.ListOptions.Page = resp.NextPage
	}
}

func (s *GitHubAPISource) listUserRepo(ctx context.Context, user string) {
	// Repository
	listOption := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 250,
		},
	}

	for {
		s.Limiter.Take()
		repos, resp, err := s.client.Repositories.List(ctx, user, listOption)
		if err != nil {
			s.logger.Warn().Err(err).Msg("resp list error")
		}
		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if *repo.StargazersCount < 1000 {
				continue
			}
			s.numRepos += 1
			s.listContributors(ctx, repo)
		}

		if resp.NextPage <= listOption.ListOptions.Page {
			break
		}
		listOption.ListOptions.Page = resp.NextPage
	}
}

func (s *GitHubAPISource) listContributors(ctx context.Context, repo *github.Repository) {
	// Contributors
	listOption := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 250,
		},
	}
	projectRecords := 0
	for {
		s.Limiter.Take()
		contributors, resp, err := s.client.Repositories.ListContributors(ctx, *repo.Owner.Login, *repo.Name, listOption)
		if err != nil {
			s.logger.Warn().Err(err).Msg("repo list contributors errors")
		}
		if len(contributors) == 0 {
			break
		}

		s.numRecords += len(contributors)
		projectRecords += len(contributors)
		for _, contributor := range contributors {
			s.userInfo(ctx, contributor, repo)
		}
		s.logger.Info().
			Str("project", *repo.Name).
			Int("page", listOption.ListOptions.Page).
			Int("next_page", resp.NextPage).
			Int("total_repo", s.numRepos).
			Int("total_record", s.numRecords).
			Int("project_record", projectRecords).
			Msg("list contributors")

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
			s.logger.Warn().Err(err).Msg("user get error")
			return
		}
		user = user0
		s.m.Store(*contributor.Login, user)
	} else {
		user = v.(*github.User)
	}

	data := make(map[string]interface{})
	data["repo"] = lodash.TernaryF(repo.Name != nil, func() string { return *repo.Name }, func() string { return "" })
	data["star"] = lodash.TernaryF(repo.StargazersCount != nil, func() int { return *repo.StargazersCount }, func() int { return 0 })
	data["org"] = lodash.TernaryF(repo.Owner.Login != nil, func() string { return *repo.Owner.Login }, func() string { return "" })
	data["url"] = lodash.TernaryF(repo.HTMLURL != nil, func() string { return *repo.HTMLURL }, func() string { return "" })
	data["uid"] = lodash.TernaryF(user.Login != nil, func() string { return *user.Login }, func() string { return "" })
	data["username"] = lodash.TernaryF(user.Name != nil, func() string { return *user.Name }, func() string { return "" })
	data["email"] = lodash.TernaryF(user.Email != nil, func() string { return *user.Email }, func() string { return "" })
	data["company"] = lodash.TernaryF(user.Company != nil, func() string { return *user.Company }, func() string { return "" })

	s.sendEvent("contributors", data["org"].(string), data)
}
