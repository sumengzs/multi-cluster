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
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Interface interface {
	Runnable
	Status
	Name() string
	Client() client.Client
	Cache() cache.Cache
	ApiExtensions() clientset.Interface
	Dynamic() dynamic.Interface
	RESTMapper() meta.RESTMapper
	Config() *rest.Config
	Discovery() discovery.DiscoveryInterface
}

type Status interface {
	Status() Code
	Disable()
}

type Runnable interface {
	Start(ctx context.Context) error
	Stop()
}
