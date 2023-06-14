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
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SourceConfigAccessor = &DouyinConfig{}

type DouyinConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	RateHourLimit int `json:"rate_hour_limit" yaml:"rate_hour_limit"`

	OpenID       string `json:"open_id" yaml:"open_id" validate:"required"`
	ClientKey    string `json:"client_key" yaml:"client_key" validate:"required"`
	ClientSecret string `json:"client_secret" yaml:"client_secret" validate:"required"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &DouyinConfig{}
}

func (c *DouyinConfig) Init() {
	if c.RateHourLimit == 0 {
		c.RateHourLimit = 3600
	}
}
