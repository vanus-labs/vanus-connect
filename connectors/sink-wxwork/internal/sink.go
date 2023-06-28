package internal

import (
	"context"
	"github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/xen0n/go-workwx"
	"net/http"
)

var _ cdkgo.Sink = &WxworkSink{}

type WxworkSink struct {
	cfg    *Config
	logger zerolog.Logger

	workwxApp   *workwx.WorkwxApp
	httpHandler *workwx.HTTPHandler
}

func NewSink() cdkgo.Sink {
	return &WxworkSink{}
}

func (s *WxworkSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) (err error) {
	s.logger = log.FromContext(ctx)
	s.cfg = cfg.(*Config)
	s.cfg.Init()

	s.workwxApp = workwx.New(s.cfg.WeworkCorpId).WithApp(s.cfg.WeworkAgentSecret, s.cfg.WeworkAgentId)
	s.workwxApp.SpawnAccessTokenRefresherWithContext(ctx)

	return nil
}

func (s *WxworkSink) Name() string {
	return "Wxwork Sink"
}

func (s *WxworkSink) Destroy() error {
	return nil
}

func (s *WxworkSink) Arrived(_ context.Context, events ...*v2.Event) cdkgo.Result {
	for idx := range events {
		e := events[idx]
		eStr, _ := e.MarshalJSON()
		content := gjson.Get(string(eStr), "data.content").String()
		userId := gjson.Get(string(eStr), "data.fromUserID").String()
		s.logger.Info().
			Str("userId", userId).
			Str("content", content).
			Msg("Msg arrived")

		err := s.workwxApp.SendTextMessage(&workwx.Recipient{UserIDs: []string{userId}}, content, true)
		if err != nil {
			s.logger.Error().Err(err).Msg("Fail SendTextMessage")
			return cdkgo.NewResult(http.StatusInternalServerError, "Fail SendTextMessage")
		}
	}
	return cdkgo.SuccessResult
}
