package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"

	ce "github.com/cloudevents/sdk-go/v2"

	cdkgo "github.com/vanus-labs/cdk-go"
)

const (
	headerTopic            = "X-Shopify-Topic"
	extendAttributesTopic  = "xvshopifytopic"
	extendAttributesAction = "xvaction"
)

func (s *shopifySource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.shopifyApp.VerifyWebhookRequest(r) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("hmac invalid"))
		return
	}
	atomic.AddInt64(&s.count, 1)
	s.logger.Info().Int64("total", atomic.LoadInt64(&s.count)).Msg("receive a new event")
	topic := r.Header.Get(headerTopic)
	//  "products/create"
	//	"products/delete"
	//	"products/update"
	topicArr := strings.Split(topic, "/")
	if len(topicArr) != 2 {
		s.logger.Info().Str("topic", topic).Msg("event topic invalid")
		return
	}
	if !s.isSync(ApiType(topicArr[0])) {
		s.logger.Info().Str("topic", topic).Msg("skip event")
		return
	}
	e := s.newEvent()
	e.SetType(topicArr[0])
	e.SetExtension(extendAttributesAction, topicArr[1])
	e.SetExtension(extendAttributesTopic, topic)
	body, _ := io.ReadAll(r.Body)
	var m map[string]interface{}
	_ = json.Unmarshal(body, &m)
	e.SetData(ce.ApplicationJSON, m)
	wg := sync.WaitGroup{}
	wg.Add(1)
	s.events <- &cdkgo.Tuple{
		Event: &e,
		Success: func() {
			defer wg.Done()
			s.logger.Info().Str("event_id", e.ID()).Msg("send event to target success")
			w.WriteHeader(http.StatusOK)
		},
		Failed: func(err2 error) {
			defer wg.Done()
			s.logger.Warn().Interface("event_id", e.ID()).Err(err2).Msg("failed to send event to target")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(
				fmt.Sprintf("failed to send event to remote server: %s", err2.Error())))
		},
	}
	wg.Wait()

}
