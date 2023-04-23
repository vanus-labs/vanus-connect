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
	"fmt"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/connector/source/chatai/internal/ernie_bot"
	"github.com/vanus-labs/connector/source/chatai/internal/gpt"
)

var _ cdkgo.SourceConfigAccessor = &chatConfig{}

func NewChatConfig() cdkgo.SourceConfigAccessor {
	return &chatConfig{}
}

type chatConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`

	Port            int              `json:"port" yaml:"port"`
	GPT             gpt.Config       `json:"gpt" yaml:"gpt"`
	ErnieBot        ernie_bot.Config `json:"ernie_bot" yaml:"ernie_bot"`
	EverydayLimit   int              `json:"everyday_limit" yaml:"everyday_limit"`
	MaxTokens       int              `json:"max_tokens" yaml:"max_tokens"`
	EnableContext   bool             `json:"enable_context" yaml:"enable_context"`
	DefaultChatMode ChatType         `json:"default_chat_mode" yaml:"default_chat_mode"`
}

func (c *chatConfig) Validate() error {
	if c.DefaultChatMode != "" {
		switch c.DefaultChatMode {
		case chatGPT, chatErnieBot:
		default:
			return fmt.Errorf("chat mode is invalid")
		}
	}
	return c.SourceConfig.Validate()
}

func (c *chatConfig) Init() {
	if c.Port <= 0 {
		c.Port = 8080
	}
	if c.EverydayLimit <= 0 {
		c.EverydayLimit = 100
	}
	if c.MaxTokens <= 0 {
		c.MaxTokens = 3500
	}
	if c.DefaultChatMode == "" {
		c.DefaultChatMode = chatGPT
	}
}
