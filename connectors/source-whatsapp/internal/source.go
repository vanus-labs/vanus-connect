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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	ce "github.com/cloudevents/sdk-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/source/chatai/chat"
)

var _ cdkgo.Source = &whatsAppSource{}

func NewWhatsAppSource() cdkgo.Source {
	return &whatsAppSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type whatsAppSource struct {
	config      *whatsAppConfig
	events      chan *cdkgo.Tuple
	client      *whatsmeow.Client
	chatService *chat.ChatService
	logger      zerolog.Logger
}

func (s *whatsAppSource) getDBFileName(cfg *whatsAppConfig) string {
	if cfg.FileName != "" {
		return cfg.FileName
	}
	if cfg.WhatsAppID != "" {
		return cfg.WhatsAppID + ".db"
	}
	return "store.db"
}

func (s *whatsAppSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*whatsAppConfig)
	dbFileName := s.getDBFileName(s.config)
	if s.config.Data != "" {
		dbBytes, err := base64.StdEncoding.DecodeString(s.config.Data)
		if err != nil {
			return err
		}
		err = os.WriteFile(dbFileName, dbBytes, 0644)
		if err != nil {
			return err
		}
		s.logger.Info().Msg("Database restored successfully.")
	}

	if s.config.EnableChatAi {
		s.chatService = chat.NewChatService(*s.config.ChatConfig, s.logger)
	}

	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFileName), dbLog)
	if err != nil {
		return err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return err
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	s.client = whatsmeow.NewClient(deviceStore, clientLog)

	s.events = make(chan *cdkgo.Tuple, 100)
	s.client.AddEventHandler(func(evt interface{}) {

		switch v := evt.(type) {
		case *events.PairSuccess:
			log.Info().Str("jid", v.ID.String()).Msg("pair success")
		case *events.Message:

			info := v.Info
			message := v.Message.GetConversation()
			if message == "" {
				// androd
				message = v.Message.GetExtendedTextMessage().GetText()
			}
			if message != "" {
				if v.Info.IsFromMe {
					if v.Info.Sender.User == v.Info.Chat.User {
						event := s.makeEvent(info, message)
						s.events <- &cdkgo.Tuple{
							Event: event,
							Success: func() {
								s.logger.Info().Str("event_id", event.ID()).Msg("send event success")
							},
							Failed: func(err error) {
								// TODO
								b, _ := json.Marshal(event)
								s.logger.Warn().Err(err).Msg("send event failed: " + string(b))
							},
						}
					}

				} else {

					event := s.makeEvent(info, message)
					s.events <- &cdkgo.Tuple{
						Event: event,
						Success: func() {
							s.logger.Info().Str("event_id", event.ID()).Msg("send event success")
						},
						Failed: func(err error) {
							// TODO
							b, _ := json.Marshal(event)
							s.logger.Warn().Err(err).Msg("send event failed: " + string(b))
						},
					}
				}
			}
		}

	})

	if s.client.Store.ID == nil {
		if s.config.Data != "" {
			return fmt.Errorf("data exist but store id is nil")
		}
		qrChan, _ := s.client.GetQRChannel(context.Background())
		err = s.client.Connect()
		if err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qr.png")

				if err != nil {
					fmt.Println("Failed to generate QR code:", err)
				} else {
					fmt.Println("QR code generated successfully")
				}

			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = s.client.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *whatsAppSource) Name() string {
	return "WhatsAppSource"
}

func (s *whatsAppSource) Destroy() error {
	if s.client != nil {
		s.client.Disconnect()
	}
	if s.chatService != nil {
		s.chatService.Close()
	}
	return nil
}

func (s *whatsAppSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *whatsAppSource) makeEvent(info types.MessageInfo, message string) *ce.Event {
	event := ce.NewEvent()
	event.SetID(info.ID)
	event.SetSource("whatsapp")
	event.SetType(info.Type)
	if s.chatService != nil {
		resp, err := s.chatService.ChatCompletion(context.Background(), s.config.ChatConfig.DefaultChatMode, info.Sender.User, message)
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to get content from Chat")
		}
		event.SetData(ce.ApplicationJSON, map[string]interface{}{
			"info":    info,
			"message": resp,
		})
	} else {
		event.SetData(ce.ApplicationJSON, map[string]interface{}{
			"info":    info,
			"message": message,
		})
	}
	return &event
}
