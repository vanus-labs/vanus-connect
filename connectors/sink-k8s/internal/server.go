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
	"net/http"

	"github.com/cloudevents/sdk-go/v2/protocol"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/linkall-labs/cdk-go/connector"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/linkall-labs/cdk-go/log"
	cdkutil "github.com/linkall-labs/cdk-go/utils"

	"k8s.io/client-go/kubernetes"
)

// KubernetesResourceKind refers to the type of kind performed on the K8s resource
type KubernetesResourceKind string

// possible values for KubernetesResourceKind
const (
	Pod KubernetesResourceKind = "Pod"
	Job KubernetesResourceKind = "Job"
)

// KubernetesResourceOperation refers to the type of operation performed on the K8s resource
type KubernetesResourceOperation string

// possible values for KubernetesResourceOperation
const (
	Create KubernetesResourceOperation = "create" // create the resource
	Update KubernetesResourceOperation = "update" // updates the resource
	Delete KubernetesResourceOperation = "delete" // deletes the resource
)

type Config struct {
	Secret *Secret `json:"-" yaml:"-"`
}

type Secret struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	AuthSource string `json:"authSource" yaml:"authSource"`
}

type k8sSink struct {
	cfg    *Config
	client *kubernetes.Clientset
	logger log.Logger
}

func NewKubernetesSink() connector.Sink {
	return &k8sSink{}
}

func (s *k8sSink) Init(cfgPath, secretPath string) error {
	cfg := &Config{}
	if err := cdkutil.ParseConfig(cfgPath, cfg); err != nil {
		return err
	}

	if connector.IsSecretEnable() {
		secret := &Secret{}
		if err := cdkutil.ParseConfig(secretPath, secret); err != nil {
			return err
		}
		cfg.Secret = secret
	}

	kubeconfig := GetKubeConfigFromEnv()
	config, err := GetInClusterOrKubeConfig(kubeconfig)
	if err != nil {
		panic(err)
	}

	s.client = kubernetes.NewForConfigOrDie(config)
	s.cfg = cfg
	return nil
}

func (s *k8sSink) Name() string {
	return "kubernetes-sink"
}

func (s *k8sSink) SetLogger(logger log.Logger) {
	s.logger = logger
}

func (s *k8sSink) Destroy() error {
	return nil
}

func (s *k8sSink) Receive(ctx context.Context, event ce.Event) protocol.Result {
	var err error
	log.Info(ctx, "receive an event", map[string]interface{}{
		"event_id": event.ID(),
	})
	data := event.Data()
	reader, err := NewResourceReader(data)
	if err != nil {
		log.Info(ctx, "new resource reader failed", map[string]interface{}{
			log.KeyError: err,
		})
		return cehttp.NewResult(http.StatusBadRequest, "")
	}

	uObj, err := Fetch(reader)
	if err != nil {
		log.Info(ctx, "FetchArtifact failed", map[string]interface{}{
			log.KeyError: err,
		})
		return cehttp.NewResult(http.StatusBadRequest, "")
	}

	gvk := GetGroupVersionKind(uObj)
	switch KubernetesResourceKind(gvk.Kind) {
	case Pod:
		err = s.handlePod(ctx, data)
		if err != nil {
			return cehttp.NewResult(http.StatusInternalServerError, "")
		}
	case Job:
		err = s.handleJob(ctx, data)
		if err != nil {
			return cehttp.NewResult(http.StatusInternalServerError, "")
		}
	default:
		return cehttp.NewResult(http.StatusNotAcceptable, "")
	}

	return cehttp.NewResult(http.StatusOK, "")
}
