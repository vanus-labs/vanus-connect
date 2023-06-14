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
	"fmt"
	cdkgo "github.com/vanus-labs/cdk-go"
)

type Type string

const (
	Contributor Type = "contributor"
	PR          Type = "pr"

	ListByOrg  Type = "org"
	ListByUser Type = "user"
)

var _ cdkgo.SourceConfigAccessor = &GitHubAPIConfig{}

type GitHubAPIConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	APIType       Type       `json:"api_type" yaml:"api_type" validate:"required"`
	ListType      Type       `json:"list_type" yaml:"list_type"`
	Organizations []string   `json:"organizations" yaml:"organizations"`
	UserList      []string   `json:"user_list" yaml:"user_list"`
	PRConfigs     []PRConfig `json:"pr_configs" yaml:"pr_configs"`

	GitHubAccessToken string `json:"github_access_token" yaml:"github_access_token" validate:"required"`
	GitHubHourLimit   int    `json:"github_hour_limit" yaml:"github_hour_limit"`
}

type PRConfig struct {
	Organization string   `json:"organization" yaml:"organization"`
	Repo         string   `json:"repo" yaml:"repo"`
	UserList     []string `json:"user_list" yaml:"user_list"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &GitHubAPIConfig{}
}

func (c *GitHubAPIConfig) Validate() error {
	if c.APIType != "" {
		switch c.APIType {
		case PR:
			if len(c.PRConfigs) == 0 {
				return fmt.Errorf("API type is '%s', should have pr_config", PR)
			}
		case Contributor:
			if c.ListType == ListByOrg && len(c.Organizations) == 0 {
				return fmt.Errorf("API type is '%s', should have organizations", ListByOrg)
			}
			if c.ListType == ListByUser && len(c.UserList) == 0 {
				return fmt.Errorf("API type is '%s', should have users", ListByUser)
			}
		default:
			return fmt.Errorf("API type is invalid")
		}
	}
	return nil
}

func (c *GitHubAPIConfig) Init() {
	if c.GitHubHourLimit == 0 {
		c.GitHubHourLimit = 3600
	}
}
