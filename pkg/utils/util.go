/*
Copyright 2023 The Multi Cluster Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
)

func GetConfigOrDie(master, kubeConfig string) *rest.Config {
	var config *rest.Config
	if len(kubeConfig) == 0 {
		config = ctrl.GetConfigOrDie()
	} else {
		var err error
		config, err = clientcmd.BuildConfigFromFlags(master, kubeConfig)
		if err != nil {
			klog.ErrorS(err, "unable to get kubeconfig", "master", master, "kubeConfig", kubeConfig)
			os.Exit(1)
		}
	}
	return config
}
