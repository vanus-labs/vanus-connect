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
	"sync"
	"time"

	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/source/chatai/internal/ernie_bot"
	"github.com/vanus-labs/connector/source/chatai/internal/gpt"
)

const (
	responseEmpty = "Get response empty."
	responseErr   = "Get response failed."
)

var (
	ErrLimit = fmt.Errorf("reached the daily limit")
)

type ChatClient interface {
	SendChatCompletion(content string) (string, error)
}

type ChatType string

const (
	chatGPT    ChatType = "chatgpt"
	chatWenxin ChatType = "wenxin"
)

type chatService struct {
	chatGpt      ChatClient
	wenxin       ChatClient
	config       *chatConfig
	lock         sync.Mutex
	day          int
	num          int
	limitContent string
}

func newChatService(config *chatConfig) *chatService {
	return &chatService{
		config:       config,
		chatGpt:      gpt.NewChatGPTService(config.GPT, config.MaxTokens),
		wenxin:       ernie_bot.NewErnieBotService(config.ErnieBot, config.MaxTokens),
		day:          today(),
		limitContent: fmt.Sprintf("You've reached the daily limit (%d/day). Your quota will be restored tomorrow.", config.EverydayLimit),
	}
}

func today() int {
	return time.Now().UTC().Day()
}

func (s *chatService) reset() {
	s.day = today()
	s.num = 0
}

func (s *chatService) ChatCompletion(chatType ChatType, content string) (resp string, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.num >= s.config.EverydayLimit {
		if today() == s.day {
			return s.limitContent, ErrLimit
		}
		s.reset()
	}
	log.Info("receive content:"+content, map[string]interface{}{
		"chat": chatType,
	})
	switch chatType {
	case chatWenxin:
		resp, err = s.wenxin.SendChatCompletion(content)
	case chatGPT:
		resp, err = s.chatGpt.SendChatCompletion(content)
	}
	if err != nil {
		return responseErr, err
	}
	if resp == "" {
		return responseEmpty, nil
	}
	s.num++
	return resp, nil
}
