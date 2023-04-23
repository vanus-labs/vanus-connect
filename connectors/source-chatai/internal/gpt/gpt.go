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

package gpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type chatGPTService struct {
	client        *openai.Client
	maxTokens     int
	enableContext bool
	messages      []openai.ChatCompletionMessage
	tokens        []int
	totalToken    int
}

func NewChatGPTService(config Config, maxTokens int, enableContext bool) *chatGPTService {
	client := openai.NewClient(config.Token)
	return &chatGPTService{
		client:        client,
		maxTokens:     maxTokens,
		enableContext: enableContext,
	}
}

func (s *chatGPTService) Reset() {
	s.messages = nil
	s.tokens = nil
	s.totalToken = 0
}

func (s *chatGPTService) SendChatCompletion(content string) (string, error) {
	currToken := s.totalToken + calTokens(content)
	if s.enableContext && currToken > s.maxTokens {
		var index, token int
		for index = range s.tokens {
			// question token
			token += s.tokens[index]
			index++
			// answer token
			token += s.tokens[index]
			if currToken-token < s.maxTokens {
				index = index + 1
				break
			}
		}
		s.totalToken -= token
		s.messages = s.messages[index:]
		s.tokens = s.tokens[index:]
	}
	messages := append(s.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	})
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	respContent := resp.Choices[0].Message.Content
	if s.enableContext {
		s.messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: respContent,
		})
		s.tokens = append(s.tokens, resp.Usage.PromptTokens-s.totalToken, resp.Usage.CompletionTokens)
		s.totalToken = resp.Usage.TotalTokens
	}
	return respContent, nil
}

func calTokens(content string) int {
	return len(content) / 4
}
