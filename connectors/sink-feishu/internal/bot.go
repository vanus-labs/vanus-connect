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
	"strings"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/go-resty/resty/v2"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/tidwall/gjson"
)

type messageType string

const (
	xChatGroupID  = "xvfeishuchatgroup"
	xMessageType  = "xvfeishumsgtype"
	xBotURL       = "xvfeishuboturls"
	xBotSignature = "xvfeishubotsigns"

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
	errInvalidPostMessage     = errors.New("feishu: invalid post message, please make sure it's the json format")
	errInvalidAttributes      = errors.New("feishu: invalid xvfeishuboturls or xvfeishubotsigs")
	errInvalidAttributeNumber = errors.New("feishu: the number of bot url and signature must be equal")
	errNoBotWebhookFound      = errors.New("feishu: no feishu bot target webhook found")
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
	ChatGroup string `json:"chat_group" yaml:"chat_group" validate:"required"`
	URL       string `json:"url" yaml:"url" validate:"required"`
	Signature string `json:"signature" yaml:"signature"`
}

type BotConfig struct {
	Webhooks     []WebHook `json:"webhooks" yaml:"webhooks" validate:"dive"`
	DynamicRoute bool      `json:"dynamic_route" yaml:"dynamic_route"`
}

func (c *BotConfig) Validate() error {
	if !c.DynamicRoute && len(c.Webhooks) == 0 {
		return errors.New("the bot.webhooks can't be empty when dynamic_route is false")
	}
	return nil
}

func (b *bot) sendMessage(e *v2.Event) (err error) {
	var (
		whs     []WebHook
		groupID string
	)
	defer func() {
		if err != nil {
			d, _ := e.MarshalJSON()
			log.Warning("failed to send message", map[string]interface{}{
				log.KeyError: err,
				"event":      string(d),
				"webhooks":   whs,
			})
		}
	}()
	v := e.Extensions()[xChatGroupID]

	groupID, err = types.ToString(v)
	if err != nil && !b.cfg.DynamicRoute {
		return errChatGroup
	} else {
		wh, exist := b.cm[groupID]
		if !exist {
			if !b.cfg.DynamicRoute {

				return errChatGroup
			}
		} else {
			whs = append(whs, wh)
		}
	}

	if b.cfg.DynamicRoute {
		v = e.Extensions()[xBotURL]
		urlAttr, ok := v.(string)
		if !ok {
			return errInvalidAttributes
		}
		v = e.Extensions()[xBotSignature]
		signatureAttr, ok := v.(string)
		if !ok {
			return errInvalidAttributes
		}
		urls := strings.Split(urlAttr, ",")
		signatures := strings.Split(signatureAttr, ",")
		if len(urls) != len(signatures) {
			return errInvalidAttributeNumber
		}
		for idx := range urls {
			whs = append(whs, WebHook{
				URL:       urls[idx],
				Signature: signatures[idx],
			})
		}
	}

	if len(whs) == 0 {
		return errNoBotWebhookFound
	}

	v = e.Extensions()[xMessageType]
	t, ok := v.(string)
	if !ok {
		return errMessageType
	}
	switch messageType(t) {
	case textMessage:
		return b.sendTextMessage(e, whs)
	case postMessage:
		return b.sendPostMessage(e, whs)
	case shareChatMessage:
		return b.sendShareChatMessage(e, whs)
	case imageMessage:
		return b.sendImageMessage(e, whs)
	case interactiveMessage:
		return b.sendInteractiveMessage(e, whs)
	default:
		return errMessageType
	}
}

func (b *bot) sendTextMessage(e *v2.Event, whs []WebHook) error {
	content := map[string]interface{}{
		"text": string(e.Data()),
	}

	for _, wh := range whs {
		if wh.URL == "" {
			continue
		}
		res, err := b.httpClient.R().SetBody(b.generatePayload(content, textMessage, wh)).Post(wh.URL)
		if err != nil {

			return err
		}
		if err = b.processResponse(e, res); err != nil {

			return err
		}
	}
	return nil
}

func (b *bot) sendPostMessage(e *v2.Event, whs []WebHook) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(trim(e.Data()), &m); err != nil {

		return errInvalidPostMessage
	}
	content := map[string]interface{}{
		"post": m,
	}
	for _, wh := range whs {
		if wh.URL == "" {
			continue
		}
		res, err := b.httpClient.R().SetBody(b.generatePayload(content, postMessage, wh)).Post(wh.URL)
		if err != nil {

			return err
		}
		if err = b.processResponse(e, res); err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) sendShareChatMessage(e *v2.Event, whs []WebHook) error {
	content := map[string]interface{}{
		"share_chat_id": string(e.Data()),
	}

	for _, wh := range whs {
		if wh.URL == "" {
			continue
		}
		res, err := b.httpClient.R().SetBody(b.generatePayload(content, shareChatMessage, wh)).Post(wh.URL)
		if err != nil {
			return err
		}
		if err = b.processResponse(e, res); err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) sendImageMessage(e *v2.Event, whs []WebHook) error {
	content := map[string]interface{}{
		"image_key": string(e.Data()),
	}
	for _, wh := range whs {
		if wh.URL == "" {
			continue
		}
		res, err := b.httpClient.R().SetBody(b.generatePayload(content, imageMessage, wh)).Post(wh.URL)
		if err != nil {
			return err
		}
		if err = b.processResponse(e, res); err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) sendInteractiveMessage(e *v2.Event, whs []WebHook) error {
	m := map[string]interface{}{}

	if err := json.Unmarshal(trim(e.Data()), &m); err != nil {
		return errInvalidPostMessage
	}

	t := time.Now()
	payload := map[string]interface{}{
		"timestamp": t.Unix(),
		"msg_type":  interactiveMessage,
		"card":      m,
	}

	for _, wh := range whs {
		if wh.URL == "" {
			continue
		}
		if wh.Signature != "" {
			payload["sign"] = b.genSignature(t, wh.Signature)
		}
		res, err := b.httpClient.R().SetBody(payload).Post(wh.URL)
		if err != nil {
			return err
		}
		if err = b.processResponse(e, res); err != nil {
			return err
		}
		delete(payload, "sign")
	}

	return nil
}

func (b *bot) genSignature(t time.Time, signature string) string {
	if signature == "" {
		return ""
	}
	strToSign := fmt.Sprintf("%d\n%s", t.Unix(), signature)
	h := hmac.New(sha256.New, []byte(strToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (b *bot) generatePayload(content map[string]interface{}, msgType messageType, wh WebHook) interface{} {
	t := time.Now()
	payload := map[string]interface{}{
		"timestamp": t.Unix(),
		"msg_type":  msgType,
		"content":   content,
	}
	if wh.Signature != "" {
		payload["sign"] = b.genSignature(t, wh.Signature)
	}
	return payload
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

func trim(data []byte) []byte {
	s := strings.ReplaceAll(string(data), "\\\"", "\"")
	idx1 := strings.Index(s, "{")
	idx2 := strings.LastIndex(s, "}")
	return []byte(s[idx1 : idx2+1])
}
