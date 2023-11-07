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
	"net/http"
	"strings"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"

	cdk "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/connector"
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
	errChatGroup = cdk.NewResult(http.StatusBadRequest, "feishu: xvfeishuchatgroup is missing or incorrect, please check your Feishu "+
		"Sink config or subscription")
	errMessageType = cdk.NewResult(http.StatusBadRequest, "feishu: xvfeishumsgtype is missing or invalid, only"+
		" [text, post, share_chat, image, interactive] are supported")
	errInvalidPostMessage     = errors.New("feishu: invalid post message, please make sure it's the json format")
	errInvalidAttributes      = cdk.NewResult(http.StatusBadRequest, "feishu: invalid xvfeishuboturls or xvfeishubotsigs")
	errInvalidAttributeNumber = cdk.NewResult(http.StatusBadRequest, "feishu: the number of bot url and signature must be equal")
	errNoBotWebhookFound      = cdk.NewResult(http.StatusBadRequest, "feishu: no feishu bot target webhook found")
)

type bot struct {
	cfg        BotConfig
	cm         map[string]WebHook
	httpClient *resty.Client
	logger     zerolog.Logger
}

func (b *bot) init(cfg BotConfig, logger zerolog.Logger) error {
	b.logger = logger
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
	Default      string    `json:"default" yaml:"default"`
	DynamicRoute bool      `json:"dynamic_route" yaml:"dynamic_route"`
}

func (c *BotConfig) Validate() error {
	if !c.DynamicRoute && len(c.Webhooks) == 0 {
		return errors.New("the bot.webhooks can't be empty when dynamic_route is false")
	}
	if len(c.Webhooks) > 0 {
		if c.Default == "" {
			c.Default = c.Webhooks[0].ChatGroup
		} else {
			for _, webhook := range c.Webhooks {
				if webhook.ChatGroup == c.Default {
					return nil
				}
			}
			return errors.New("the bot.default not exist in webhooks.chatGroup")
		}
	}
	return nil
}

func (b *bot) sendMessage(e *v2.Event) cdk.Result {
	var (
		whs     []WebHook
		groupID string
	)
	groupID, ok := e.Extensions()[xChatGroupID].(string)
	if ok {
		wh, exist := b.cm[groupID]
		if !exist {
			if !b.cfg.DynamicRoute {
				return errChatGroup
			}
		} else {
			whs = append(whs, wh)
		}
	} else {
		if !b.cfg.DynamicRoute {
			whs = append(whs, b.cm[b.cfg.Default])
		}
	}

	if b.cfg.DynamicRoute {
		urlAttr, ok := e.Extensions()[xBotURL].(string)
		if !ok {
			return errInvalidAttributes
		}
		signatureAttr, ok := e.Extensions()[xBotSignature].(string)
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

	t, ok := e.Extensions()[xMessageType].(string)
	if !ok {
		t = string(textMessage)
	}
	botMsg, err := b.event2BotMessage(e, messageType(t))
	if err != nil {
		return cdk.NewResult(http.StatusBadRequest, "event parse error:"+err.Error())
	}
	now := time.Now().Unix()
	for _, wh := range whs {
		if wh.Signature != "" {
			botMsg.Timestamp = &now
			botMsg.Sign = b.genSignature(now, wh.Signature)
		} else {
			botMsg.Timestamp = nil
			botMsg.Sign = ""
		}
		res, err := b.httpClient.R().SetBody(botMsg).Post(wh.URL)
		if err != nil {
			return cdk.NewResult(http.StatusInternalServerError, "call feishu error: "+err.Error())
		}
		if code, err := b.processResponse(e, res); err != nil {
			return cdk.NewResult(connector.Code(code), "call feishu response error:"+err.Error())
		}
	}
	return cdk.SuccessResult
}

func (b *bot) event2BotMessage(e *v2.Event, msgType messageType) (*botMessage, error) {
	msg := &botMessage{
		MsgType: msgType,
	}
	switch msgType {
	case textMessage:
		var text string
		if isJSONString(e) {
			err := json.Unmarshal(e.Data(), &text)
			if err != nil {
				return nil, err
			}
		} else {
			text = string(e.Data())
		}
		msg.Content = &botContent{
			Text: text,
		}
	case postMessage:
		m := map[string]interface{}{}
		if err := json.Unmarshal(trim(e.Data()), &m); err != nil {
			return nil, err
		}
		msg.Content = &botContent{
			Post: m,
		}
	case shareChatMessage:
		msg.Content = &botContent{
			ShareChatID: string(e.Data()),
		}
	case imageMessage:
		msg.Content = &botContent{
			ImageKey: string(e.Data()),
		}
	case interactiveMessage:
		m := map[string]interface{}{}
		if err := json.Unmarshal(trim(e.Data()), &m); err != nil {
			return nil, err
		}
		msg.Card = m
	default:
		return nil, errMessageType.Error()
	}
	return msg, nil
}

func isJSONString(e *v2.Event) bool {
	if e.DataContentType() != v2.ApplicationJSON {
		return false
	}
	for i := range e.Data() {
		c := e.Data()[i]
		switch c {
		case '"':
			return true
		case '\t', ' ':
			continue
		default:
			return false
		}
	}
	return false
}

func (b *bot) genSignature(timestamp int64, signature string) string {
	if signature == "" {
		return ""
	}
	strToSign := fmt.Sprintf("%d\n%s", timestamp, signature)
	h := hmac.New(sha256.New, []byte(strToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type botContent struct {
	Text        string                 `json:"text,omitempty"`
	Post        map[string]interface{} `json:"post,omitempty"`
	ShareChatID string                 `json:"share_chat_id,omitempty"`
	ImageKey    string                 `json:"image_key,omitempty"`
}

type botMessage struct {
	MsgType messageType            `json:"msg_type"`
	Content *botContent            `json:"content,omitempty"`
	Card    map[string]interface{} `json:"card,omitempty"`

	Sign      string `json:"sign,omitempty"`
	Timestamp *int64 `json:"timestamp,omitempty"`
}

type botResponse struct {
	Code int `json:"code"`
	//Msg  string `json:"msg"`
}

func (b *bot) processResponse(e *v2.Event, res *resty.Response) (int, error) {
	// docs: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN?lang=zh-CN#756b882f
	var resp botResponse
	err := json.Unmarshal(res.Body(), &resp)
	if err != nil {
		b.logger.Info().Err(err).Str("event_id", e.ID()).
			Str("body", string(res.Body())).Msg("unmarshal error")
		return http.StatusBadRequest, err
	}
	if resp.Code == 0 {
		b.logger.Info().Str("event_id", e.ID()).Msg("success send message to feishu Bot")
		return 0, nil
	}
	var code int
	switch resp.Code {
	case 9499:
		code = http.StatusTooManyRequests
	default:
		code = http.StatusBadRequest
	}
	return code, fmt.Errorf("failed to call feishu: %s", string(res.Body()))
}

func trim(data []byte) []byte {
	s := strings.ReplaceAll(string(data), "\\\"", "\"")
	idx1 := strings.Index(s, "{")
	idx2 := strings.LastIndex(s, "}")
	return []byte(s[idx1 : idx2+1])
}
