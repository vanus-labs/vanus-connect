// Copyright 2023 Linkall Inc.
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

package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	defaultEventType   = "vanus-chatAI-type"
	defaultEventSource = "vanus-chatAI-source"
	headerSource       = "Vanus-Source"
	headerType         = "Vanus-Type"
	headerContentType  = "Content-Type"
	headerChatModeOld  = "Chat_Mode"
	headerChatMode     = "Chat-Mode"
	headerProcessMode  = "Process-Mode"
	applicationJSON    = "application/json"
)

var _ cdkgo.Source = &chatSource{}

func NewChatSource() cdkgo.Source {
	return &chatSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type chatSource struct {
	config  *chatConfig
	events  chan *cdkgo.Tuple
	number  int
	day     string
	lock    sync.Mutex
	server  *http.Server
	service *chatService
}

func (s *chatSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*chatConfig)
	s.config.Init()
	s.service = newChatService(s.config)
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s,
	}
	go func() {
		log.Info("http server is ready to start", map[string]interface{}{
			"port": s.config.Port,
		})
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("cloud not listen on %d, error:%s", s.config.Port, err.Error()))
		}
		log.Info("http server stopped", nil)
	}()
	return nil
}

func (s *chatSource) Name() string {
	return "ChatAiSource"
}

func (s *chatSource) Destroy() error {
	if s.server != nil {
		s.server.Shutdown(context.Background())
	}
	return nil
}

func (s *chatSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *chatSource) getMessage(req *http.Request) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return nil, errors.New("read body error")
	}
	message := make(map[string]interface{})
	contentType := req.Header.Get(headerContentType)
	if contentType != "" && strings.HasPrefix(contentType, applicationJSON) {
		err = json.Unmarshal(body, &message)
		if err != nil {
			return nil, errors.New("invalid JSON body")
		}
		if _, ok := message["message"].(string); !ok {
			return nil, errors.New("body message doesn't exist")
		}
	} else {
		message["message"] = string(body)
	}
	return message, nil
}

func (s *chatSource) isSync(req *http.Request) bool {
	processMode := req.Header.Get(headerProcessMode)
	return processMode == "sync"
}

func (s *chatSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data, err := s.getMessage(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"status":%d,"msg":"%s"}`, http.StatusBadRequest, err.Error())))
		return
	}
	chatType := chatGPT
	chatMode := req.Header.Get(headerChatMode)
	if chatMode == "" {
		chatMode = req.Header.Get(headerChatModeOld)
	}
	if chatMode != "" {
		chatType = ChatType(chatMode)
		switch chatType {
		case chatGPT, chatErnieBot:
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(`{"status":%d,"msg":"%s"}`, http.StatusBadRequest, "chat_mode invalid")))
			return
		}
	}
	eventSource := req.Header.Get(headerSource)
	if eventSource == "" {
		eventSource = defaultEventSource
	}
	eventType := req.Header.Get(headerType)
	if eventType == "" {
		eventType = defaultEventType
	}
	event := ce.NewEvent()
	event.SetID(uuid.New().String())
	event.SetTime(time.Now())
	event.SetType(eventType)
	event.SetSource(eventSource)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(event *ce.Event, data map[string]interface{}) {
		defer wg.Done()
		content, err := s.service.ChatCompletion(chatType, data["message"].(string))
		if err != nil {
			log.Warning("failed to get content from Chat", map[string]interface{}{
				log.KeyError: err,
				"chatType":   chatType,
			})
		}
		data["result"] = content
		event.SetData(ce.ApplicationJSON, data)
		s.events <- &cdkgo.Tuple{
			Event: event,
			Success: func() {
				log.Info("send event to target success", map[string]interface{}{
					"event": event.ID(),
				})
			},
			Failed: func(err2 error) {
				log.Warning("failed to send event to target", map[string]interface{}{
					log.KeyError: err2,
					"event":      event.ID(),
				})
			},
		}
	}(&event, data)
	if !s.isSync(req) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respSuccess))
		return
	}
	wg.Wait()
	w.WriteHeader(http.StatusOK)
	w.Write(event.Data())
}

var (
	respSuccess = fmt.Sprintf(`{"status":%d,"msg":"%s"}`, 200, "Get API data successfully.")
)
