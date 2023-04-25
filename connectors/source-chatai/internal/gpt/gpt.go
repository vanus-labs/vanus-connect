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
	"sync"

	"github.com/sashabaranov/go-openai"
)

type chatGPTService struct {
	client        *openai.Client
	maxTokens     int
	enableContext bool
	userMap       map[string]*userMessage
	lock          sync.Mutex
}

type userMessage struct {
	messages   []openai.ChatCompletionMessage
	tokens     []int
	totalToken int
}

func (m *userMessage) cal(newToken, maxTokens int) {
	currToken := m.totalToken + newToken
	if currToken < maxTokens {
		return
	}
	var index, token int
	for index = range m.tokens {
		// question token
		token += m.tokens[index]
		index++
		// answer token
		token += m.tokens[index]
		if currToken-token < maxTokens {
			index = index + 1
			break
		}
	}
	m.totalToken -= token
	m.messages = m.messages[index:]
	m.tokens = m.tokens[index:]
}

func NewChatGPTService(config Config, maxTokens int, enableContext bool) *chatGPTService {
	client := openai.NewClient(config.Token)
	return &chatGPTService{
		client:        client,
		maxTokens:     maxTokens,
		enableContext: enableContext,
		userMap:       map[string]*userMessage{},
	}
}

func (s *chatGPTService) getUser(userIdentifier string) *userMessage {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, ok := s.userMap[userIdentifier]
	if !ok {
		user = &userMessage{}
		s.userMap[userIdentifier] = user
	}
	return user
}

func (s *chatGPTService) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.userMap = map[string]*userMessage{}
}

func (s *chatGPTService) SendChatCompletion(userIdentifier, content string) (string, error) {
	user := s.getUser(userIdentifier)
	if s.enableContext {
		s.lock.Lock()
		user.cal(calTokens(content), s.maxTokens)
		s.lock.Unlock()
	}
	messages := append(user.messages, openai.ChatCompletionMessage{
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
		s.lock.Lock()
		defer s.lock.Unlock()
		user.messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: respContent,
		})
		user.tokens = append(user.tokens, resp.Usage.PromptTokens-user.totalToken, resp.Usage.CompletionTokens)
		user.totalToken = resp.Usage.TotalTokens
	}
	return respContent, nil
}

func calTokens(content string) int {
	return len(content) / 4
}
