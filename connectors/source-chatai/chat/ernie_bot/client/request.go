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
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const URL = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions"

func (c *Client) newHttpRequest(ctx context.Context, request interface{}) (*http.Request, error) {
	t, err := c.tokenSource.Token()
	if err != nil {
		return nil, err
	}
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		URL,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return nil, err
	}
	parameters := url.Values{}
	parameters.Set("access_token", t.AccessToken)
	req.URL.RawQuery = parameters.Encode()
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
