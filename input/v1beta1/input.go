// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=template.fn.crossplane.io
// +versionName=v1beta1
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// Input can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type Pkl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec PklSpec `json:"spec,omitempty"`
}

type PklSpec struct {
	// Pkl Template of the CompositeResourceDefinition (XRD), which will be amended by the CompositeResource (XR)
	XrdTemplate string `json:"xrdTemplate,omitempty"`

	// Source from which the Project is imported
	// +kubebuilder:validation:Enum=inline;configMap;uri
	Source string `json:"source,omitempty"`
	// Contains a stringified Pkl file
	Inline string `json:"inline,omitempty"`

	// Use URI Scheme to load Project/Package
	Uri string `json:"uri,omitempty"`

	// Load Project/Package from ConfigMap. Will evaluate PklProject and *.pkl files within the ConfigMap.
	ConfigMapRef string `json:"configMapRef,omitempty"`
}

type ConfigMapRef struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
}
