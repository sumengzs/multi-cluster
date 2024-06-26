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

package controllers

import (
	"context"
	"github.com/sumengzs/multi-cluster/pkg/cluster"
	"github.com/sumengzs/multi-cluster/pkg/pool"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/sumengzs/multi-cluster/api/v1beta1"
)

// ClusterController reconciles a Cluster object
type ClusterController struct {
	client.Client
	Pool   pool.Interface
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sumengzs.cn,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sumengzs.cn,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sumengzs.cn,resources=clusters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Cluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ClusterController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Cluster{}).
		WithEventFilter(r.Predicate()).
		Complete(r)
}

func (r *ClusterController) Predicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(event event.DeleteEvent) bool {
			clu := event.Object.(*v1beta1.Cluster)
			r.Pool.Remove(clu.Name)
			return false
		},
		UpdateFunc: func(event event.UpdateEvent) bool {
			return false
		},
		CreateFunc: func(event event.CreateEvent) bool {
			clu := event.Object.(*v1beta1.Cluster)
			cc, err := cluster.
				By(r.Client).
				WithScheme(r.Scheme).
				Named(clu.Name).
				WithOptions().
				Complete()
			if err != nil {
				klog.Errorf("error creating cluster %s: %v", clu.Name, err)
				return false
			}
			err = r.Pool.Add(cc)
			if err != nil {
				klog.Errorf("error add cluster to pool %s: %v", cc.Name(), err)
				r.Pool.Remove(cc.Name())
				return false
			}
			return false
		},
	}
}
