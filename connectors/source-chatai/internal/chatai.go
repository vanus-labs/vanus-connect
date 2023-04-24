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
	"context"
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
	SendChatCompletion(userIdentifier, content string) (string, error)
	Reset()
}

type ChatType string

const (
	chatGPT      ChatType = "chatgpt"
	chatErnieBot ChatType = "wenxin"
)

type chatService struct {
	chatGpt      ChatClient
	ernieBot     ChatClient
	config       *chatConfig
	lock         sync.RWMutex
	day          int
	limitContent string
	userNum      map[string]int
	ctx          context.Context
	cancel       context.CancelFunc
}

func newChatService(config *chatConfig) *chatService {
	s := &chatService{
		config:       config,
		chatGpt:      gpt.NewChatGPTService(config.GPT, config.MaxTokens, config.EnableContext),
		ernieBot:     ernie_bot.NewErnieBotService(config.ErnieBot, config.MaxTokens),
		day:          today(),
		limitContent: fmt.Sprintf("You've reached the daily limit (%d/day). Your quota will be restored tomorrow.", config.EverydayLimit),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	go func() {
		now := time.Now().UTC()
		next := now.Add(time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		t := time.NewTicker(next.Sub(now))
		select {
		case <-s.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			s.reset()
		}
		t.Stop()
		tk := time.NewTicker(time.Hour)
		defer tk.Stop()
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-tk.C:
				s.reset()
			}
		}
	}()
	return s
}

func today() int {
	return time.Now().UTC().Day()
}

func (s *chatService) Close() {
	s.cancel()
}

func (s *chatService) addNum(userIdentifier string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	num, ok := s.userNum[userIdentifier]
	if !ok {
		num = 0
	}
	num++
	s.userNum[userIdentifier] = num
}

func (s *chatService) getNum(userIdentifier string) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	num, ok := s.userNum[userIdentifier]
	if !ok {
		return 0
	}
	return num
}

func (s *chatService) reset() {
	s.lock.Lock()
	defer s.lock.Unlock()
	t := today()
	if s.day == t {
		return
	}
	s.day = t
	s.userNum = map[string]int{}
	s.chatGpt.Reset()
	s.ernieBot.Reset()
}

func (s *chatService) ChatCompletion(chatType ChatType, userIdentifier, content string) (resp string, err error) {
	num := s.getNum(userIdentifier)
	if num >= s.config.EverydayLimit {
		return s.limitContent, ErrLimit
	}
	log.Info("receive content:"+content, map[string]interface{}{
		"chat": chatType,
		"user": userIdentifier,
	})
	switch chatType {
	case chatErnieBot:
		resp, err = s.ernieBot.SendChatCompletion(userIdentifier, content)
	case chatGPT:
		resp, err = s.chatGpt.SendChatCompletion(userIdentifier, content)
	}
	if err != nil {
		return responseErr, err
	}
	if resp == "" {
		return responseEmpty, nil
	}
	s.addNum(userIdentifier)
	return resp, nil
}
