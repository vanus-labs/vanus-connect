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
	"context"
	"strconv"

	"github.com/linkall-labs/cdk-go/config"
	"github.com/linkall-labs/cdk-go/log"
)

type Config struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Endpoint        string `json:"endpoint"`
	PullHour        int    `json:"pull_hour"`
}

func getConfig(ctx context.Context) Config {
	c := config.Accessor
	conf := Config{
		AccessKeyID:     c.Get("access_key_id"),
		SecretAccessKey: c.Get("secret_access_key"),
		Endpoint:        c.Get("endpoint"),
	}
	if conf.Endpoint == "" {
		conf.Endpoint = "https://ce.us-east-1.amazonaws.com"
	}
	if c.Get("pull_hour") != "" {
		pullHour, err := strconv.Atoi(c.Get("pull_hour"))
		if err != nil {
			log.FromContext(ctx).Info("pull hour parse to int error", "error", err)
		} else {
			conf.PullHour = pullHour
		}
	}
	if conf.PullHour <= 0 || conf.PullHour >= 24 {
		conf.PullHour = 2
	}
	return conf
}
