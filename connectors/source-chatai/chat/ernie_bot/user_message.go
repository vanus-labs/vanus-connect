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
	"sync"

	"github.com/vanus-labs/connector/source/chatai/chat/ernie_bot/client"
)

type userMessage struct {
	messages   []client.ChatCompletionMessage
	tokens     []int
	totalToken int
	lock       sync.Mutex
}

type tokens struct {
	prompt     int
	completion int
	total      int
}

func (m *userMessage) reset() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.tokens = nil
	m.messages = nil
	m.totalToken = 0
}

func (m *userMessage) set(message []client.ChatCompletionMessage, tokens tokens) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.messages = append(m.messages, message...)
	m.tokens = append(m.tokens, tokens.prompt-m.totalToken, tokens.completion)
	m.totalToken = tokens.total
}

func (m *userMessage) cal(newToken, maxTokens int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	currToken := m.totalToken + newToken
	if currToken < maxTokens {
		return
	}
	var index, token int
	for index < len(m.tokens) {
		// question token
		token += m.tokens[index]
		// answer token
		token += m.tokens[index+1]
		index += 2
		if currToken-token < maxTokens {
			break
		}
	}
	m.totalToken -= token
	m.messages = m.messages[index:]
	m.tokens = m.tokens[index:]
}
