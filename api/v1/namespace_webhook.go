/*
Copyright 2023 Karthikeyan Govindaraj <github.gkarthiks@gmail.com>.

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

package v1

import (
	"context"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/mutate-core-v1-namespace,mutating=true,failurePolicy=fail,sideEffects=None,groups=core,resources=namespaces,verbs=create;update,versions=v1,name=namespace-annotate.github.gkarthiks.io,admissionReviewVersions=v1

var (
	namespacelog                          = logf.Log.WithName("namespace-annotator")
	_            webhook.AdmissionHandler = &namespaceAnnotator{}
)

// namespaceAnnotator struct used to handle admission control for Kubernetes namespaces
type namespaceAnnotator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// NewNamespaceAnnotator function creates a new instance of a namespace annotator
func NewNamespaceAnnotator(client client.Client, scheme *runtime.Scheme) admission.Handler {
	return &namespaceAnnotator{
		Client:  client,
		decoder: admission.NewDecoder(scheme),
	}
}

func (n *namespaceAnnotator) Handle(ctx context.Context, request admission.Request) admission.Response {

	namespacelog.Info("handling the namespace CREATE/UPDATE event")

	namespace := &v1.Namespace{}

	err := n.decoder.Decode(request, namespace)
	if err != nil {
		namespacelog.Error(err, "error occurred while decoding the admission request")
		return admission.Errored(http.StatusBadRequest, err)
	}

	namespacelog.Info("handling the namespace CREATE/UPDATE event for", "namespace", namespace.Name)

	if namespace.Annotations == nil {
		namespace.Annotations = make(map[string]string)
	}

	namespace.Annotations["githu.gkarthiks.io/annotation"] = "added"

	marshaledNamespace, err := json.Marshal(namespace)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	namespacelog.Info("the following namespace has been annotated successfully", "namespace", namespace.Name)
	return admission.PatchResponseFromRaw(request.AdmissionRequest.Object.Raw, marshaledNamespace)
}
