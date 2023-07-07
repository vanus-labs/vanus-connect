package internal

import (
	"context"
	"github.com/cloudevents/sdk-go/v2"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Sink = &DingtalkSink{}

type DingtalkSink struct {
	cfg    *Config
	logger zerolog.Logger
}

func NewSink() cdkgo.Sink {
	return &DingtalkSink{}
}

func (s *DingtalkSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) (err error) {
	s.logger = log.FromContext(ctx)
	s.cfg = cfg.(*Config)
	s.cfg.Init()

	return nil
}

func (s *DingtalkSink) Name() string {
	return "Dingtalk Sink"
}

func (s *DingtalkSink) Destroy() error {
	return nil
}

func (s *DingtalkSink) Arrived(ctx context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		e := events[idx]
		eStr, _ := e.MarshalJSON()
		content := gjson.Get(string(eStr), "data.content").String()
		webhook := gjson.Get(string(eStr), "data.webhook").String()
		s.logger.Info().
			Str("webhook", webhook).
			Str("content", content).
			Msg("Msg arrived")

		s.replyText(ctx, webhook, content)
	}
	return cdkgo.SuccessResult
}

func (s *DingtalkSink) replyText(ctx context.Context, webhook, content string) {
	replier := chatbot.NewChatbotReplier()
	err := replier.SimpleReplyText(ctx, webhook, []byte(content))
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed SimpleReplyText")
	}
}