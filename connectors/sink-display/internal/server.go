// Copyright 2022 Linkall Inc.
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
	"bytes"
	"context"
	"encoding/json"
	"sync/atomic"

	v2 "github.com/cloudevents/sdk-go/v2"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/config"
	"github.com/vanus-labs/cdk-go/connector"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name = "Display Sink"
)

func NewConfig() cdkgo.SinkConfigAccessor {
	return &cdkgo.SinkConfig{}
}

var _ cdkgo.Sink = &displaySink{}

func NewDisplaySink() cdkgo.Sink {
	return &displaySink{}
}

type displaySink struct {
	count int64
}

func (ds *displaySink) Initialize(_ context.Context, _ config.ConfigAccessor) error {
	return nil
}

func (ds *displaySink) Name() string {
	return name
}

func (ds *displaySink) Destroy() error {
	return nil
}

func (ds *displaySink) Arrived(_ context.Context, events ...*v2.Event) connector.Result {
	for idx := range events {
		e := events[idx]
		atomic.AddInt64(&ds.count, 1)
		log.Info().Int64("total", atomic.LoadInt64(&ds.count)).Msg("receive a new event")
		d, err := e.MarshalJSON()
		if err != nil {
			log.Warn().Err(err).Str("event", e.String()).Msg("received a new event, but failed to marshal to JSON")
		} else {
			buf := bytes.NewBuffer([]byte{})
			_ = json.Indent(buf, d, "", "  ")
			println(string(buf.Bytes()))
		}
		println()
	}
	return cdkgo.SuccessResult
}
