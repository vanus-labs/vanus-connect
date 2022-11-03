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

import (
	"github.com/pkg/errors"
)

type InsertMode string

const (
	Insert InsertMode = "insert"
	Upsert InsertMode = "upsert"
)

type Config struct {
	Address  string `json:"address" yaml:"address"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	IndexName string `json:"index_name" yaml:"index_name"`

	Timeout    int        `json:"timeout" yaml:"timeout"`
	PrimaryKey string     `json:"primary_key" yaml:"primary_key"`
	InsertMode InsertMode `json:"insert_mode" yaml:"insert_mode"`
}

func (cfg *Config) Validate() error {
	if cfg == nil {
		return errors.New("cfg is nil")
	}
	if cfg.IndexName == "" {
		return errors.New("config index name is empty")
	}
	if cfg.InsertMode == Upsert && cfg.PrimaryKey == "" {
		return errors.New("insert mode is upsert but primary key is empty")
	}
	return nil
}
