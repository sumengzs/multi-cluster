package cluster

import (
	"context"
	"fmt"
	"github.com/sumengzs/multi-cluster/api/v1beta1"
	"github.com/sumengzs/multi-cluster/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Builder struct {
	clusterName string
	master      client.Client
	scheme      *runtime.Scheme
	options     []InitOptions
}

func By(master client.Client) *Builder {
	return &Builder{master: master}
}

func (b *Builder) WithScheme(scheme *runtime.Scheme) *Builder {
	if scheme == nil {

	}
	b.scheme = scheme
	return b
}

func (b *Builder) WithOptions(opts ...InitOptions) *Builder {
	b.options = opts
	return b
}

func (b *Builder) Named(clusterName string) *Builder {
	b.clusterName = clusterName
	return b
}

func (b *Builder) Complete() (Interface, error) {
	if b.scheme == nil {
		return nil, fmt.Errorf("must provide a non-nil scheme")
	}
	clusterCR, err := b.loadClusterCR()
	if err != nil {
		return nil, fmt.Errorf("failed to load cluster resource: %s", err)
	}
	config, err := b.loadConfig(clusterCR.Spec.Connect)
	if err != nil {
		return nil, fmt.Errorf("failed to load client rest config: %s", err)
	}
	cluster, err := New(config, b.scheme, b.options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %s", err)
	}
	if clusterCR.Spec.Disabled {
		cluster.Disable()
	}
	return cluster, nil
}

func (b *Builder) loadClusterCR() (*v1beta1.Cluster, error) {
	if b.master == nil {
		return nil, fmt.Errorf("must provide a non-nil master cluster client")
	}
	if len(b.clusterName) == 0 {
		return nil, fmt.Errorf("must provide a non-empty cluster name")
	}
	return b.clusterGetter(b.clusterName)
}

func (b *Builder) loadConfig(connect v1beta1.ConnectConfig) (*rest.Config, error) {
	return utils.BuildConfig(b.clusterName, connect, b.secretGetter)
}

func (b *Builder) clusterGetter(name string) (*v1beta1.Cluster, error) {
	cluster := &v1beta1.Cluster{}
	err := b.master.Get(context.Background(), types.NamespacedName{Name: name}, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (b *Builder) secretGetter(key types.NamespacedName) (*v1.Secret, error) {
	secret := &v1.Secret{}
	err := b.master.Get(context.TODO(), key, secret)
	return secret, err
}
