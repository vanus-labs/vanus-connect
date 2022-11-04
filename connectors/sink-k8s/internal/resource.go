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
	"encoding/json"
	"fmt"
	"time"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/types"
)

// Reader enables reading from an external store
type Reader interface {
	Read() ([]byte, error)
}

// ResourceReader implements the reader interface for k8s resource
type ResourceReader struct {
	resourceArtifact *unstructured.Unstructured
}

// NewResourceReader creates a new Reader for resource
func NewResourceReader(resource []byte) (Reader, error) {
	if resource == nil {
		return nil, fmt.Errorf("resource does not exist")
	}
	object := make(map[string]interface{})
	err := json.Unmarshal(resource, &object)
	if err != nil {
		return nil, err
	}
	un := &unstructured.Unstructured{Object: object}
	return &ResourceReader{un}, nil
}

func (reader *ResourceReader) Read() ([]byte, error) {
	return yaml.Marshal(reader.resourceArtifact.Object)
}

// Fetch from the location, decode it using explicit types, and unstructure it
func Fetch(reader Reader) (*unstructured.Unstructured, error) {
	var obj []byte
	backoff := wait.Backoff{
		Duration: time.Second,
		Factor:   1,
		Jitter:   1,
		Steps:    5,
		Cap:      time.Duration(0),
	}

	if err := DoWithRetry(&backoff, func() error {
		var e error
		obj, e = reader.Read()
		return e
	}); err != nil {
		return nil, fmt.Errorf("failed to fetch artifact, %w", err)
	}
	return decodeAndUnstructure(obj)
}

func DoWithRetry(backoff *wait.Backoff, f func() error) error {
	var err error
	_ = wait.ExponentialBackoff(*backoff, func() (bool, error) {
		if err = f(); err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("failed after retries: %w", err)
	}
	return nil
}

func decodeAndUnstructure(b []byte) (*unstructured.Unstructured, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: result}, nil
}

func GetGroupVersionKind(obj *unstructured.Unstructured) schema.GroupVersionKind {
	return obj.GroupVersionKind()
}

func GetGroupVersionResource(obj *unstructured.Unstructured) schema.GroupVersionResource {
	gvk := obj.GroupVersionKind()
	pluralExceptions := map[string]string{
		"EventBus": "eventbus",
	}
	resource := namer.NewAllLowercasePluralNamer(pluralExceptions).Name(&types.Type{
		Name: types.Name{
			Name: gvk.Kind,
		},
	})

	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: resource,
	}
}
