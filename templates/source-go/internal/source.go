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
	"math/rand"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.Source = &exampleSource{}

func NewExampleSource() cdkgo.Source {
	return &exampleSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type exampleSource struct {
	config *exampleConfig
	events chan *cdkgo.Tuple
	number int
}

func (s *exampleSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	// TODO
	s.config = cfg.(*exampleConfig)
	go s.loopProduceEvent()
	return nil
}

func (s *exampleSource) Name() string {
	// TODO
	return "ExampleSource"
}

func (s *exampleSource) Destroy() error {
	// TODO
	return nil
}

func (s *exampleSource) Chan() <-chan *cdkgo.Tuple {
	// TODO
	return s.events
}

func (s *exampleSource) loopProduceEvent() *ce.Event {
	for {
		event := s.makeEvent()
		b, _ := json.Marshal(event)
		success := func() {
			fmt.Println("send event success: " + string(b))
		}
		failed := func(err error) {
			fmt.Println("send event failed: " + string(b) + ", error: " + err.Error())
		}
		s.events <- cdkgo.NewTuple(event, success, failed)
	}
}

func (s *exampleSource) makeEvent() *ce.Event {
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)+100))
	s.number++
	event := ce.NewEvent()
	event.SetID(fmt.Sprintf("id-%d", s.number))
	event.SetSource("testSource")
	event.SetType("testType")
	event.SetExtension("t", time.Now())
	event.SetData(ce.ApplicationJSON, map[string]interface{}{
		"number": s.number,
		"string": fmt.Sprintf("str-%d", s.number),
	})
	return &event
}
