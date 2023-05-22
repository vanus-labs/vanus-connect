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
	"fmt"
	ce "github.com/cloudevents/sdk-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	cdkgo "github.com/vanus-labs/cdk-go"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
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
}

type JID struct {
	User   string
	Agent  uint8
	Device uint8
	Server string
	AD     bool
}

func (s *whatsappSink) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	// TODO\
	s.config = cfg.(*WhatsappConfig)
	s.whatsappConnect()
	return nil
}

func (s *whatsappSink) Name() string {
	// TODO
	return "whatsappSink"
}

func (s *whatsappSink) Destroy() error {
	// TODO
	s.client.Disconnect()
	return nil
}

type Data struct {
	Info    types.MessageInfo `json:"info"`
	Message string            `json:"message"`
}

func (s *whatsappSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	// TODO
	for _, event := range events {
		var data Data
		_ = json.Unmarshal(event.Data(), &data)

		message := &waProto.Message{
			Conversation: proto.String(data.Message),
		}
		s.client.SendMessage(ctx, data.Info.Sender, message)

	}
	return cdkgo.SuccessResult
}

func (s *whatsappSink) whatsappConnect() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:Store.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	s.client = whatsmeow.NewClient(deviceStore, clientLog)

	if s.client.Store.ID == nil {
		qrChan, _ := s.client.GetQRChannel(context.Background())
		err = s.client.Connect()
		if err != nil {
			panic(err)
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
			panic(err)
		}

	}

}
