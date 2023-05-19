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

package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

type Client struct {
	client      *resty.Client
	httpClient  *http.Client
	tokenSource oauth2.TokenSource
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		httpClient:  http.DefaultClient,
		client:      resty.New(),
		tokenSource: NewTokenSource(clientID, clientSecret),
	}
}

func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (*ChatCompletionResponse, error) {
	request.Stream = false
	req, err := c.newHttpRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp ChatCompletionResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorCode != 0 {
		return nil, fmt.Errorf("response error code:%d, msg:%s", resp.ErrorCode, resp.ErrorMsg)
	}
	return &resp, nil
}

func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (*ChatCompletionStream, error) {
	request.Stream = true
	req, err := c.newHttpRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, c.handleErrorResp(resp)
	}
	return &ChatCompletionStream{
		reader:         bufio.NewReader(resp.Body),
		response:       resp,
		errAccumulator: newErrorAccumulator(),
	}, nil
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes ErrResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil {
		return err
	}
	if errRes.ErrorCode != 0 {
		return fmt.Errorf("response error code:%d, msg:%s", errRes.ErrorCode, errRes.ErrorMsg)
	}
	return nil
}
