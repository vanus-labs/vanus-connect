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
	"net/http"
	"os"

	ce "github.com/cloudevents/sdk-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

var _ cdkgo.Sink = &whatsappSink{}

func NewWhatsAppSink() cdkgo.Sink {
	return &whatsappSink{}
}

type whatsappSink struct {
	config *WhatsappConfig
	events chan *cdkgo.Tuple
	number int
	client *whatsmeow.Client
	logger zerolog.Logger
}

type JID struct {
	User   string
	Agent  uint8
	Device uint8
	Server string
	AD     bool
}

func (s *whatsappSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.logger = log.FromContext(ctx)
	s.config = cfg.(*WhatsappConfig)
	err := s.whatsappConnect()
	return err
}

func (s *whatsappSink) Name() string {
	return "whatsappSink"
}

func (s *whatsappSink) Destroy() error {
	if s.client != nil {
		s.client.Disconnect()
	}
	return nil
}

type Data struct {
	Info    types.MessageInfo `json:"info"`
	Message string            `json:"message"`
}

func (s *whatsappSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	for _, event := range events {
		result := s.processEvent(ctx, event)
		if cdkgo.SuccessResult != result {
			s.logger.Warn().Err(result.Error()).Str("event_id", event.ID()).Msg("event process failed")
			return result
		}
	}
	return cdkgo.SuccessResult
}

func (s *whatsappSink) processEvent(ctx context.Context, event *ce.Event) cdkgo.Result {
	var data Data
	err := json.Unmarshal(event.Data(), &data)
	if err != nil {
		return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
	}

	message := &waProto.Message{
		Conversation: proto.String(data.Message),
	}
	_, err = s.client.SendMessage(ctx, data.Info.Sender, message)
	if err != nil {
		return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
	}
	return cdkgo.SuccessResult
}

func (s *whatsappSink) whatsappConnect() error {
	dbFileName := "store.db"
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
