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

package mapper

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sync"
	"time"
)

// GetGroupVersionResource is a helper to map GVK(schema.GroupVersionKind) to GVR(schema.GroupVersionResource).
func GetGroupVersionResource(restMapper meta.RESTMapper, gvk schema.GroupVersionKind) (schema.GroupVersionResource, error) {
	restMapping, err := restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return restMapping.Resource, nil
}

type cacheRESTMapper struct {
	restMapper  meta.RESTMapper
	gvkToGVR    sync.Map
	gvrToGVK    sync.Map
	gvrToGVKErr sync.Map
}

var defaultRecheckTime = time.Second * 300
var NearErr = fmt.Errorf("recent query gvk error")

func (c *cacheRESTMapper) KindFor(resource schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	value, ok := c.gvrToGVK.Load(resource.String())
	if ok {
		return value.(schema.GroupVersionKind), nil
	}
	_time, ok := c.gvrToGVKErr.Load(resource.String())
	if ok {
		record := _time.(*time.Time)
		if time.Since(*record) < defaultRecheckTime {
			return schema.GroupVersionKind{}, NearErr
		}
	}
	gvk, err := c.restMapper.KindFor(resource)
	if err != nil {
		now := time.Now()
		c.gvrToGVKErr.Store(resource.String(), &now)
		return schema.GroupVersionKind{}, err
	}
	c.gvrToGVK.Store(resource.String(), gvk)
	return gvk, nil
}

func (c *cacheRESTMapper) KindsFor(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error) {
	return c.restMapper.KindsFor(resource)
}

func (c *cacheRESTMapper) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	return c.restMapper.ResourceFor(input)
}

func (c *cacheRESTMapper) ResourcesFor(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error) {
	return c.restMapper.ResourcesFor(input)
}

func (c *cacheRESTMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	if len(versions) > 1 {
		return c.restMapper.RESTMapping(gk, versions...)
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("expected at least one version")
	}
	gvk := gk.WithVersion(versions[0])
	value, ok := c.gvkToGVR.Load(gvk)
	if ok {
		return value.(*meta.RESTMapping), nil
	}
	restMapping, err := c.restMapper.RESTMapping(gk, versions...)
	if err != nil {
		return restMapping, err
	}
	c.gvkToGVR.Store(gvk, restMapping)
	return restMapping, nil
}

func (c *cacheRESTMapper) RESTMappings(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error) {
	return c.restMapper.RESTMappings(gk, versions...)
}

func (c *cacheRESTMapper) ResourceSingularizer(resource string) (singular string, err error) {
	return c.restMapper.ResourceSingularizer(resource)
}

func NewCachedRESTMapper(cfg *rest.Config, underlyingMapper meta.RESTMapper) (meta.RESTMapper, error) {
	cachedMapper := cacheRESTMapper{}

	if underlyingMapper != nil {
		cachedMapper.restMapper = underlyingMapper
		return &cachedMapper, nil
	}
	client, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	option := apiutil.WithCustomMapper(func() (meta.RESTMapper, error) {
		groupResources, err := restmapper.GetAPIGroupResources(client)
		if err != nil {
			return nil, err
		}
		cachedMapper.gvkToGVR = sync.Map{}
		cachedMapper.gvrToGVK = sync.Map{}
		cachedMapper.gvrToGVKErr = sync.Map{}
		return restmapper.NewDiscoveryRESTMapper(groupResources), nil
	})

	underlyingMapper, err = apiutil.NewDynamicRESTMapper(cfg, option)
	if err != nil {
		return nil, err
	}
	cachedMapper.restMapper = underlyingMapper
	return &cachedMapper, nil
}

func Provider(cfg *rest.Config) (meta.RESTMapper, error) {
	return NewCachedRESTMapper(cfg, nil)
}
