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

package auth

import "github.com/pkg/errors"

type Type string

const (
	Basic Type = "basic"
	Hmac  Type = "hmac"

	DefaultHeaderSignature = "X-Vanus-Signature"
)

type Config struct {
	Type  Type      `json:"type" yaml:"type"`
	Basic BasicAuth `json:"basic" yaml:"basic"`
	HMAC  HMAC      `json:"hmac" yaml:"hmac"`
}

type HMAC struct {
	Header string `json:"header" yaml:"header"`
	Secret string `json:"secret" yaml:"secret" `
}

type BasicAuth struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (c *Config) Validate() error {
	if c == nil {
		return nil
	}
	switch c.Type {
	case Hmac:
		if c.HMAC.Secret == "" {
			return errors.New("auth hmac secret can not empty")
		}
	case Basic:
		if c.Basic.Username == "" || c.Basic.Password == "" {
			return errors.New("auth basic username and password can not empty")
		}
	default:
		return errors.New("auth type invalid")
	}
	return nil
}
