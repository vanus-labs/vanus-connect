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

	ce "github.com/cloudevents/sdk-go/v2"
	cdkgo "github.com/linkall-labs/cdk-go"
	"github.com/linkall-labs/cdk-go/log"
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
	cdkgo.SinkConfig
	Secret *Secret `json:"secret" yaml:"secret"`
}

func NewConfig() cdkgo.SinkConfigAccessor {
	return &Config{
		Secret: &Secret{},
	}
}

func (c *Config) GetSecret() cdkgo.SecretAccessor {
	return c.Secret
}

type Secret struct {
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	AuthSource string `json:"authSource" yaml:"authSource"`
}

type k8sSink struct {
	cfg    *Config
	client *kubernetes.Clientset
}

func NewKubernetesSink() cdkgo.Sink {
	return &k8sSink{}
}

func (s *k8sSink) Initialize(_ context.Context, cfg cdkgo.ConfigAccessor) error {
	s.cfg = cfg.(*Config)
	kubeconfig := GetKubeConfigFromEnv()
	config, err := GetInClusterOrKubeConfig(kubeconfig)
	if err != nil {
		return err
	}

	s.client = kubernetes.NewForConfigOrDie(config)
	return nil
}

func (s *k8sSink) Arrived(ctx context.Context, events ...*ce.Event) cdkgo.Result {
	if len(events) != 1 {
		return cdkgo.NewResult(http.StatusInternalServerError, "the event number must be 1")
	}
	event := events[0]
	var err error
	log.Info("receive an event", map[string]interface{}{
		"event_id": event.ID(),
	})
	data := event.Data()
	reader, err := NewResourceReader(data)
	if err != nil {
		log.Info("new resource reader failed", map[string]interface{}{
			log.KeyError: err,
		})
		return cdkgo.NewResult(http.StatusBadRequest, "resource read failed")
	}

	uObj, err := Fetch(reader)
	if err != nil {
		log.Info("FetchArtifact failed", map[string]interface{}{
			log.KeyError: err,
		})
		return cdkgo.NewResult(http.StatusBadRequest, "fetch artifact failed")
	}

	gvk := GetGroupVersionKind(uObj)
	switch KubernetesResourceKind(gvk.Kind) {
	case Pod:
		err = s.handlePod(ctx, data)
		if err != nil {
			return cdkgo.NewResult(http.StatusInternalServerError, "handle resource error")
		}
	case Job:
		err = s.handleJob(ctx, data)
		if err != nil {
			return cdkgo.NewResult(http.StatusInternalServerError, "handler resource error")
		}
	default:
		return cdkgo.NewResult(http.StatusNotAcceptable, "can't handle resource")
	}
	return cdkgo.SuccessResult
}

func (s *k8sSink) Name() string {
	return "kubernetes-sink"
}

func (s *k8sSink) Destroy() error {
	return nil
}
