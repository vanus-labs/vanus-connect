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

import cdkgo "github.com/vanus-labs/cdk-go"

var _ cdkgo.SourceConfigAccessor = &shopifyConfig{}

type shopifyConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	ShopName           string `json:"shop_name" yaml:"shop_name" validate:"required"`
	ApiAccessToken     string `json:"api_access_token" yaml:"api_access_token" validate:"required"`
	SyncBeginDate      string `json:"sync_begin_date" yaml:"sync_begin_date" validate:"required"`
	DelaySecond        int    `json:"delay_second" yaml:"delay_second"`

	EventTypes []ApiType `json:"event_type" yaml:"event_type"`

	Port         int    `json:"port" yaml:"port"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &shopifyConfig{}
}
