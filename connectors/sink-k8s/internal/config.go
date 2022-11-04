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
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func GetKubeConfigFromEnv() string {
	home := os.Getenv("HOME")
	if home != "" {
		fpath := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(fpath); err != nil {
			return ""
		}
		return fpath
	}
	return ""
}

func GetInClusterOrKubeConfig(kubeconfig string) (config *rest.Config, rerr error) {
	config, rerr = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if rerr != nil {
		klog.Errorf("auth from kubeconfig failed:%v", rerr)
		return nil, rerr
	}
	return config, nil
}
