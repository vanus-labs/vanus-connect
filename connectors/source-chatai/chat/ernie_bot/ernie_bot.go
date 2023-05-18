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

package ernie_bot

import (
	"context"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"

	"github.com/vanus-labs/connector/source/chatai/chat/ernie_bot/client"
	"github.com/vanus-labs/connector/source/chatai/chat/model"
)

const url = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions"

type ernieBotService struct {
	client        *client.Client
	maxTokens     int
	enableContext bool
	userMap       map[string]*userMessage
	lock          sync.Mutex
}

func NewErnieBotService(config Config, maxTokens int, enableContext bool) *ernieBotService {
	return &ernieBotService{
		maxTokens:     1500,
		enableContext: enableContext,
		userMap:       map[string]*userMessage{},
		client:        client.NewClient(config.AccessKey, config.SecretKey),
	}
}

func (s *ernieBotService) getUser(userIdentifier string) *userMessage {
	s.lock.Lock()
	defer s.lock.Unlock()
	user, ok := s.userMap[userIdentifier]
	if !ok {
		user = &userMessage{}
		s.userMap[userIdentifier] = user
	}
	return user
}

func (s *ernieBotService) Reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.userMap = map[string]*userMessage{}
}

func (s *ernieBotService) SendChatCompletion(ctx context.Context, userIdentifier, content string) (string, error) {
	user := s.getUser(userIdentifier)
	if s.enableContext {
		user.cal(calTokens(content), s.maxTokens)
	}
	question := client.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
	messages := append(user.messages, question)
	req := client.ChatCompletionRequest{
		Message: messages,
		User:    userIdentifier,
	}
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	respContent := resp.Result
	if s.enableContext {
		if resp.NeedClearHistory {
			user.reset()
			return respContent, nil
		}
		user.set([]client.ChatCompletionMessage{question, {
			Role:    openai.ChatMessageRoleAssistant,
			Content: respContent,
		}}, tokens{
			prompt:     resp.Usage.PromptTokens,
			completion: resp.Usage.CompletionTokens,
			total:      resp.Usage.TotalTokens})
	}
	return respContent, nil
}

func (s *ernieBotService) SendChatCompletionStream(ctx context.Context, userIdentifier, content string) (model.ChatCompletionStream, error) {
	user := s.getUser(userIdentifier)
	if s.enableContext {
		user.cal(calTokens(content), s.maxTokens)
	}
	question := client.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}
	messages := append(user.messages, question)
	req := client.ChatCompletionRequest{
		Message: messages,
		User:    userIdentifier,
		Stream:  true,
	}
	stream, err := s.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return newChatCompletionStream(stream, question, user, s.enableContext), nil
}

type chatCompletionStream struct {
	chat             *client.ChatCompletionStream
	buffer           strings.Builder
	question         client.ChatCompletionMessage
	user             *userMessage
	promptTokens     int
	completionTokens int
	totalTokens      int
	isFinish         bool
	needClearHistory bool
	enableContext    bool
}

func newChatCompletionStream(chat *client.ChatCompletionStream, question client.ChatCompletionMessage, user *userMessage, enableContext bool) model.ChatCompletionStream {
	return &chatCompletionStream{
		chat:          chat,
		user:          user,
		question:      question,
		enableContext: enableContext,
	}
}

func (s *chatCompletionStream) Recv() (*model.StreamMessage, error) {
	if s.isFinish {
		return nil, nil
	}
	resp, err := s.chat.Recv()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		return nil, err
	}
	s.isFinish = resp.IsEnd
	if s.enableContext {
		if !s.needClearHistory {
			s.needClearHistory = resp.NeedClearHistory
			s.promptTokens = resp.Usage.PromptTokens
			s.totalTokens = resp.Usage.TotalTokens
			s.completionTokens += resp.Usage.CompletionTokens
			s.buffer.WriteString(resp.Result)
		}
		if s.isFinish {
			s.doFinish()
		}
	}
	return &model.StreamMessage{
		ID:      resp.ID,
		Index:   resp.SentenceID,
		IsEnd:   resp.IsEnd,
		Content: resp.Result,
	}, nil
}

func (s *chatCompletionStream) Close() {
	s.chat.Close()
}

func (s *chatCompletionStream) doFinish() {
	if s.needClearHistory {
		s.user.reset()
		return
	}
	s.user.set([]client.ChatCompletionMessage{s.question, {
		Role:    openai.ChatMessageRoleAssistant,
		Content: s.buffer.String(),
	}}, tokens{
		prompt:     s.promptTokens,
		completion: s.completionTokens,
		total:      s.totalTokens,
	})
}

func calTokens(content string) int {
	return len(content)
}
