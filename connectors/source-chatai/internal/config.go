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
	"github.com/vanus-labs/connector/source/chatai/chat"
)

var _ cdkgo.SourceConfigAccessor = &chatConfig{}

func NewChatConfig() cdkgo.SourceConfigAccessor {
	return &chatConfig{}
}

type chatConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	chat.ChatConfig      `json:",inline" yaml:",inline"`
	Port                 int    `json:"port" yaml:"port"`
	DefaultProcessMode   string `json:"default_process_mode" yaml:"default_process_mode"`
	UserIdentifierHeader string `json:"user_identifier_header" yaml:"user_identifier_header"`
	Auth                 *Auth  `json:"auth" yaml:"auth"`
}

type Auth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (a *Auth) IsEmpty() bool {
	return a == nil || a.Username == "" || a.Password == ""
}

func (c *chatConfig) Validate() error {
	err := c.ChatConfig.Validate()
	if err != nil {
		return err
	}
	return c.SourceConfig.Validate()
}
