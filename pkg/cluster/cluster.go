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

package cluster

import (
	"context"
	"fmt"
	"github.com/sumengzs/multi-cluster/pkg/mapper"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ Interface = &cluster{}

type InitOptions func(Interface) error

type cluster struct {
	name       string
	ctx        context.Context
	cancelFunc context.CancelFunc
	status     Code
	scheme     *runtime.Scheme
	client     client.Client
	cache      cache.Cache
	mapper     meta.RESTMapper
	dynamic    dynamic.Interface
	discovery  discovery.DiscoveryInterface
	extensions clientset.Interface
	config     *rest.Config
}

// New returns a new cluster or error
// default status code is Stopped
func New(config *rest.Config, scheme *runtime.Scheme, options ...InitOptions) (Interface, error) {
	var clu *cluster
	var err error

	clu.mapper, err = mapper.Provider(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create mapper: %s", err)
	}

	if clu.client, err = client.New(config, client.Options{Scheme: scheme, Mapper: clu.mapper}); err != nil {
		return nil, fmt.Errorf("failed to create runtime client: %s", err)
	}

	if clu.cache, err = cache.New(config, cache.Options{Scheme: scheme, Mapper: clu.mapper}); err != nil {
		return nil, fmt.Errorf("failed to create runtime cache: %s", err)
	}

	if clu.dynamic, err = dynamic.NewForConfig(config); err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %s", err)
	}

	if clu.extensions, err = clientset.NewForConfig(config); err != nil {
		return nil, fmt.Errorf("failed to create api-extensions client: %s", err)
	}
	clu.extensions.ApiextensionsV1beta1().CustomResourceDefinitions()

	if clu.discovery, err = discovery.NewDiscoveryClientForConfig(config); err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %s", err)
	}

	clu.status = Stopped
	clu.scheme = scheme

	for _, option := range options {
		if err = option(clu); err != nil {
			return nil, fmt.Errorf("failed to initialize options: %s", err)
		}
	}

	return clu, nil
}

func (c *cluster) Start(ctx context.Context) error {
	switch {
	case c.status == Disabled:
		klog.Infof("%s,needs to be enabled first", c)
	case c.status >= Started:
		klog.Infof("%s,no need to start", c)
	case c.status == Stopped:
		c.ctx, c.cancelFunc = context.WithCancel(ctx)
		go func() {
			_ = c.cache.Start(c.ctx)
		}()
		c.status = Started
		klog.Infof("%s", c)
	default:
		return fmt.Errorf("%s", c)
	}
	return nil
}

func (c *cluster) Stop() {
	if c.status == Disabled {
		klog.Infof("%s,no need stop", c)
	}
	if c.status > Stopped {
		c.cancelFunc()
		c.status = Stopped
	}
	klog.Infof("%s", c)
}

func (c *cluster) Name() string {
	return c.name
}

func (c *cluster) Status() Code {
	return c.status
}

func (c *cluster) Disable() {
	if c.status >= Started {
		c.Stop()
	}
	c.status = Disabled
}

func (c *cluster) Client() client.Client {
	return c.client
}

func (c *cluster) Cache() cache.Cache {
	return c.cache
}

func (c *cluster) ApiExtensions() clientset.Interface {
	return c.extensions
}

func (c *cluster) Dynamic() dynamic.Interface {
	return c.dynamic
}

func (c *cluster) RESTMapper() meta.RESTMapper {
	return c.mapper
}

func (c *cluster) Config() *rest.Config {
	return c.config
}

func (c *cluster) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *cluster) String() string {
	return fmt.Sprintf("cluster %s is %s", c.name, c.status)
}
