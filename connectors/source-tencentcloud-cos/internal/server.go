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
	"fmt"
	cdkutil "github.com/linkall-labs/cdk-go/utils"
	"math/rand"
	"strconv"
	"sync"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
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
	runtime            = "Go1"
	handler            = "main"
	funcDesc           = "auto-created function by Vanus for syncing COS event"
	funcMemSize        = int64(64)
	functionNamePrefix = "vanus-cos-source-function"
	defaultFunction    = Function{
		Bucket: "vanus-1253760853",
		Region: "ap-beijing",
		Path:   "/vanus/cos-source/dev/main.zip",
	}

	triggerType   = "cos"
	triggerEnable = "OPEN"
	triggerDesc   = `{"event":"cos:ObjectCreated:*"}`
)

type Config struct {
	Target         string   `json:"v_target" yaml:"v_target"`
	BucketEndpoint string   `json:"bucket_endpoint" yaml:"bucket_endpoint"`
	Function       Function `json:"function" yaml:"function"`
	Region         string   `json:"function_region" yaml:"function_region"`
	Namespace      string   `json:"namespace" yaml:"namespace"`
	Debug          bool     `json:"debug" yaml:"debug"`
	Eventbus       string   `json:"eventbus" yaml:"eventbus"`
	Secret         *Secret  `json:"-" yaml:"-"`
}

type Function struct {
	Bucket string `yaml:"bucket" json:"bucket"`
	Region string `yaml:"region" json:"region"`
	Path   string `yaml:"path" json:"path"`
}

func (f Function) isValid() bool {
	return f.Bucket != "" && f.Region != "" && f.Path != ""
}

type Secret struct {
	SecretID  string `json:"secret_id" yaml:"secret_id"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
}

func NewCosSink() connector.Source {
	return &cosSource{}
}

type cosSource struct {
	scfClient *v20180416.Client
	logger    log.Logger
	cfg       *Config
	funcName  string
	mutex     sync.Mutex
}

func (c *cosSource) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}

	secret := &Secret{}
	if err := cdkutil.ParseConfig(secretPath, secret); err != nil {
		return err
	}
	cfg.Secret = secret

	if cfg.Namespace == "" {
		cfg.Namespace = "default"
	}
	if !cfg.Function.isValid() {
		cfg.Function = defaultFunction
	}
	c.cfg = cfg

	cli, err := v20180416.NewClient(&common.Credential{
		SecretId:  c.cfg.Secret.SecretID,
		SecretKey: c.cfg.Secret.SecretKey,
	}, c.cfg.Region, profile.NewClientProfile())

	if err != nil {
		return err
	}

	c.scfClient = cli
	r := rand.New(rand.NewSource(time.Now().Unix()))
	c.funcName = fmt.Sprintf("%s-%d", functionNamePrefix, r.Uint64())

	return nil
}

func (c *cosSource) Name() string {
	return name
}

func (c *cosSource) SetLogger(logger log.Logger) {
	c.logger = logger
}

func (c *cosSource) Run() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// TODO 检查cos配置
	debugStr := strconv.FormatBool(c.cfg.Debug)
	req := v20180416.NewCreateFunctionRequest()
	req.FunctionName = &c.funcName
	req.Description = &funcDesc
	req.MemorySize = &funcMemSize
	req.Runtime = &runtime
	req.Handler = &handler
	req.Namespace = &c.cfg.Namespace
	req.Environment = &v20180416.Environment{
		Variables: []*v20180416.Variable{
			{
				Key:   &EnvEventGateway,
				Value: &c.cfg.Target,
			},
			{
				Key:   &EnvFuncName,
				Value: &c.funcName,
			},
			{
				Key:   &EnvVanusEventbus,
				Value: &c.cfg.Eventbus,
			},
			{
				Key:   &EnvDebugMode,
				Value: &debugStr,
			},
		},
	}

	req.Code = &v20180416.Code{
		CosBucketName:   &c.cfg.Function.Bucket,
		CosBucketRegion: &c.cfg.Function.Region,
		CosObjectName:   &c.cfg.Function.Path,
	}

	res, err := c.scfClient.CreateFunction(req)
	if err != nil {
		return err
	}

	log.Info("success to create function", map[string]interface{}{
		"response":      res.ToJsonString(),
		"function_name": c.funcName,
	})

	for {
		getReq := v20180416.NewGetFunctionRequest()
		getReq.FunctionName = &c.funcName
		getReq.Namespace = &c.cfg.Namespace
		getRes, err := c.scfClient.GetFunction(getReq)
		if err != nil {
			return err
		}
		if *getRes.Response.Status == "Active" {
			break
		}
		log.Info("function isn't ready", map[string]interface{}{
			"function_name": c.funcName,
			"status":        *getRes.Response.Status,
		})
		time.Sleep(time.Second)
	}

	log.Info("function is ready to create trigger", map[string]interface{}{
		"function_name": c.funcName,
	})

	createTriggerReq := v20180416.NewCreateTriggerRequest()
	createTriggerReq.FunctionName = &c.funcName
	createTriggerReq.Namespace = &c.cfg.Namespace
	createTriggerReq.TriggerName = &c.cfg.BucketEndpoint
	createTriggerReq.Type = &triggerType
	createTriggerReq.TriggerDesc = &triggerDesc
	createTriggerReq.Enable = &triggerEnable

	createTriggerRes, err := c.scfClient.CreateTrigger(createTriggerReq)
	if err != nil {
		return err
	}

	log.Info("success to create trigger", map[string]interface{}{
		"response":      createTriggerRes.ToJsonString(),
		"function_name": c.funcName,
	})
	return nil
}

func (c *cosSource) Destroy() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	req := v20180416.NewDeleteFunctionRequest()
	req.FunctionName = &c.funcName
	res, err := c.scfClient.DeleteFunction(req)
	if err != nil {
		return err
	}

	log.Info("success to delete function", map[string]interface{}{
		"response":      res.ToJsonString(),
		"function_name": c.funcName,
	})
	return nil
}

func (c *cosSource) Adapt(_ ...interface{}) ce.Event {
	panic(fmt.Sprintf("%s doesn't support adaptor", name))
}
