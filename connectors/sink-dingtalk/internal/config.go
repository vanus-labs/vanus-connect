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
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/connector/sink/dingtalk/internal/bot"
)

var _ cdkgo.SinkConfigAccessor = &dingtalkConfig{}

type dingtalkConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	Bot              bot.Config `json:"bot" yaml:"bot"`
}

func (fc *dingtalkConfig) Validate() error {
	if err := fc.Bot.Validate(); err != nil {
		return err
	}
	return fc.SinkConfig.Validate()
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &dingtalkConfig{}
}
