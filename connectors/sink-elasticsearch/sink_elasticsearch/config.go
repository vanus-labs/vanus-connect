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

package sink_elasticsearch

import (
	"github.com/linkall-labs/cdk-go/config"
	"strings"
)

type Config struct {
	Addresses  []string `json:"addresses"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	SkipVerify bool     `json:"skipVerify"`

	IndexName string `json:"indexName"`
}

func getConfig() Config {
	c := config.Accessor
	conf := Config{
		Addresses:  strings.Split(c.Get("addresses"), ","),
		Username:   c.Get("username"),
		Password:   c.Get("password"),
		IndexName:  c.Get("indexName"),
		SkipVerify: true,
	}
	return conf
}
