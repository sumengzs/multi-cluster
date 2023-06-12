package cluster

import (
	"github.com/sumengzs/multi-cluster/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

var (
	fakeScheme *runtime.Scheme
	fakeClient client.Client
	clusterCR  *v1beta1.Cluster
)

func init() {
	fakeScheme = runtime.NewScheme()
	v1beta1.AddToScheme(fakeScheme)
	clusterCR = GetMockCluster()
	fakeClient = fake.NewClientBuilder().WithScheme(fakeScheme).WithObjects(clusterCR).Build()
}

func GetMockCluster() *v1beta1.Cluster {
	return &v1beta1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: v1beta1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		Spec: v1beta1.ClusterSpec{
			Provider: "mock-provider",
			Disabled: false,
			Connect: v1beta1.ConnectConfig{
				Secret: &v1beta1.SecretRef{
					Namespace: "default",
					Name:      "test-secret",
				},
				Config: &v1beta1.ConfigRef{
					Secret: &v1beta1.SecretRef{
						Namespace: "default",
						Name:      "private-test-secret",
					},
					Config: []byte("XX"),
				},
				Token: &v1beta1.TokenRef{
					CABundle: nil,
					Token:    "A78DDS464Z",
				},
				InsecureSkipTLSVerification: false,
				Endpoint:                    "https://10.10.0.1:6443",
				ProxyURL:                    "",
				ProxyHeader:                 nil,
			},
			Region: v1beta1.Region{
				Zone:     "mock-zone",
				Country:  "mock-country",
				Province: "mock-province",
				City:     "mock-city",
			},
		},
		Status: v1beta1.ClusterStatus{
			Version: "1.19.2",
			APIEnablements: []v1beta1.APIEnablement{
				{
					GroupVersion: "apps/v1",
					Resources: []v1beta1.APIResource{
						{
							Name: "deployments",
							Kind: "Deployment",
						},
						{
							Name: "pods",
							Kind: "Pod",
						},
					},
				},
			},
			Conditions: []metav1.Condition{
				{
					Type:               "Ready",
					Status:             "True",
					LastTransitionTime: metav1.Now(),
				},
			},
			NodeSummary: &v1beta1.NodeSummary{
				TotalNum: 10,
				ReadyNum: 5,
			},
		},
	}
}

func TestBuilder_Complete(t *testing.T) {
	type fields struct {
		clusterName string
		master      client.Client
		scheme      *runtime.Scheme
	}
	tests := []struct {
		name    string
		fields  fields
		want    Interface
		wantErr bool
	}{
		{
			name: "test-cluster",
			fields: fields{
				clusterName: "test-cluster",
				master:      fakeClient,
				scheme:      fakeScheme,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := By(tt.fields.master).WithScheme(tt.fields.scheme).Named(tt.fields.clusterName)
			_, err := b.Complete()
			if (err != nil) != tt.wantErr {
				t.Errorf("Complete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
