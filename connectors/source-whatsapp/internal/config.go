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

var _ cdkgo.SourceConfigAccessor = &whatsAppConfig{}

func NewExampleConfig() cdkgo.SourceConfigAccessor {
	return &whatsAppConfig{}
}

type whatsAppConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	// TODO
	Secret Secret `json:"secret" yaml:"secret"`

	*chat.ChatConfig `json:",inline" yaml:",inline"`
	EnableChatAi     bool `json:"enable_chatai" yaml:"enable_chatai"`
}

func (c *whatsAppConfig) GetSecret() cdkgo.SecretAccessor {
	return &c.Secret
}

func (c *whatsAppConfig) Validate() error {
	// TODO
	return c.SourceConfig.Validate()
}

type Secret struct {
}
