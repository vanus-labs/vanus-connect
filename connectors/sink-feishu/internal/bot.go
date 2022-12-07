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

package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
)

func (c *functionSink) sendTextToFeishuBot(e *v2.Event) error {
	t := time.Now()

	payload := map[string]interface{}{
		"sign":      c.genSignature(t),
		"timestamp": t.Unix(),
		"msg_type":  "text",
		"content": map[string]interface{}{
			"text": string(e.Data()),
		},
	}
	res, err := c.httpClient.R().SetBody(payload).Post(c.cfg.Bot.Webhook)
	if err != nil {
		return err
	}
	// {"code":19021,"data":{},"msg":"sign match fail or timestamp is not within one hour from current time"}
	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed call feishu: %s with HTTP Code %d",
			string(res.Body()), res.StatusCode())
	}
	return nil
}

func (c *functionSink) genSignature(t time.Time) string {
	strToSign := fmt.Sprintf("%d\n%s", t.Unix(), c.cfg.Secret.BotSignature)
	h := hmac.New(sha256.New, []byte(strToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
