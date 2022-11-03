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
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cdkutil "github.com/linkall-labs/cdk-go/utils"
	"math/rand"
	"net/http"
	"time"

	"github.com/cloudevents/sdk-go/v2"
	"github.com/linkall-labs/cdk-go/connector"
	"github.com/linkall-labs/cdk-go/log"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	v20180416 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
)

const (
	name = "Tencent Cloud COS Source"
)

var (
	functionNamePrefix = "vanus-cos-source-function"
)

type Config struct {
	F      Function `json:"function" yaml:"function"`
	Debug  bool     `json:"debug" yaml:"debug"`
	Secret *Secret  `json:"-" yaml:"-"`
}

type Function struct {
	Name      string `yaml:"name" json:"name"`
	Region    string `yaml:"region" json:"region"`
	Namespace string `yaml:"namespace" json:"namespace" default:"default"`
}

func (f Function) isValid() bool {
	return f.Name != ""
}

type Secret struct {
	SecretID  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}

func NewFunctionSink() connector.Sink {
	return &functionSink{}
}

type functionSink struct {
	scfClient *v20180416.Client
	logger    log.Logger
	cfg       *Config
	funcName  string
}

func (c *functionSink) Receive(_ context.Context, event v2.Event) protocol.Result {
	req := v20180416.NewInvokeRequest()
	req.FunctionName = &c.cfg.F.Name
	req.Namespace = &c.cfg.F.Namespace
	payload := string(event.Data())
	req.ClientContext = &payload

	res, err := c.scfClient.Invoke(req)
	if err != nil {
		c.logger.Debug("failed to invoke function", map[string]interface{}{
			log.KeyError: err,
		})
		return v2.NewHTTPResult(http.StatusInternalServerError, err.Error())
	}
	c.logger.Debug("invoke function success", map[string]interface{}{
		"response": res,
	})
	return protocol.ResultACK
}

func (c *functionSink) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}

	if cfg.F.Namespace == "" {
		cfg.F.Namespace = "default"
	}

	secret := &Secret{}
	if err := cdkutil.ParseConfig(secretPath, secret); err != nil {
		return err
	}
	cfg.Secret = secret

	if !cfg.F.isValid() {
		return errors.New("invalid function configuration")
	}
	c.cfg = cfg

	if c.cfg.Debug {
		c.logger.SetLevel("debug")
	}

	cli, err := v20180416.NewClient(&common.Credential{
		SecretId:  c.cfg.Secret.SecretID,
		SecretKey: c.cfg.Secret.SecretKey,
	}, c.cfg.F.Region, profile.NewClientProfile())

	if err != nil {
		return err
	}

	c.scfClient = cli
	r := rand.New(rand.NewSource(time.Now().Unix()))
	c.funcName = fmt.Sprintf("%s-%d", functionNamePrefix, r.Uint64())
	return nil
}

func (c *functionSink) Name() string {
	return name
}

func (c *functionSink) SetLogger(logger log.Logger) {
	c.logger = logger
}

func (c *functionSink) Destroy() error {
	// nothing to do
	return nil
}
