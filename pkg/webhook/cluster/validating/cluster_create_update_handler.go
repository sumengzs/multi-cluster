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

package validating

import (
	"context"
	"github.com/sumengzs/multi-cluster/api/v1beta1"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var _ admission.Handler = &ClusterCreateUpdateHandler{}
var _ admission.DecoderInjector = &ClusterCreateUpdateHandler{}
var _ inject.Client = &ClusterCreateUpdateHandler{}

// ClusterCreateUpdateHandler handles Cluster
type ClusterCreateUpdateHandler struct {
	// To use the client, you need to do the following:
	// - uncomment it
	// - import sigs.k8s.io/controller-runtime/pkg/client
	// - uncomment the InjectClient method at the bottom of this file.
	Client client.Client

	// Decoder decodes objects
	Decoder *admission.Decoder
}

// Handle handles admission requests.
func (h *ClusterCreateUpdateHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	switch req.AdmissionRequest.Operation {
	case admissionv1.Create:
		obj := &v1beta1.Cluster{}
		if err := h.Decoder.Decode(req, obj); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		errList := h.validateCluster(obj)
		if len(errList) != 0 {
			klog.ErrorS(errList.ToAggregate(), "Invalid cluster error")
			return admission.Errored(http.StatusUnprocessableEntity, errList.ToAggregate())
		}
	case admissionv1.Update:
		obj := &v1beta1.Cluster{}
		if err := h.Decoder.Decode(req, obj); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		oldObj := &v1beta1.Cluster{}
		if err := h.Decoder.DecodeRaw(req.AdmissionRequest.OldObject, oldObj); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
		errList := h.validateClusterUpdate(oldObj, obj)
		if len(errList) != 0 {
			klog.ErrorS(errList.ToAggregate(), "Invalid cluster error")
			return admission.Errored(http.StatusUnprocessableEntity, errList.ToAggregate())
		}
	}
	klog.Infof("handle cluster create update request successfully")
	return admission.ValidationResponse(true, "")
}

func (h *ClusterCreateUpdateHandler) validateClusterUpdate(oldObj, newObj *v1beta1.Cluster) field.ErrorList {
	latestObject := &v1beta1.Cluster{}
	key := client.ObjectKeyFromObject(newObj)
	err := h.Client.Get(context.TODO(), key, latestObject)
	if err != nil {
		return field.ErrorList{field.InternalError(field.NewPath("cluster"), err)}
	}
	if errorList := h.validateCluster(newObj); errorList != nil {
		return errorList
	}

	return nil
}

func (h *ClusterCreateUpdateHandler) validateCluster(obj *v1beta1.Cluster) field.ErrorList {
	return h.validateClusterSpec(&obj.Spec, field.NewPath("Spec"))
}

func (h *ClusterCreateUpdateHandler) validateClusterSpec(spec *v1beta1.ClusterSpec, path *field.Path) field.ErrorList {
	return h.validateSpecConnectConfig(spec.Connect, path.Child("ConnectConfig"))
}

func (h *ClusterCreateUpdateHandler) validateSpecConnectConfig(config v1beta1.ConnectConfig, path *field.Path) field.ErrorList {
	if config.Secret == nil && config.Config == nil && config.Token == nil {
		return field.ErrorList{field.Invalid(path, "", "Secret, Config and Token cannot be empty as the same time")}
	}
	return nil
}

// InjectClient injects the client into the ClusterCreateUpdateHandler
func (h *ClusterCreateUpdateHandler) InjectClient(c client.Client) error {
	h.Client = c
	return nil
}

// InjectDecoder injects the decoder into the ClusterCreateUpdateHandler
func (h *ClusterCreateUpdateHandler) InjectDecoder(d *admission.Decoder) error {
	h.Decoder = d
	return nil
}
