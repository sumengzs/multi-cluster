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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// Provider of the cluster, this field is just for description
	// +optional
	Provider string `json:"provider,omitempty"`
	// Desired state of the cluster
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// Kubernetes API Server endpoint.
	// hostname:port, IP or IP:port.
	// Example: https://10.10.0.1:6443
	// +optional
	Endpoint string `json:"endpoint,omitempty"`
	// ProxyURL is the proxy URL for the cluster.
	// If not empty, the multi-cluster control plane will use this proxy to talk to the cluster.
	// More details please refer to: https://github.com/kubernetes/client-go/issues/351
	// +optional
	ProxyURL string `json:"proxyURL,omitempty"`
	// ProxyHeader is the HTTP header required by proxy server.
	// The key in the key-value pair is HTTP header key and value is the associated header payloads.
	// For the header with multiple values, the values should be separated by comma(e.g. 'k1': 'v1,v2,v3').
	// +optional
	ProxyHeader map[string]string `json:"proxyHeader,omitempty"`
	// KubeConfig content used to connect to cluster api server
	KubeConfig []byte `json:"kubeconfig,omitempty"`
	// Region represents the region of the member cluster locate in.
	// +optional
	Region Region `json:"region,omitempty"`
}

type Region struct {
	// Zone represents the zone of the member cluster locate in.
	// +optional
	Zone string `json:"zone,omitempty"`
	// Country represents the country of the member cluster locate in.
	// +optional
	Country string `json:"country,omitempty"`
	// Province represents the province of the member cluster locate in.
	// +optional
	Province string `json:"province,omitempty"`
	// City represents the city of the member cluster locate in.
	// +optional
	City string `json:"city,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// Version represents version of the member cluster.
	// +optional
	Version string `json:"version,omitempty"`

	// APIEnablements represents the list of APIs installed in the member cluster.
	// +optional
	APIEnablements []APIEnablement `json:"apiEnablements,omitempty"`

	// Conditions is an array of current cluster conditions.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// NodeSummary represents the summary of nodes status in the member cluster.
	// +optional
	NodeSummary *NodeSummary `json:"nodeSummary,omitempty"`
}

// APIEnablement is a list of API resource, it is used to expose the name of the
// resources supported in a specific group and version.
type APIEnablement struct {
	// GroupVersion is the group and version this APIEnablement is for.
	GroupVersion string `json:"groupVersion,omitempty"`

	// Resources is a list of APIResource.
	// +optional
	Resources []APIResource `json:"resources,omitempty"`
}

// APIResource specifies the name and kind names for the resource.
type APIResource struct {
	// Name is the plural name of the resource.
	// +required
	Name string `json:"name,omitempty"`

	// Kind is the kind for the resource (e.g. 'Deployment' is the kind for resource 'deployments')
	// +required
	Kind string `json:"kind,omitempty"`
}

// NodeSummary represents the summary of nodes status in a specific cluster.
type NodeSummary struct {
	// TotalNum is the total number of nodes in the cluster.
	// +optional
	TotalNum int32 `json:"total,omitempty"`

	// ReadyNum is the number of ready nodes in the cluster.
	// +optional
	ReadyNum int32 `json:"ready,omitempty"`
}

// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ENDPOINT",type="string",priority=1,JSONPath=".spec.endpoint",description="The cluster endpoint"
// +kubebuilder:printcolumn:name="ENABLE",type="boolean",priority=1,JSONPath=".spec.enabled",description="The cluster enable status"
// +kubebuilder:printcolumn:name="PROVIDER",type="string",priority=1,JSONPath=".spec.provider",description="The cluster provider"
// +kubebuilder:printcolumn:name="VERSION",type="string",JSONPath=".status.version",description="The cluster version"
// +kubebuilder:printcolumn:name="TOTAL",type="integer",JSONPath=".status.nodeSummary.total",description="The total number of node"
// +kubebuilder:printcolumn:name="READY",type="integer",JSONPath=".status.nodeSummary.ready",description="The ready number of node"
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=".metadata.creationTimestamp"

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterSpec `json:"spec,omitempty"`
	// +optional
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
