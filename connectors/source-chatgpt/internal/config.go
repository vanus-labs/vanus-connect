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

var _ cdkgo.SourceConfigAccessor = &chatGPTConfig{}

func NewChatGPTConfig() cdkgo.SourceConfigAccessor {
	return &chatGPTConfig{}
}

type chatGPTConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	Port          int    `json:"port" yaml:"port"`
	Token         string `json:"token" yaml:"token" validate:"required"`
	EverydayLimit int    `json:"everyday_limit" yaml:"everyday_limit"`
}

func (c *chatGPTConfig) Validate() error {
	return c.SourceConfig.Validate()
}

func (c *chatGPTConfig) Init() {
	if c.Port <= 0 {
		c.Port = 8080
	}
	if c.EverydayLimit <= 0 {
		c.EverydayLimit = 100
	}
}
