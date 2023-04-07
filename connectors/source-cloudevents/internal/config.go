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

var _ cdkgo.SourceConfigAccessor = &cloudEventsConfig{}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &cloudEventsConfig{}
}

type cloudEventsConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	Port               int               `json:"port" yaml:"port"`
	Path               string            `json:"path" yaml:"path"`
	Headers            map[string]string `json:"headers" yaml:"headers"`
	Auth               Auth              `json:"auth" yaml:"auth"`
}

func (c *cloudEventsConfig) GetSecret() cdkgo.SecretAccessor {
	return &c.Auth
}

type Auth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (a *Auth) IsEmpty() bool {
	return a.Username == "" || a.Password == ""
}
