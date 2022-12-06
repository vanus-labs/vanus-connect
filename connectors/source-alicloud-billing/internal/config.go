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
	cdkgo "github.com/linkall-labs/cdk-go"
)

type billingConfig struct {
	cdkgo.SourceConfig
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Endpoint        string `json:"endpoint" yaml:"endpoint"`
	PullHour        int    `json:"pull_hour" yaml:"pull_hour"`
	Secret          Secret
}

func Config() cdkgo.SourceConfigAccessor {
	return &billingConfig{}
}

func (cfg *billingConfig) GetSecret() cdkgo.SecretAccessor {
	return &cfg.Secret
}

type Secret struct {
	AccessKeyID     string `json:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key"`
}
