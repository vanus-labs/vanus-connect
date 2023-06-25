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

func WhatsAppConfig() cdkgo.SinkConfigAccessor {
	return &WhatsappConfig{}
}

type WhatsappConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	FileName         string `json:"file_name" yaml:"file_name"`
	Data             string `json:"data" yaml:"data"`
}

func (c *WhatsappConfig) Validate() error {
	return c.SinkConfig.Validate()
}
