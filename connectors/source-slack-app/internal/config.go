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
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SourceConfigAccessor = &slackConfig{}

type slackConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	BotToken string `json:"bot_token" yaml:"bot_token" validate:"required"`
	AppToken string `json:"app_token" yaml:"app_token" validate:"required"`
	UserID   string `json:"user_id" yaml:"user_id"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &slackConfig{}
}
