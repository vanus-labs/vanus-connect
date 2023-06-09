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
	"errors"

	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SinkConfigAccessor = &FacebookLeadAdsConfig{}

func NewFacebookLeadAdsConfig() cdkgo.SinkConfigAccessor {
	return &FacebookLeadAdsConfig{}
}

type FacebookLeadAdsConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	// Facebook Lead Ads credentials
	PageId      string `json:"page_id" yaml:"page_id" validate:"required`
	AccessToken string `json:"access_token" yaml:"access_token" validate:"required"`
}

func (c *FacebookLeadAdsConfig) GetSecret() cdkgo.SecretAccessor {
	return &c.AccessToken
}

func (c *FacebookLeadAdsConfig) Validate() error {
	if len(c.AccessToken) == 0 {
		return errors.New("Access Token must be set")
	}
	if len(c.PageId) == 0 {
		return errors.New("Page ID must be set")
	}
	return c.SinkConfig.Validate()
}
