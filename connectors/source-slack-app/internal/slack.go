package internal

import (
	"context"
	"fmt"
	"strings"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	cdkgo "github.com/vanus-labs/cdk-go"
)

type Slack struct {
	logger            zerolog.Logger
	cfg               *slackConfig
	cancel            context.CancelFunc
	groupMsgPrefix    string
	groupMsgPrefixLen int
	userID            string
	events            chan *cdkgo.Tuple
}

func NewSlack(cfg *slackConfig, logger zerolog.Logger, events chan *cdkgo.Tuple) *Slack {
	return &Slack{
		logger: logger,
		cfg:    cfg,
		events: events,
	}
}

func (d *Slack) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
}

func (d *Slack) Start() error {
	api := slack.New(d.cfg.BotToken, slack.OptionAppLevelToken(d.cfg.AppToken))
	resp, err := api.AuthTest()
	if err != nil {
		return err
	}
	// userID fmt: U04M1L7L64U
	d.userID = resp.UserID
	d.groupMsgPrefix = fmt.Sprintf("<@%s> ", d.userID)
	d.groupMsgPrefixLen = len(d.groupMsgPrefix)
	slackClient := socketmode.New(api)
	handler := socketmode.NewSocketmodeHandler(slackClient)
	handler.HandleEvents(slackevents.Message, d.directMessageEvent)
	handler.HandleDefault(d.defaultEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	d.cancel = cancel
	go handler.RunEventLoopContext(ctx)
	return nil
}

func (d *Slack) directMessageEvent(evt *socketmode.Event, client *socketmode.Client) {
	client.Ack(*evt.Request)
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return
	}
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		return
	}
	if ev.BotID != "" {
		// bot message
		return
	}
	if ev.ChannelType != "channel" {
		// not channel message
		return
	}
	if d.cfg.UserID != "" && ev.User != d.cfg.UserID {
		// not use msg
		return
	}
	if ev.ThreadTimeStamp == "" || ev.ThreadTimeStamp == ev.EventTimeStamp {
		// not reply msg
		return
	}
	eventData, err := d.getTsMsgContent(ev, client)
	if err != nil {
		d.logger.Warn().Err(err).Str("thread_ts", ev.ThreadTimeStamp).Msg("get ts msg error")
		return
	}
	if eventData == nil {
		d.logger.Info().Str("thread_ts", ev.ThreadTimeStamp).Msg("ts msg is nil")
		return
	}
	answer := ev.Text
	eventData.Answer = answer
	eventData.AnswerUser = ev.User
	d.logger.Info().
		Str("thread_ts", ev.ThreadTimeStamp).
		Str("answer", answer).
		Str("question", eventData.Question).
		Msg("event receive")
	e := ce.NewEvent()
	e.SetID(ev.ClientMsgID)
	e.SetSource("vanus-slack-app")
	e.SetType("question-answer")
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
	return
}

func (d *Slack) getTsMsgContent(ev *slackevents.MessageEvent, client *socketmode.Client) (*EventData, error) {
	resp, err := client.GetConversationHistory(&slack.GetConversationHistoryParameters{
		Latest:    ev.ThreadTimeStamp,
		ChannelID: ev.Channel,
		Oldest:    ev.ThreadTimeStamp,
		Inclusive: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get conversation history error")
	}
	if len(resp.Messages) != 1 {
		return nil, nil
	}
	msg := resp.Messages[0]
	texts := strings.SplitN(msg.Text, " ", 2)
	if len(texts) != 2 {
		// not @ msg
		return nil, nil
	}
	if !strings.HasPrefix(texts[0], "<@") || !strings.HasSuffix(texts[0], ">") {
		return nil, nil
	}
	text := strings.TrimSpace(texts[1])
	if text == "" {
		return nil, nil
	}
	if strings.HasPrefix(text, "<@") {
		//@ other user
		return nil, nil
	}
	botUser := texts[0][2 : len(texts[0])-1]
	return &EventData{
		Question:       text,
		QuestionUser:   msg.User,
		QuestionAtUser: botUser,
	}, nil
}
func (d *Slack) defaultEvent(evt *socketmode.Event, client *socketmode.Client) {
	// Unexpected event type received
	d.logger.Info().Interface("eventType", evt.Type).Msg("unexpected event")
}
