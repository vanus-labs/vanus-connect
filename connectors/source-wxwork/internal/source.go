package internal

import (
	"context"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/zerolog"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/xen0n/go-workwx"
	"net/http"
	"time"
)

var _ cdkgo.HTTPSource = &WxworkSource{}

type WxworkSource struct {
	config *Config
	logger zerolog.Logger
	events chan *cdkgo.Tuple

	cache *ttlcache.Cache[string, any]

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
	s.cache = ttlcache.New[string, any](
		ttlcache.WithTTL[string, any](time.Minute),
	)
	go s.cache.Start()
	s.config = cfg.(*Config)
	s.config.Init()

	s.workwxApp = workwx.New(s.config.WeworkCorpId).WithApp(s.config.WeworkAgentSecret, s.config.WeworkAgentId)
	s.workwxApp.SpawnAccessTokenRefresherWithContext(ctx)

	s.httpHandler, err = workwx.NewHTTPHandler(s.config.WeworkToken, s.config.WeworkEncodingAESKey, &WxworkMessageHandler{s})
	if err != nil {
		s.logger.Error().Err(err).Msg("workwx.NewHTTPHandler fail")
		return err
	}

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
	if s.cache.Get(req.RequestURI) != nil {
		s.logger.Info().
			Str("RequestURI", req.RequestURI).
			Msg("Duplicated request come, just return")
		return
	}
	s.cache.Set(req.RequestURI, struct{}{}, time.Minute)
	s.logger.Info().
		Str("RequestURI", req.RequestURI).
		Msg("Request first come")
	s.httpHandler.ServeHTTP(w, req)
}

func (s *WxworkSource) sendEvent(data map[string]interface{}) []byte {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetTime(time.Now())
	event.SetType("Conversion")
	event.SetSource(s.Name())
	event.SetData(ce.ApplicationJSON, data)
	s.events <- &cdkgo.Tuple{
		Event: &event,
	}
	return event.Data()
}
