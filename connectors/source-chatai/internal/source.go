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
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/source/chatai/chat"
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

var _ cdkgo.HTTPSource = &chatSource{}

func NewChatSource() cdkgo.HTTPSource {
	return &chatSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type chatSource struct {
	config     *chatConfig
	events     chan *cdkgo.Tuple
	number     int
	day        string
	lock       sync.Mutex
	service    *chat.ChatService
	authEnable bool
}

func (s *chatSource) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*chatConfig)
	s.authEnable = !s.config.Auth.IsEmpty()
	s.service = chat.NewChatService(s.config.ChatConfig)
	return nil
}

func (s *chatSource) Name() string {
	return "ChatAiSource"
}

func (s *chatSource) Destroy() error {
	if s.service != nil {
		s.service.Close()
	}
	return nil
}

func (s *chatSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *chatSource) getMessage(req *http.Request) (map[string]interface{}, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return nil, errors.New("read body error")
	}
	message := make(map[string]interface{})
	contentType := req.Header.Get(headerContentType)
	if contentType != "" && strings.HasPrefix(contentType, applicationJSON) {
		err := json.Unmarshal(body, &message)
		if err != nil {
			return nil, errors.New("invalid JSON body")
		}
		if _, ok := message["message"].(string); !ok {
			return nil, errors.New("body message doesn't exist")
		}
	} else {
		message["message"] = string(body)
	}
	if message["message"] == "" {
		return nil, errors.New("message is empty")
	}
	return message, nil
}

func (s *chatSource) isSync(req *http.Request) bool {
	processMode := req.Header.Get(headerProcessMode)
	if processMode == "" {
		processMode = s.config.DefaultProcessMode
	}
	return processMode == "sync"
}

func (s *chatSource) writeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set(headerContentType, ce.ApplicationJSON)
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(`{"status":%d,"msg":"%s"}`, code, err.Error())))
}

func (s *chatSource) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.authEnable {
		username, password, ok := req.BasicAuth()
		if !ok || s.config.Auth.Username != username || s.config.Auth.Password != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="vanus connector", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	data, err := s.getMessage(req)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err)
		return
	}
	var chatType chat.Type
	chatMode := req.Header.Get(headerChatMode)
	if chatMode == "" {
		chatMode = req.Header.Get(headerChatModeOld)
	}
	if chatMode != "" {
		chatType = chat.Type(chatMode)
		switch chatType {
		case chat.ChatGPT, chat.ChatErnieBot:
		default:
			s.writeError(w, http.StatusBadRequest, errors.New("chat_mode invalid"))
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
	var wg sync.WaitGroup
	wg.Add(1)
	var userIdentifier string
	if s.config.UserIdentifierHeader != "" {
		userIdentifier = req.Header.Get(s.config.UserIdentifierHeader)
		if userIdentifier == "" {
			s.writeError(w, http.StatusBadRequest, errors.New("header userIdentifier is empty"))
			return
		}
	}
	sync := s.isSync(req)
	go func() {
		defer wg.Done()
		if !s.config.Stream {
			content, err := s.service.ChatCompletion(context.Background(), chatType, userIdentifier, data["message"].(string))
			if err != nil {
				log.Warning("failed to get content from Chat", map[string]interface{}{
					log.KeyError: err,
					"chatType":   chatType,
				})
			}
			data["result"] = content
			dataBytes := s.sendEvent(eventType, eventSource, data)
			if sync {
				w.Header().Set(headerContentType, ce.ApplicationJSON)
				w.WriteHeader(http.StatusOK)
				w.Write(dataBytes)
			}
		} else {
			stream, err := s.service.ChatCompletionStream(context.Background(), chatType, userIdentifier, data["message"].(string))
			if err != nil {
				log.Warning("failed to get chat with stream", map[string]interface{}{
					log.KeyError: err,
					"chatType":   chatType,
				})
				data["result"] = err.Error()
				s.sendEvent(eventType, eventSource, data)
				s.writeError(w, http.StatusInternalServerError, err)
				return
			}
			defer stream.Close()
			var flusher http.Flusher
			if sync {
				flusher, _ = w.(http.Flusher)
				w.Header().Set(headerContentType, "text/event-stream;charset=utf-8")
				w.WriteHeader(http.StatusOK)
			}
			for {
				msg, err := stream.Recv()
				if err != nil {
					data["result"] = err.Error()
					s.sendEvent(eventType, eventSource, data)
					if sync {
						w.Write([]byte(fmt.Sprintf(`{"status":%d,"msg":"%s"}`, http.StatusInternalServerError, err.Error())))
					}
					return
				}
				if msg == nil {
					return
				}
				data["is_end"] = msg.IsEnd
				data["result"] = msg.Content
				data["stream_id"] = msg.ID
				data["index"] = msg.Index
				if s.config.UserIdentifierHeader != "" {
					data[s.config.UserIdentifierHeader] = userIdentifier
				}
				dataBytes := s.sendEvent(eventType, eventSource, data)
				if sync {
					w.Write([]byte(fmt.Sprintf("data: %s\n\n", string(dataBytes))))
					flusher.Flush()
				}
			}
		}
	}()
	if !sync {
		w.Header().Set(headerContentType, ce.ApplicationJSON)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respSuccess))
		return
	}
	wg.Wait()

}

func (s *chatSource) sendEvent(eventType, eventSource string, data map[string]interface{}) []byte {
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetTime(time.Now())
	event.SetType(eventType)
	event.SetSource(eventSource)
	event.SetData(ce.ApplicationJSON, data)
	s.events <- &cdkgo.Tuple{
		Event: &event,
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
	return event.Data()
}

var (
	respSuccess = fmt.Sprintf(`{"status":%d,"msg":"%s"}`, 200, "Get API data successfully.")
)
