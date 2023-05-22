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
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"math/rand"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	cdkgo "github.com/vanus-labs/cdk-go"
	"go.mau.fi/whatsmeow/types/events"
)

var _ cdkgo.Source = &whatsAppSource{}

func NewExampleSource() cdkgo.Source {
	return &whatsAppSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type whatsAppSource struct {
	config *whatsAppConfig
	events chan *cdkgo.Tuple
	number int
	client *whatsmeow.Client
}

func (s *whatsAppSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {

	s.config = cfg.(*whatsAppConfig)

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:exampleStore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	s.client = whatsmeow.NewClient(deviceStore, clientLog)

	s.events = make(chan *cdkgo.Tuple, 100)
	s.client.AddEventHandler(func(evt interface{}) {

		switch v := evt.(type) {
		case *events.Message:
			message := v.Message.GetConversation()
			number := v.Info.Sender.User

			if message != "" {
				event := s.makeEvent(message, number)
				s.events <- &cdkgo.Tuple{
					Event: event,
					Success: func() {
						// TODO
						b, _ := json.Marshal(event)
						fmt.Println("send event success: " + string(b))
					},
					Failed: func(err error) {
						// TODO
						b, _ := json.Marshal(event)
						fmt.Println("send event failed: " + string(b) + ", error: " + err.Error())
					},
				}
			}
		}
	})

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

	return nil
}

func (s *whatsAppSource) Name() string {
	// TODO
	return "WhatsAppSource"
}

func (s *whatsAppSource) Destroy() error {
	s.client.Disconnect()
	return nil
}

func (s *whatsAppSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *whatsAppSource) makeEvent(str string, num string) *ce.Event {
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+100))
	s.number++
	event := ce.NewEvent()
	event.SetID(fmt.Sprintf("id-%d", s.number))
	event.SetSource("whatsAppSource")
	event.SetType("testType")
	event.SetExtension("t", time.Now())
	event.SetData(ce.ApplicationJSON, map[string]interface{}{
		"number":  fmt.Sprintf(num),
		"message": fmt.Sprintf(str),
	})
	return &event
}
