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
	cdkgo "github.com/linkall-labs/cdk-go"
)

type GitHubConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	Port   int       `json:"port" yaml:"port"`
	GitHub GitHubCfg `json:"github" yaml:"github"`
}

func Config() cdkgo.SourceConfigAccessor {
	return &GitHubConfig{}
}

func (cfg *GitHubConfig) GetSecret() cdkgo.SecretAccessor {
	return &cfg.GitHub
}

type GitHubCfg struct {
	AccessToken   string `json:"access_token" yaml:"access_token"`
	WebHookSecret string `json:"webhook_secret" yaml:"webhook_secret"`
}
