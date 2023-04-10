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

package bot

import (
	"errors"
	"net/url"
)

type WebHook struct {
	ChatGroup string `json:"chat_group" yaml:"chat_group" validate:"required"`
	URL       string `json:"url" yaml:"url" validate:"required"`
	Signature string `json:"signature" yaml:"signature" validate:"required"`
}

type Config struct {
	Webhooks []WebHook `json:"webhooks" yaml:"webhooks" validate:"dive"`
	Default  string    `json:"default" yaml:"default"`
}

func (c *Config) Validate() error {
	if len(c.Webhooks) == 0 {
		return errors.New("the Bot.webhooks can't be empty when dynamic_route is false")
	}
	if c.Default != "" && !c.defaultExist() {
		return errors.New("the Bot.default not exist in webhooks.chatGroup")
	}
	for _, wh := range c.Webhooks {
		_, err := url.Parse(wh.URL)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) defaultExist() bool {
	for _, webhook := range c.Webhooks {
		if webhook.ChatGroup == c.Default {
			return true
		}
	}
	return false
}
