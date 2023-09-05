package internal

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	cdk "github.com/vanus-labs/cdk-go"
	cdkgo "github.com/vanus-labs/cdk-go"
)

type Feishu struct {
	logger zerolog.Logger
	cli    *lark.Client
	cfg    *Config
	cache  *cache.Cache
	events chan *cdk.Tuple
}

func NewFeishu(logger zerolog.Logger, cfg *Config, events chan *cdk.Tuple) *Feishu {
	return &Feishu{
		logger: logger,
		cfg:    cfg,
		events: events,
		cli:    lark.NewClient(cfg.AppID, cfg.AppSecret),
		cache:  cache.New(time.Minute*10, time.Minute*15),
	}
}

func (d *Feishu) OnChatBotMessageReceived(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	message := event.Event.Message
	if *message.MessageType != "text" {
		return nil
	}
	if *message.ChatType != "group" {
		return nil
	}
	if d.cfg.UserID != "" && *event.Event.Sender.SenderId.UserId != d.cfg.UserID {
		//not cfg user message
		return nil
	}
	if message.ParentId == nil || *message.ParentId == "" {
		//not reply message
		return nil
	}
	if message.RootId != nil && *message.RootId != *message.ParentId {
		// reply multi message
		return nil
	}
	if len(message.Mentions) != 1 {
		//not @ msg
		return nil
	}
	eventData, err := d.getParentMsgContent(ctx, *message.ParentId)
	if err != nil {
		return err
	}
	if eventData == nil {
		return nil
	}
	var text TextMsg
	err = json.Unmarshal([]byte(*message.Content), &text)
	if err != nil {
		d.logger.Info().Err(err).Str("content", *message.Content).Msg("unmarshal content error")
		return err
	}
	answer := text.Text
	if len(answer) < 9 {
		return nil
	}
	answer = strings.TrimSpace(answer[9:])
	if answer == "" {
		return nil
	}
	eventID := event.EventV2Base.Header.EventID
	_, exist := d.cache.Get(eventID)
	if exist {
		d.logger.Info().Str("event", eventID).Msg("repeated receive")
		return nil
	}
	d.logger.Info().
		Str("msgID", *message.MessageId).
		Str("answer", answer).
		Str("question", eventData.Question).
		Msg("event receive")
	d.cache.SetDefault(eventID, true)
	eventData.Answer = answer
	eventData.AnswerUser = *event.Event.Sender.SenderId.UserId
	e := ce.NewEvent()
	e.SetID(eventID)
	e.SetSource("vanus-feishu-app")
	e.SetType("reply")
	e.SetData(ce.ApplicationJSON, eventData)
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

func (d *Feishu) getParentMsgContent(ctx context.Context, msgID string) (*EventData, error) {
	req := larkim.NewGetMessageReqBuilder().MessageId(msgID).Build()
	resp, err := d.cli.Im.Message.Get(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "get message error")
	}
	if !resp.Success() {
		return nil, errors.Errorf("resp error,code:%d,msg:%s", resp.StatusCode, resp.Msg)
	}
	if len(resp.Data.Items) != 1 {
		return nil, nil
	}
	msg := resp.Data.Items[0]
	if *msg.MsgType != "text" {
		return nil, nil
	}
	if len(msg.Mentions) != 1 {
		//not @ msg
		return nil, nil
	}
	//todo @ is bot
	var text TextMsg
	err = json.Unmarshal([]byte(*msg.Body.Content), &text)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal content %s error", *msg.Body.Content)
	}
	prompt := text.Text
	//format @_user_1 hi"
	if len(prompt) < 9 {
		return nil, nil
	}
	prompt = strings.TrimSpace(prompt[9:])
	if prompt == "" {
		return nil, nil
	}
	return &EventData{
		Question:     prompt,
		QuestionUser: *msg.Mentions[0].Id,
	}, nil
}

type TextMsg struct {
	Text string `json:"text"`
}

type EventData struct {
	QuestionUser string `json:"question_user"`
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	AnswerUser   string `json:"answer_user"`
}
