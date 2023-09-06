package internal

import (
	"context"
	"strings"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"

	cdkgo "github.com/vanus-labs/cdk-go"
)

type Slack struct {
	logger     zerolog.Logger
	cfg        *slackConfig
	cancel     context.CancelFunc
	userID     string
	events     chan *cdkgo.Tuple
	cache      *cache.Cache
	eventTypes map[MessageType]struct{}
}

func NewSlack(cfg *slackConfig, logger zerolog.Logger, events chan *cdkgo.Tuple) *Slack {
	eventTypes := make(map[MessageType]struct{}, len(cfg.EventType))
	for _, t := range cfg.EventType {
		eventTypes[t] = struct{}{}
	}
	return &Slack{
		logger:     logger,
		cfg:        cfg,
		events:     events,
		eventTypes: eventTypes,
		cache:      cache.New(time.Minute*10, time.Minute*15),
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
	slackClient := socketmode.New(api)
	handler := socketmode.NewSocketmodeHandler(slackClient)
	handler.HandleEvents(slackevents.Message, d.messageEvent)
	handler.HandleDefault(d.defaultEvent)
	ctx, cancel := context.WithCancel(context.TODO())
	d.cancel = cancel
	go handler.RunEventLoopContext(ctx)
	return nil
}

func (d *Slack) containsEventType(msgType MessageType) bool {
	if len(d.eventTypes) == 0 {
		return true
	}
	_, exist := d.eventTypes[msgType]
	return exist
}

func (d *Slack) messageEvent(evt *socketmode.Event, client *socketmode.Client) {
	client.Ack(*evt.Request)
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return
	}
	ev, ok := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
	if !ok {
		return
	}
	_, exist := d.cache.Get(ev.ClientMsgID)
	if exist {
		d.logger.Info().Str("msgID", ev.ClientMsgID).
			Str("text", ev.Text).
			Msg("repeated receive msg")
		return
	}
	d.cache.SetDefault(ev.ClientMsgID, true)
	d.logger.Info().Str("msgID", ev.ClientMsgID).
		Str("text", ev.Text).
		Str("user", ev.User).
		Msg("receive msg")
	msgType := MessageText
	mentionUsers, content := d.parseText(ev.Text)
	if len(mentionUsers) > 0 {
		msgType = MessageTextAt
	}
	ed := &MessageData{
		Channel:      ev.Channel,
		ChannelType:  ev.ChannelType,
		BotID:        ev.BotID,
		User:         ev.User,
		MentionUsers: mentionUsers,
		Content:      content,
		Text:         ev.Text,
	}
	if ev.ThreadTimeStamp != "" && ev.ThreadTimeStamp != ev.EventTimeStamp {
		threadMsg := d.getThreadMsg(ev, client)
		ed.ThreadMessage = threadMsg
	}
	if ed.ThreadMessage != nil {
		msgType = MessageTextReply
	}
	if !d.containsEventType(msgType) {
		return
	}
	e := ce.NewEvent()
	e.SetID(ev.ClientMsgID)
	e.SetSource("vanus-slack-app")
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
	return
}

func (d *Slack) parseText(text string) ([]string, string) {
	var mentionUsers []string
	for {
		if !strings.HasPrefix(text, "<@") {
			break
		}
		index := strings.Index(text, "> ")
		if index <= 0 {
			break
		}
		mentionUsers = append(mentionUsers, text[2:index])
		text = text[index+2:]
	}
	return mentionUsers, strings.TrimSpace(text)
}

func (d *Slack) getThreadMsg(ev *slackevents.MessageEvent, client *socketmode.Client) *MessageData {
	messages, _, _, err := client.GetConversationReplies(&slack.GetConversationRepliesParameters{
		ChannelID: ev.Channel,
		Timestamp: ev.ThreadTimeStamp,
		Latest:    ev.ThreadTimeStamp,
		Oldest:    ev.ThreadTimeStamp,
		Inclusive: true,
	})
	if err != nil {
		d.logger.Warn().Err(err).Str("thread_ts", ev.ThreadTimeStamp).
			Msg("get conversation replies error")
		return nil
	}
	if len(messages) != 1 {
		d.logger.Info().Str("thread_ts", ev.ThreadTimeStamp).
			Int("length", len(messages)).
			Msg("resp message length not equal to 1")
	}
	msg := messages[0]
	mentionUsers, content := d.parseText(msg.Text)
	return &MessageData{
		User:         msg.User,
		BotID:        msg.BotID,
		MentionUsers: mentionUsers,
		Content:      content,
		Text:         msg.Text,
	}
}
func (d *Slack) defaultEvent(evt *socketmode.Event, client *socketmode.Client) {
	// Unexpected event type received
	d.logger.Info().Interface("eventType", evt.Type).Msg("unexpected event")
}
