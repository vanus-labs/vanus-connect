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

type Config struct {
	Fenodes   string `json:"fenodes" yaml:"fenodes"`
	Username  string `json:"username" yaml:"username"`
	Password  string `json:"password" yaml:"password"`
	DbName    string `json:"db_name" yaml:"db_name"`
	TableName string `json:"table_name" yaml:"table_name"`

	StreamLoad map[string]string `json:"stream_load" yaml:"stream_load"`

	Timeout int `json:"timeout" yaml:"timeout"`
}

func (cfg *Config) Validate() error {
	if cfg == nil {
		return errors.New("cfg is nil")
	}
	if cfg.Fenodes == "" {
		return errors.New("fenodes is mepty")
	}
	if cfg.DbName == "" {
		return errors.New("db name is mepty")
	}
	if cfg.TableName == "" {
		return errors.New("table name is mepty")
	}
	if cfg.Username == "" {
		return errors.New("username is mepty")
	}
	return nil
}
