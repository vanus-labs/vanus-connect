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
	"net/url"

	"github.com/pkg/errors"
	cdkgo "github.com/vanus-labs/cdk-go"
)

func NewConfig() cdkgo.SinkConfigAccessor {
	return &httpConfig{}
}

type httpConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`

	Target  string            `json:"target" yaml:"target"`
	Method  string            `json:"method" yaml:"method"`
	Headers map[string]string `json:"headers" yaml:"headers"`
	Auth    Auth              `json:"auth" yaml:"auth"`
}

func (c *httpConfig) GetSecret() cdkgo.SecretAccessor {
	return &c.Auth
}

type Auth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (c *httpConfig) Validate() error {
	_, err := url.Parse(c.Target)
	if err != nil {
		return errors.Wrap(err, "target url parse error")
	}
	return c.Config.Validate()
}
