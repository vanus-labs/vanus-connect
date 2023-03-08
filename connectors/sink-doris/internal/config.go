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
	cdkgo "github.com/vanus-labs/cdk-go"
)

type dorisConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`

	StreamLoad map[string]string `json:"stream_load" yaml:"stream_load"`

	Timeout      int    `json:"timeout" yaml:"timeout"`
	LoadInterval int    `json:"load_interval" yaml:"load_interval"`
	LoadSize     int    `json:"load_size" yaml:"load_size"`
	Secret       Secret `json:"secret" yaml:"secret"`
}

func Config() cdkgo.SinkConfigAccessor {
	return &dorisConfig{}
}

func (cfg *dorisConfig) GetSecret() cdkgo.SecretAccessor {
	return &cfg.Secret
}

type Secret struct {
	Fenodes   string `json:"fenodes" yaml:"fenodes" validate:"required"`
	DbName    string `json:"db_name" yaml:"db_name" validate:"required"`
	TableName string `json:"table_name" yaml:"table_name" validate:"required"`
	Username  string `json:"username" yaml:"username" validate:"required"`
	Password  string `json:"password" yaml:"password" validate:"required"`
}
