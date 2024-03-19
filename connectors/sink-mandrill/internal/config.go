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
	cdkgo "github.com/vanus-labs/cdk-go"
)

func NewConfig() cdkgo.SinkConfigAccessor {
	return &sinkConfig{}
}

type sinkConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`

	ApiKey       string `json:"api_key" yaml:"api_key" validate:"required"`
	TemplateName string `json:"template_name" yaml:"template_name"`
}
