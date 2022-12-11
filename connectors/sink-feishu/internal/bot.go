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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/go-resty/resty/v2"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/tidwall/gjson"
)

type messageType string

const (
	xChatGroupID = "xvfeishuchatgroup"
	xMessageType = "xvfeishumsgtype"

	textMessage        = messageType("text")
	postMessage        = messageType("post")
	shareChatMessage   = messageType("share_chat")
	imageMessage       = messageType("image")
	interactiveMessage = messageType("interactive")
)

var (
	errChatGroup = errors.New("feishu: xvfeishuchatgroup is missing or incorrect, please check your Feishu " +
		"Sink config or subscription")
	errMessageType = errors.New("feishu: xvfeishumsgtype is missing or invalid, only" +
		" [text, post, share_chat, image, interactive] are supported")
	errInvalidPostMessage = errors.New("feishu: invalid post message, please make sure it's the json format")
)

type bot struct {
	cfg        BotConfig
	cm         map[string]WebHook
	httpClient *resty.Client
}

func (b *bot) init(cfg BotConfig) error {
	b.cfg = cfg
	b.cm = make(map[string]WebHook, len(cfg.Webhooks))
	for _, wh := range cfg.Webhooks {
		_, exist := b.cm[wh.ChatGroup]
		if exist {
			return fmt.Errorf("the chat_group has conflicted with name: %s", wh.ChatGroup)
		}
		b.cm[wh.ChatGroup] = wh
	}
	return nil
}

type WebHook struct {
	Address   string `json:"address" yaml:"address"`
	Signature string `json:"signature" yaml:"signature"`
	ChatGroup string `json:"chat_group" yaml:"chat_group"`
}

func (wh WebHook) Validate() error {
	if wh.Address == "" {
		return errors.New("webhook address is nil")
	}
	if wh.Signature == "" {
		return errors.New("webhook signature is nil")
	}
	if wh.ChatGroup == "" {
		return errors.New("webhook chat_group is nil")
	}
	return nil
}

type BotConfig struct {
	Webhooks []WebHook `json:"webhooks" yaml:"webhooks"`
}

func (bc BotConfig) Validate() error {
	if len(bc.Webhooks) == 0 {
		return errors.New("feishu: bot webhooks can't be empty")
	}
	for _, v := range bc.Webhooks {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) sendMessage(e *v2.Event) error {

	v := e.Extensions()[xChatGroupID]
	groupID, err := types.ToString(v)
	if err != nil {
		return errChatGroup
	}

	wh, exist := b.cm[groupID]
	if !exist {
		return errChatGroup
	}

	v = e.Extensions()[xMessageType]
	t, ok := v.(string)
	if !ok {
		return errMessageType
	}
	switch messageType(t) {
	case textMessage:
		return b.sendTextMessage(e, wh)
	case postMessage:
		return b.sendPostMessage(e, wh)
	case shareChatMessage:
		return b.sendShareChatMessage(e, wh)
	case imageMessage:
		return b.sendImageMessage(e, wh)
	case interactiveMessage:
		return b.sendInteractiveMessage(e, wh)
	default:
		return errMessageType
	}
}

func (b *bot) sendTextMessage(e *v2.Event, wh WebHook) error {
	content := map[string]interface{}{
		"text": string(e.Data()),
	}
	res, err := b.httpClient.R().SetBody(b.generatePayload(content, textMessage, wh)).Post(wh.Address)
	if err != nil {
		return err
	}

	return b.processResponse(e, res)
}

func (b *bot) sendPostMessage(e *v2.Event, wh WebHook) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(e.Data(), &m); err != nil {
		return errInvalidPostMessage
	}
	content := map[string]interface{}{
		"post": m,
	}

	res, err := b.httpClient.R().SetBody(b.generatePayload(content, postMessage, wh)).Post(wh.Address)
	if err != nil {
		return err
	}
	return b.processResponse(e, res)
}

func (b *bot) sendShareChatMessage(e *v2.Event, wh WebHook) error {
	content := map[string]interface{}{
		"share_chat_id": string(e.Data()),
	}

	res, err := b.httpClient.R().SetBody(b.generatePayload(content, shareChatMessage, wh)).Post(wh.Address)
	if err != nil {
		return err
	}
	return b.processResponse(e, res)
}

func (b *bot) sendImageMessage(e *v2.Event, wh WebHook) error {
	content := map[string]interface{}{
		"image_key": string(e.Data()),
	}

	res, err := b.httpClient.R().SetBody(b.generatePayload(content, imageMessage, wh)).Post(wh.Address)
	if err != nil {
		return err
	}
	return b.processResponse(e, res)
}

func (b *bot) sendInteractiveMessage(e *v2.Event, wh WebHook) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(e.Data(), &m); err != nil {
		return errInvalidPostMessage
	}

	t := time.Now()
	payload := map[string]interface{}{
		"sign":      b.genSignature(t, wh.Signature),
		"timestamp": t.Unix(),
		"msg_type":  interactiveMessage,
		"card":      m,
	}
	res, err := b.httpClient.R().SetBody(payload).Post(wh.Address)
	if err != nil {
		return err
	}
	return b.processResponse(e, res)
}

func (b *bot) genSignature(t time.Time, signature string) string {
	strToSign := fmt.Sprintf("%d\n%s", t.Unix(), signature)
	h := hmac.New(sha256.New, []byte(strToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (b *bot) generatePayload(content map[string]interface{}, msgType messageType, wh WebHook) interface{} {
	t := time.Now()
	return map[string]interface{}{
		"sign":      b.genSignature(t, wh.Signature),
		"timestamp": t.Unix(),
		"msg_type":  msgType,
		"content":   content,
	}
}

func (b *bot) processResponse(e *v2.Event, res *resty.Response) error {
	// docs: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN?lang=zh-CN#756b882f
	obj := gjson.ParseBytes(res.Body())
	if obj.Get("StatusCode").Int() == 0 &&
		obj.Get("StatusMessage").String() == "success" {
		log.Debug("success send message to Feishu Bot", map[string]interface{}{
			"id": e.ID(),
		})
		return nil
	}
	return fmt.Errorf("failed to call feishu: %s", string(res.Body()))
}
