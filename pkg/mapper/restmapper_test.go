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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

func getConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/zhousong/admin.conf")
	if err != nil {
		panic(err)
	}
	return config
}

var podGVR = schema.GroupVersionResource{
	Group:    "",
	Version:  "v1",
	Resource: "pods",
}
var nodeImagesGVR = schema.GroupVersionResource{
	Group:    "apps.kruise.io",
	Version:  "v1alpha1",
	Resource: "nodeimages",
}
var nodePoolsGVR = schema.GroupVersionResource{
	Group:    "apps.openyurt.io",
	Version:  "v1alpha1",
	Resource: "nodepools",
}

func Test_cacheRESTMapper_KindFor(t *testing.T) {
	config := getConfig()
	mapper, err := Provider(config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	tests := []struct {
		name     string
		resource schema.GroupVersionResource
		wantErr  bool
	}{
		{
			name:     "std pod gvr",
			resource: podGVR,
			wantErr:  false,
		},
		{
			name:     "std crd gvr",
			resource: nodeImagesGVR,
			wantErr:  false,
		},
		{
			name:     "err crd gvr",
			resource: nodePoolsGVR,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mapper
			_, err = c.KindFor(tt.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("KindFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func BenchmarkRESTMapper_KindFor(b *testing.B) {
	config := getConfig()
	mapper, err := Provider(config)
	if err != nil {
		b.Errorf("unexpected error: %v", err)
	}
	tests := []struct {
		name     string
		resource schema.GroupVersionResource
		wantErr  bool
	}{
		{
			name:     "std pod gvr",
			resource: podGVR,
			wantErr:  false,
		},
		{
			name:     "std crd gvr",
			resource: nodeImagesGVR,
			wantErr:  false,
		},
		{
			name:     "err crd gvr",
			resource: nodePoolsGVR,
			wantErr:  true,
		},
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, tt := range tests {
				b.Run(tt.name, func(b *testing.B) {
					c := mapper
					_, err = c.KindFor(tt.resource)
					if (err != nil) != tt.wantErr {
						b.Errorf("KindFor() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				})
			}
		}
	})
}
