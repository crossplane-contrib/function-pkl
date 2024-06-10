// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=template.fn.crossplane.io
// +versionName=v1beta1
package v1beta1

import (
	"fmt"
	"strings"

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
	Composition *PklFileRef `json:"composition,omitempty"`

	// +kubebuilder:validation:Required
	Resources []PklFileRef `json:"resources,omitempty"`

	Full *PklFileRef `json:"full,omitempty"`

	// Pkl Template of the CompositeResourceDefinition (XRD), which will be amended by the CompositeResource (XR)
	CRDs []PklCrdRef `json:"crds,omitempty"`

	Requirements *PklFileRef `json:"requirements,omitempty"`

	// Packages is a list of Pkl Packages that can be used as a shorthand for the full package Path. This is similar to PklProject dependencies
	Packages []Package `json:"packages,omitempty"`
}

func (p PklSpec) ParseUri(uri string) string {
	if !strings.Contains(uri, "@") {
		return uri
	}
	for _, v := range p.Packages {
		if filePath, found := strings.CutPrefix(uri, "@"+v.Name); found {
			return fmt.Sprintf("%s#%s", v.Uri, filePath)
		}
	}

	// If no match was found try the full path
	return uri
}

type Package struct {
	Name string `json:"name,omitempty"`

	// Core specifies this packages as the one
	// providing essential capabilities for this function
	Core bool `json:"core,omitempty"`

	Uri string `json:"uri,omitempty"`
}

type PklCrdRef struct {
	// Use URI Scheme to load CRD Template
	Uri string `json:"uri,omitempty"`

	// +kubebuilder:validation:Required
	Kind string `json:"kind,omitempty"`

	// +kubebuilder:validation:Required
	ApiVersion string `json:"apiVersion,omitempty"`
}

type PklFileRef struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Enum=uri;inline
	Type string `json:"type,omitempty"`

	// Use URI Scheme to load Project/Package
	Uri string `json:"uri,omitempty"`
	// Contains a stringified Pkl file
	Inline string `json:"inline,omitempty"`
}
