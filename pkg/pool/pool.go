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
	"context"
	"fmt"
	"github.com/sumengzs/multi-cluster/pkg/cluster"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

type Pool struct {
	mu       sync.RWMutex
	client   client.Client
	clusters map[string]cluster.Interface
}

func New(config *rest.Config) (Interface, error) {
	clusters := make(map[string]cluster.Interface)
	config = rest.AddUserAgent(config, UserAgentName)

	cli, err := client.New(config, client.Options{
		Scheme: Scheme,
	})
	if err != nil {
		return nil, err
	}

	return &Pool{
		client:   cli,
		clusters: clusters,
	}, nil
}

func (p *Pool) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, clu := range p.clusters {
		if err := clu.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, clu := range p.clusters {
		clu.Stop()
	}
}

func (p *Pool) Add(clu cluster.Interface) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if oldC, ok := p.clusters[clu.Name()]; ok {
		switch oldC.Status() {
		case cluster.Started, cluster.Ready, cluster.Waiting:
			return fmt.Errorf("%s,can not replace", clu)
		}
	}
	p.clusters[clu.Name()] = clu
	return nil
}

func (p *Pool) Remove(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.clusters[name]; ok {
		p.clusters[name].Stop()
		delete(p.clusters, name)
	}
}

func (p *Pool) Cluster(name string) cluster.Interface {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.clusters[name]
}

func (p *Pool) Clusters() map[string]cluster.Interface {
	clusters := make(map[string]cluster.Interface, len(p.clusters))
	p.mu.RLock()
	defer p.mu.RUnlock()
	for name, clu := range p.clusters {
		clusters[name] = clu
	}
	return clusters
}
