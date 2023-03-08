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
	"fmt"
	"math/rand"
	"net/http"
	"time"

	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	v20180416 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
)

const (
	name = "Tencent SCF Sink"
)

var (
	functionNamePrefix = "vanus-cos-source-function"
)

var _ cdkgo.SinkConfigAccessor = &scfConfig{}

type scfConfig struct {
	cdkgo.SinkConfig `json:",inline" yaml:",inline"`
	F                Function `json:"function" yaml:"function"`
	Debug            bool     `json:"debug" yaml:"debug"`
	Secret           *Secret  `json:"secret" yaml:"secret"`
}

func (c *scfConfig) GetSecret() cdkgo.SecretAccessor {
	return c.Secret
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &scfConfig{
		Secret: &Secret{},
	}
}

type Function struct {
	Name      string `yaml:"name" json:"name" validate:"required"`
	Region    string `yaml:"region" json:"region"`
	Namespace string `yaml:"namespace" json:"namespace" default:"default"`
}

type Secret struct {
	SecretID  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}

func NewFunctionSink() cdkgo.Sink {
	return &functionSink{}
}

var _ cdkgo.Sink = &functionSink{}

type functionSink struct {
	scfClient *v20180416.Client
	cfg       *scfConfig
	funcName  string
}

func (c *functionSink) Arrived(_ context.Context, events ...*v2.Event) cdkgo.Result {
	req := v20180416.NewInvokeRequest()
	req.FunctionName = &c.cfg.F.Name
	req.Namespace = &c.cfg.F.Namespace

	for idx := range events {
		e := events[idx]
		payload := string(e.Data())
		req.ClientContext = &payload

		res, err := c.scfClient.Invoke(req)
		if err != nil {
			log.Debug("failed to invoke function", map[string]interface{}{
				log.KeyError: err,
			})
			return cdkgo.NewResult(http.StatusInternalServerError, err.Error())
		}
		log.Debug("invoke function success", map[string]interface{}{
			"response": res,
		})
	}

	return cdkgo.SuccessResult
}

func (c *functionSink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	_cfg, ok := cfg.(*scfConfig)
	if !ok {
		return nil
	}

	if _cfg.F.Namespace == "" {
		_cfg.F.Namespace = "default"
	}

	c.cfg = _cfg

	if c.cfg.Debug {
		// log.SetLevel("debug")
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

func (c *functionSink) Destroy() error {
	// nothing to do
	return nil
}
