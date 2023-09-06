package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"

	cdk "github.com/vanus-labs/cdk-go"
	cdkgo "github.com/vanus-labs/cdk-go"
)

type Feishu struct {
	logger     zerolog.Logger
	cli        *lark.Client
	cfg        *Config
	cache      *cache.Cache
	events     chan *cdk.Tuple
	eventTypes map[MessageType]struct{}
}

func NewFeishu(logger zerolog.Logger, cfg *Config, events chan *cdk.Tuple) *Feishu {
	eventTypes := make(map[MessageType]struct{}, len(cfg.EventType))
	for _, t := range cfg.EventType {
		eventTypes[t] = struct{}{}
	}
	return &Feishu{
		logger:     logger,
		cfg:        cfg,
		events:     events,
		eventTypes: eventTypes,
		cli:        lark.NewClient(cfg.AppID, cfg.AppSecret),
		cache:      cache.New(time.Minute*10, time.Minute*15),
	}
}

func (d *Feishu) containsMsgType(msgType MessageType) bool {
	if len(d.eventTypes) == 0 {
		return true
	}
	_, exist := d.eventTypes[msgType]
	return exist
}

func (d *Feishu) OnChatBotMessageReceived(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	message := event.Event.Message
	if *message.MessageType != "text" {
		return nil
	}
	eventID := event.EventV2Base.Header.EventID
	_, exist := d.cache.Get(eventID)
	if exist {
		d.logger.Info().Str("event", eventID).
			Str("content", *message.Content).Msg("repeated receive")
		return nil
	}
	d.cache.SetDefault(eventID, true)
	content, text := d.parseText(*message.Content)
	if text == "" {
		return nil
	}
	msgType := MessageText
	var mentionUsers []string
	if len(message.Mentions) > 0 {
		for _, m := range message.Mentions {
			mentionUsers = append(mentionUsers, *m.Id.OpenId)
		}
		msgType = MessageTextAt
	}
	ed := &MessageData{
		ChatID:       *message.ChatId,
		ChatType:     *message.ChatType,
		Content:      content,
		Text:         text,
		MentionUsers: mentionUsers,
		User:         *event.Event.Sender.SenderId.OpenId,
	}
	parentID := message.RootId
	if parentID == nil || *parentID == "" {
		parentID = message.ParentId
	}
	if parentID != nil && *parentID != "" {
		parentMsg := d.getParentMsg(ctx, *parentID)
		ed.ParentMessage = parentMsg
	}
	if ed.ParentMessage != nil {
		msgType = MessageTextReply
	}
	e := ce.NewEvent()
	e.SetID(eventID)
	e.SetSource("vanus-feishu-app")
	e.SetType(string(msgType))
	e.SetData(ce.ApplicationJSON, ed)
	d.events <- &cdkgo.Tuple{
		Event: &e,
		Success: func() {
			d.logger.Info().Msg("send event to target success")
		},
		Failed: func(err2 error) {
			d.logger.Warn().Err(err2).Msg("failed to send event to target")
		},
	}
	return nil
}

func (d *Feishu) parseText(content string) (string, string) {
	var textMsg TextMsg
	err := json.Unmarshal([]byte(content), &textMsg)
	if err != nil {
		d.logger.Warn().Err(err).
			Str("content", content).
			Msg("unmarshal content error")
		return "", ""
	}
	text := textMsg.Text
	for i := 1; ; i++ {
		t := fmt.Sprintf("@_user_%d ", i)
		if !strings.HasPrefix(text, t) {
			break
		}
		text = text[len(t):]
	}
	return strings.TrimSpace(text), textMsg.Text
}

func (d *Feishu) getParentMsg(ctx context.Context, msgID string) *MessageData {
	req := larkim.NewGetMessageReqBuilder().MessageId(msgID).Build()
	resp, err := d.cli.Im.Message.Get(ctx, req)
	if err != nil {
		d.logger.Warn().Err(err).Str("msg_id", msgID).
			Msg("get message error")
		return nil
	}
	if !resp.Success() {
		d.logger.Warn().Str("msg_id", msgID).
			Str("resp", resp.Error()).
			Msg("get message resp error")
		return nil
	}
	if len(resp.Data.Items) != 1 {
		d.logger.Info().Str("msg_id", msgID).
			Int("length", len(resp.Data.Items)).
			Msg("resp message length not equal to 1")
	}
	msg := resp.Data.Items[0]
	if *msg.MsgType != "text" {
		return nil
	}
	content, text := d.parseText(*msg.Body.Content)
	if text == "" {
		return nil
	}
	var mentionUsers []string
	for _, m := range msg.Mentions {
		mentionUsers = append(mentionUsers, *m.Id)
	}
	return &MessageData{
		Content:      content,
		MentionUsers: mentionUsers,
		Text:         text,
		User:         *msg.Sender.Id,
	}
}
