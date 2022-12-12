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
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	v2 "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
)

var _ cdkgo.SinkConfigAccessor = &displayConfig{}

type displayConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &displayConfig{}
}

func NewDisplaySink() cdkgo.Sink {
	return &displaySink{}
}

var _ cdkgo.Sink = &displaySink{}

type displaySink struct {
	cfg   *displayConfig
	count int64
}

func (f *displaySink) Arrived(_ context.Context, events ...*v2.Event) cdkgo.Result {
	for _, event := range events {
		atomic.AddInt64(&f.count, 1)
		log.Info(fmt.Sprintf("receive a new event, in total: %d", f.count), map[string]interface{}{})
		v, _ := json.MarshalIndent(event, "", "  ")
		log.Info(string(v), nil)
	}
	return cdkgo.SuccessResult
}

func (f *displaySink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	return nil
}

func (f *displaySink) Name() string {
	return "Display Sink"
}

func (f *displaySink) Destroy() error {
	return nil
}
