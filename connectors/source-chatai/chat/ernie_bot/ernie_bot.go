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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"

	"github.com/vanus-labs/connector/source/chatai/chat/ernie_bot/oauth"
)

const url = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions"

type ernieBotService struct {
	client      *resty.Client
	tokenSource oauth2.TokenSource
	config      Config
	maxTokens   int
}

func NewErnieBotService(config Config, maxTokens int) *ernieBotService {
	return &ernieBotService{
		config:      config,
		maxTokens:   maxTokens,
		client:      resty.New(),
		tokenSource: oauth.NewTokenSource(config.AccessKey, config.SecretKey),
	}
}
func (s *ernieBotService) Reset() {

}

func (s *ernieBotService) SendChatCompletion(userIdentifier, content string) (string, error) {
	token, err := s.tokenSource.Token()
	if err != nil {
		return "", err
	}
	req := ChatCompletionRequest{
		Message: []ChatCompletionMessage{{
			Role:    "user",
			Content: content,
		}},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	res, err := s.client.R().SetQueryParam("access_token", token.AccessToken).
		SetHeader("Content-Type", "application/json").SetBody(bytes.NewBuffer(reqBytes)).Post(url)
	if err != nil {
		return "", err
	}
	var resp ChatCompletionResponse
	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != 0 {
		return "", fmt.Errorf("response error code:%d, msg:%s", resp.ErrorCode, resp.ErrorMsg)
	}
	return resp.Result, nil
}
