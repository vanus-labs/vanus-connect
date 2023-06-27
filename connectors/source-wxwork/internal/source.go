package internal

import (
	"context"
	"github.com/rs/zerolog"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/xen0n/go-workwx"
	"net/http"
)

var _ cdkgo.HTTPSource = &WxworkSource{}

type WxworkSource struct {
	config *Config
	logger zerolog.Logger
	events chan *cdkgo.Tuple

	workwxApp   *workwx.WorkwxApp
	httpHandler *workwx.HTTPHandler
}

func NewSource() cdkgo.HTTPSource {
	return &WxworkSource{
		events: make(chan *cdkgo.Tuple, 1024),
	}
}

func (s *WxworkSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) (err error) {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*Config)
	s.config.Init()

	s.workwxApp = workwx.New(s.config.WeworkCorpId).WithApp(s.config.WeworkAgentSecret, s.config.WeworkAgentId)

	s.httpHandler, err = workwx.NewHTTPHandler(s.config.WeworkToken, s.config.WeworkEncodingAESKey, WxworkMessageHandler{})
	if err != nil {
		s.logger.Error().Err(err).Msg("workwx.NewHTTPHandler fail")
		return err
	}

	println("Initialize success")
	return nil
}

func (s *WxworkSource) Name() string {
	return "Wxwork Source"
}

func (s *WxworkSource) Destroy() error {
	return nil
}

func (s *WxworkSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *WxworkSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.httpHandler.ServeHTTP(w, req)
}
