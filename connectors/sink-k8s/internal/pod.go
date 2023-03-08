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
	"encoding/json"

	"github.com/vanus-labs/cdk-go/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *k8sSink) handlePod(ctx context.Context, data []byte) error {
	var err error
	pod := &corev1.Pod{}
	err = json.Unmarshal(data, pod)
	if err != nil {
		log.Info("unmarshal to pod failed", map[string]interface{}{
			log.KeyError: err,
			"data":       string(data),
		})
		return err
	}

	var operation KubernetesResourceOperation = Create
	if op, ok := pod.Annotations["operation"]; ok {
		operation = KubernetesResourceOperation(op)
	}

	switch operation {
	case Create:
		createOpts := metav1.CreateOptions{}
		_, err = s.client.CoreV1().Pods(pod.Namespace).Create(ctx, pod, createOpts)
		if err != nil {
			log.Info("create pod failed", map[string]interface{}{
				log.KeyError: err,
				"pod":        pod,
			})
			return err
		}
		log.Info("create pod success", map[string]interface{}{
			"pod": pod.Name,
		})
	case Update:
		updateOpts := metav1.UpdateOptions{}
		_, err = s.client.CoreV1().Pods(pod.Namespace).Update(ctx, pod, updateOpts)
		if err != nil {
			log.Info("update pod failed", map[string]interface{}{
				log.KeyError: err,
				"pod":        pod,
			})
			return err
		}
		log.Info("update pod success", map[string]interface{}{
			"pod": pod.Name,
		})
	case Delete:
		deleteOpts := metav1.DeleteOptions{}
		err = s.client.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, deleteOpts)
		if err != nil {
			log.Info("delete pod failed", map[string]interface{}{
				log.KeyError: err,
				"pod":        pod,
			})
			return err
		}
		log.Info("delete pod success", map[string]interface{}{
			"pod": pod.Name,
		})
	default:
		log.Warning("unknown operation type", map[string]interface{}{
			"operation": operation,
		})
	}
	return nil
}
