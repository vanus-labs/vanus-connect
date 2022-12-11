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
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"sync/atomic"

	v2 "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"k8s.io/utils/strings/slices"
)

const (
	botService = "bot"

	name                      = "Feishu Sink"
	vanceServiceNameAttribute = "xvfeishuservice"
)

var (
	errFeishuSinkWrongEventNumber = cdkgo.NewResult(http.StatusBadRequest,
		"feishu: the event number must be 1")
	errFeishuSinkEventMissingServiceName = cdkgo.NewResult(http.StatusBadRequest,
		"feishu: missing or invalid service name, please check xvfeishuservice in attributes")
	errFeishuSinkUnsupportedService = cdkgo.NewResult(http.StatusBadRequest, "feishu: unsupported service")
)

var _ cdkgo.SinkConfigAccessor = &feishuConfig{}

type Secret struct {
}

type feishuConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	Enable           []string  `json:"enable" yaml:"enable"`
	Bot              BotConfig `json:"bot" yaml:"bot"`
	Secret           Secret    `json:"secret" yaml:"secret"`
}

func (fc *feishuConfig) Validate() error {
	if err := fc.SinkConfig.Validate(); err != nil {
		return err
	}

	for _, s := range fc.Enable {
		switch s {
		case botService:
			if err := fc.Bot.Validate(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported service %s in enable parameter", s)
		}
	}
	return nil
}

func (fc *feishuConfig) GetSecret() cdkgo.SecretAccessor {
	return &fc.Secret
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &feishuConfig{}
}

func NewFeishuSink() cdkgo.Sink {
	return &feishuSink{
		b: &bot{
			httpClient: resty.New(),
		},
	}
}

var _ cdkgo.Sink = &feishuSink{}

type feishuSink struct {
	cfg   *feishuConfig
	count int64
	b     *bot
}

func (f *feishuSink) Arrived(_ context.Context, events ...*v2.Event) cdkgo.Result {
	// optimize(wenfeng) give an argument to control if this sink support batch?
	if len(events) != 1 {
		return errFeishuSinkWrongEventNumber
	}

	atomic.AddInt64(&f.count, int64(len(events)))

	e := events[0]
	val, exist := e.Extensions()[vanceServiceNameAttribute]
	if !exist {
		return errFeishuSinkEventMissingServiceName
	}

	service, ok := val.(string)
	if !ok {
		return errFeishuSinkEventMissingServiceName
	}
	if !slices.Contains(f.cfg.Enable, service) {
		return errFeishuSinkUnsupportedService
	}

	switch service {
	case botService:
		if err := f.b.sendMessage(e); err != nil {
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		}
	default:
		return errFeishuSinkUnsupportedService
	}

	return cdkgo.SuccessResult
}

func (f *feishuSink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*feishuConfig)
	if !ok {
		return nil
	}
	if err := _cfg.Validate(); err != nil {
		return err
	}
	f.cfg = _cfg
	return f.b.init(_cfg.Bot)
}

func (f *feishuSink) Name() string {
	return name
}

func (f *feishuSink) Destroy() error {
	// nothing to do
	return nil
}
