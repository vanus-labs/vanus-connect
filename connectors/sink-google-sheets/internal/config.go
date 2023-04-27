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
	"fmt"

	cdkgo "github.com/vanus-labs/cdk-go"
)

var _ cdkgo.SinkConfigAccessor = &GoogleSheetConfig{}

func NewGoogleSheetConfig() cdkgo.SinkConfigAccessor {
	return &GoogleSheetConfig{}
}

type GoogleSheetConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	// Google  Credentials JSON
	Credentials string `json:"credentials" yaml:"credentials"`
	OAuth       *OAuth `json:"oauth" yaml:"oauth"`

	SheetID       string `json:"sheet_id" yaml:"sheet_id" validate:"required"`
	SheetName     string `json:"sheet_name" yaml:"sheet_name" validate:"required"`
	FlushInterval int    `json:"flush_interval" yaml:"flush_interval"`
	FlushSize     int    `json:"flush_size" yaml:"flush_size"`

	Summary []SummaryConfig `json:"summary" yaml:"summary"`
}

func (cfg *GoogleSheetConfig) Validate() error {
	if cfg.OAuth == nil && cfg.Credentials == "" {
		return errors.New("credentials or oauth must set one")
	}
	if len(cfg.Summary) > 0 {
		sheetNameMap := map[string]struct{}{cfg.SheetName: {}}
		for _, summary := range cfg.Summary {
			if _, exist := sheetNameMap[summary.SheetName]; exist {
				return errors.New(fmt.Sprintf("sheetName %s is repeated", summary.SheetName))
			} else {
				sheetNameMap[summary.SheetName] = struct{}{}
			}
		}
	}
	return cfg.SinkConfig.Validate()
}
