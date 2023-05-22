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

import cdkgo "github.com/vanus-labs/cdk-go"

var _ cdkgo.SinkConfigAccessor = &WhatsappConfig{}

func NewExampleConfig() cdkgo.SinkConfigAccessor {
	return &WhatsappConfig{}
}

type WhatsappConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	// TODO
	Secret Secret `json:"secret" yaml:"secret"`
}

func (c *WhatsappConfig) GetSecret() cdkgo.SecretAccessor {
	return &c.Secret
}

func (c *WhatsappConfig) Validate() error {
	// TODO
	return c.SinkConfig.Validate()
}

type Secret struct {
	Host     string `json:"host" yaml:"host" validate:"required"`
	Username string `json:"username" yaml:"username" validate:"required"`
	Password string `json:"password" yaml:"password" validate:"required"`
}
