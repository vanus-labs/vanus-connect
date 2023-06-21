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

var _ cdkgo.SourceConfigAccessor = &googleAnalyticsConfig{}

func GoogleAnalyticsConfig() cdkgo.SourceConfigAccessor {
	return &googleAnalyticsConfig{}
}

type googleAnalyticsConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	// TODO
	Credentials string `json:"credentials" yaml:"credentials"`
	PropertyID  string `json:"property_id" yaml:"property_id" validate:"required"`
	Start_date  string `json:"start_date" yaml:"start_date" validate:"required"`
	End_date    string `json:"end_date" yaml:"end_date" validate:"required"`
}

func (c *googleAnalyticsConfig) GetSecret() cdkgo.SecretAccessor {
	return &c.Credentials
}

func (c *googleAnalyticsConfig) Validate() error {
	// TODO
	return c.SourceConfig.Validate()
}
