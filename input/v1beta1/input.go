/*
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

// Pkl struct can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type Pkl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec PklSpec `json:"spec,omitempty"`
}

// PklSpec specifies references for the function
type PklSpec struct {
	// +kubebuilder:validation:Enum=uri;inline;local
	Type string `json:"type,omitempty"`

	// Use URI Scheme to load Project/Package
	URI string `json:"uri,omitempty"`
	// Contains a stringified Pkl file
	Inline string `json:"inline,omitempty"`

	// Reference to a Pklfile and Project
	Local *Local `json:"local,omitempty"`
}

// Local contains Reference to a Local Pkl Project and a Pkl file within it
type Local struct {
	// Path to file relative from the Project Dir
	File string `json:"file,omitempty"`
	// Path to the Project containing a Pklfile
	ProjectDir string `json:"projectDir,omitempty"`
}
