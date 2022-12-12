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

	"github.com/linkall-labs/cdk-go/log"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *k8sSink) handleJob(ctx context.Context, data []byte) error {
	var err error
	job := &batchv1.Job{}
	err = json.Unmarshal(data, job)
	if err != nil {
		log.Info("unmarshal to job failed", map[string]interface{}{
			log.KeyError: err,
			"data":       string(data),
		})
		return err
	}

	var operation KubernetesResourceOperation = Create
	if op, ok := job.Annotations["operation"]; ok {
		operation = KubernetesResourceOperation(op)
	}

	switch operation {
	case Create:
		createOpts := metav1.CreateOptions{}
		_, err = s.client.BatchV1().Jobs(job.Namespace).Create(ctx, job, createOpts)
		if err != nil {
			log.Info("create job failed", map[string]interface{}{
				log.KeyError: err,
				"job":        job,
			})
			return err
		}
		log.Info("create job success", map[string]interface{}{
			"job": job.Name,
		})
	case Update:
		updateOpts := metav1.UpdateOptions{}
		_, err = s.client.BatchV1().Jobs(job.Namespace).Update(ctx, job, updateOpts)
		if err != nil {
			log.Info("update job failed", map[string]interface{}{
				log.KeyError: err,
				"job":        job,
			})
			return err
		}
		log.Info("update job success", map[string]interface{}{
			"job": job.Name,
		})
	case Delete:
		background := metav1.DeletePropagationBackground
		deleteOpts := metav1.DeleteOptions{
			PropagationPolicy: &background,
		}
		err = s.client.BatchV1().Jobs(job.Namespace).Delete(ctx, job.Name, deleteOpts)
		if err != nil {
			log.Info("delete job failed", map[string]interface{}{
				log.KeyError: err,
				"job":        job,
			})
			return err
		}
		log.Info("delete job success", map[string]interface{}{
			"job": job.Name,
		})
	default:
		log.Warning("unknown operation type", map[string]interface{}{
			"operation": operation,
		})
	}
	return nil
}
