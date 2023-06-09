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

package bot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

type messageType string

const (
	xChatGroupID = "xvchatgroup"
	xMessageType = "xvmsgtype"

	textMessage       = messageType("text")
	linkMessage       = messageType("link")
	markdownMessage   = messageType("markdown")
	actionCardMessage = messageType("actionCard")
	feedCardMessage   = messageType("feedCard")
)

type Bot struct {
	defaultChatGroup string
	hookMap          map[string]WebHook
	httpClient       *resty.Client
	logger           zerolog.Logger
}

func NewBot(logger zerolog.Logger) *Bot {
	return &Bot{
		httpClient: resty.New(),
		logger:     logger,
	}
}

func (b *Bot) Init(cfg Config) error {
	b.hookMap = make(map[string]WebHook, len(cfg.Webhooks))
	if cfg.Default == "" {
		b.defaultChatGroup = cfg.Webhooks[0].ChatGroup
	} else {
		b.defaultChatGroup = cfg.Default
	}
	for i, wh := range cfg.Webhooks {
		b.hookMap[wh.ChatGroup] = cfg.Webhooks[i]
	}
	return nil
}

type botMessage struct {
	MsgType    messageType            `json:"msgtype"`
	Text       map[string]interface{} `json:"text,omitempty"`
	Link       map[string]interface{} `json:"link,omitempty"`
	MarkDown   map[string]interface{} `json:"markdown,omitempty"`
	ActionCard map[string]interface{} `json:"actionCard,omitempty"`
	FeedCard   map[string]interface{} `json:"feedCard,omitempty"`
}

func (b *Bot) SendMessage(e *v2.Event) (err error) {
	defer func() {
		if err != nil {
			d, _ := e.MarshalJSON()
			b.logger.Warn().Str("event", string(d)).Err(err).Msg("failed to send message")
		} else {
			b.logger.Info().Str("event_id", e.ID()).Msg("success send message")
		}
	}()
	chatGroup, ok := e.Extensions()[xChatGroupID].(string)
	if ok {
		_, exist := b.hookMap[chatGroup]
		if !exist {
			return fmt.Errorf("chat group %s is not exist", chatGroup)
		}
	} else {
		chatGroup = b.defaultChatGroup
	}
	msgType, ok := e.Extensions()[xMessageType].(string)
	if !ok {
		msgType = string(textMessage)
	}

	botMsg, err := b.event2Message(e, messageType(msgType))
	if err != nil {
		return err
	}
	return b.postMessage(botMsg, b.hookMap[chatGroup])
}

func (b *Bot) event2Message(e *v2.Event, msgType messageType) (*botMessage, error) {
	switch msgType {
	case textMessage:
		return &botMessage{MsgType: textMessage, Text: map[string]interface{}{
			"content": string(e.Data()),
		}}, nil
	default:
		var data map[string]interface{}
		err := json.Unmarshal(e.Data(), &data)
		if err != nil {
			return nil, err
		}
		botMsg := &botMessage{MsgType: msgType}
		switch msgType {
		case linkMessage:
			botMsg.Link = data
		case markdownMessage:
			botMsg.MarkDown = data
		case actionCardMessage:
			botMsg.ActionCard = data
		case feedCardMessage:
			botMsg.FeedCard = data
		default:
			return nil, fmt.Errorf("invalid message type:%s", msgType)
		}
		return botMsg, nil
	}
}

func (b *Bot) postMessage(botMsg *botMessage, hook WebHook) error {
	body, err := json.Marshal(botMsg)
	if err != nil {
		return err
	}
	t := time.Now().UnixMilli()
	sign := b.genSignature(t, hook.Signature)
	res, err := b.httpClient.R().
		SetQueryParam(paramTimestamp, strconv.FormatInt(t, 10)).
		SetQueryParam(paramSign, sign).
		SetHeader(http.ContentType, "application/json; charset=utf-8").
		SetBody(body).Post(hook.URL)
	if err != nil {
		return err
	}
	return b.processResponse(res)
}

func (b *Bot) genSignature(milli int64, secret string) string {
	strToSign := fmt.Sprintf("%d\n%s", milli, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(strToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (b *Bot) processResponse(res *resty.Response) error {
	obj := gjson.ParseBytes(res.Body())
	if obj.Get("errcode").Int() == 0 {
		return nil
	}
	return fmt.Errorf("failed to call dingtalk: %s", string(res.Body()))
}
