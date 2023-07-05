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

package chat

import (
	"fmt"
	"github.com/vanus-labs/connector/source/chatai/chat/vanus_ai"

	"github.com/vanus-labs/connector/source/chatai/chat/ernie_bot"
	"github.com/vanus-labs/connector/source/chatai/chat/gpt"
)

type ChatConfig struct {
	GPT             gpt.Config       `json:"gpt" yaml:"gpt"`
	ErnieBot        ernie_bot.Config `json:"ernie_bot" yaml:"ernie_bot"`
	VanusAI         vanus_ai.Config  `json:"vanus_ai" yaml:"vanusai"`
	EverydayLimit   int              `json:"everyday_limit" yaml:"everyday_limit"`
	MaxTokens       int              `json:"max_tokens" yaml:"max_tokens"`
	EnableContext   bool             `json:"enable_context" yaml:"enable_context"`
	DefaultChatMode Type             `json:"default_chat_mode" yaml:"default_chat_mode"`
	Stream          bool             `json:"stream" yaml:"stream"`
}

func (c *ChatConfig) init() {
	if c.DefaultChatMode == "" {
		c.DefaultChatMode = ChatGPT
	}
	if c.EverydayLimit <= 0 {
		c.EverydayLimit = 1000
	}
	if c.MaxTokens <= 0 {
		c.MaxTokens = 3500
	}
}

func (c *ChatConfig) Validate() error {
	if c.DefaultChatMode != "" {
		switch c.DefaultChatMode {
		case ChatGPT, ChatErnieBot:
		case ChatVanusAI:
			return c.VanusAI.Validate()
		default:
			return fmt.Errorf("chat mode is invalid")
		}
	}
	return nil
}
