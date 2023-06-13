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

	"github.com/pandodao/tokenizer-go"
	"github.com/sashabaranov/go-openai"

	"github.com/vanus-labs/connector/source/chatai/chat/model"
)

type chatGPTService struct {
	client        *openai.Client
	maxTokens     int
	enableContext bool
	userMap       map[string]*userMessage
	lock          sync.Mutex
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

func (s *chatGPTService) SendChatCompletion(ctx context.Context, userIdentifier, content string) (string, error) {
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
			User:     userIdentifier,
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

func (s *chatGPTService) SendChatCompletionStream(ctx context.Context, userIdentifier, content string) (model.ChatCompletionStream, error) {
	content, err := s.SendChatCompletion(ctx, userIdentifier, content)
	if err != nil {
		return nil, err
	}
	return newChatCompletionStream(content), nil
}

type chatCompletionStream struct {
	isFinish bool
	content  string
}

func newChatCompletionStream(content string) model.ChatCompletionStream {
	return &chatCompletionStream{
		content: content,
	}
}

func (s *chatCompletionStream) Recv() (*model.StreamMessage, error) {
	if s.isFinish {
		return nil, nil
	}
	s.isFinish = true
	return &model.StreamMessage{
		Index:   0,
		IsEnd:   true,
		Content: s.content,
	}, nil
}

func (s *chatCompletionStream) Close() {

}

func calTokens(content string) int {
	t, err := tokenizer.CalToken(content)
	if err != nil {
		t = len(content) / 4
	}
	return t
}
