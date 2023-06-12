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

package pool

import (
	"os"
	"os/user"
	"path"
	"time"

	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
)

type Options struct {
	KubeConfig   string        `json:"kubeConfig,omitempty" yaml:"kubeConfig,omitempty"`
	Master       string        `json:"master,omitempty" yaml:"master,omitempty"`
	QPS          float32       `json:"qps,omitempty" yaml:"qps,omitempty"`
	Burst        int           `json:"burst,omitempty" yaml:"burst,omitempty"`
	InformerSync time.Duration `json:"informerSync,omitempty" yaml:"informerSync,omitempty"`
}

func NewKubernetesOptions() (option *Options) {
	option = &Options{
		QPS:   1e6,
		Burst: 1e6,
	}
	homePath := homedir.HomeDir()
	if homePath == "" {
		if u, err := user.Current(); err == nil {
			homePath = u.HomeDir
		}
	}
	userHomeConfig := path.Join(homePath, ".kube/config")
	if _, err := os.Stat(userHomeConfig); err == nil {
		option.KubeConfig = userHomeConfig
	}
	return
}

func (k *Options) Validate() []error {
	errs := make([]error, 0)
	if k.KubeConfig != "" {
		if _, err := os.Stat(k.KubeConfig); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (k *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	fs.StringVar(&k.KubeConfig, "kubeConfig", c.KubeConfig,
		"Path for kubernetes kubeConfig file, if left blank, will use in cluster way.")
	fs.StringVar(&k.Master, "master", c.Master,
		"Used to generate kubeConfig for downloading, if not specified, will use host in kubeConfig.")
	fs.DurationVar(&k.InformerSync, "informerSync", c.InformerSync,
		"Used to sync kubernetes resource")
}
