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

var _ cdkgo.SourceConfigAccessor = &googleSheetsConfig{}

func NewGoogleSheetsConfig() cdkgo.SourceConfigAccessor {
	return &googleSheetsConfig{}
}

type googleSheetsConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	// TODO
	Credentials string `json:"credentials" yaml:"credentials" validate:"required"`
	SheetID     string `json:"sheet_id" yaml:"sheet_id" validate:"required"`
	SheetName   string `json:"sheet_name" yaml:"sheet_name" validate:"required"`
}

func (c *googleSheetsConfig) Validate() error {
	// TODO
	return c.SourceConfig.Validate()
}
