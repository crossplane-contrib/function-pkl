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

package internal

import (
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

// ExtraResourceSelectors allows selecting resources
type ExtraResourceSelectors struct {
	ExtraResources map[string]*ExtraResourceSelector `json:"extraResourceSelectors,omitempty"`
}

// ToResourceSelectors converts this to the upstream format
func (e *ExtraResourceSelectors) ToResourceSelectors() map[string]*fnv1beta1.ResourceSelector {
	out := map[string]*fnv1beta1.ResourceSelector{}
	for name, extraResourceSelector := range e.ExtraResources {
		out[name] = extraResourceSelector.ToResourceSelector()
	}
	return out
}

// ExtraResourceSelector allows setting a selector to lookup ExtraResources
type ExtraResourceSelector struct {
	APIVersion  string            `json:"apiVersion"`
	Kind        string            `json:"kind"`
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
	Name        string            `json:"name,omitempty"`
}

// ToResourceSelector converts this to the upstream format
func (e *ExtraResourceSelector) ToResourceSelector() *fnv1beta1.ResourceSelector {
	out := &fnv1beta1.ResourceSelector{
		ApiVersion: e.APIVersion,
		Kind:       e.Kind,
	}
	if len(e.MatchLabels) == 0 {
		out.Match = &fnv1beta1.ResourceSelector_MatchName{
			MatchName: e.Name,
		}
		return out
	}

	out.Match = &fnv1beta1.ResourceSelector_MatchLabels{
		MatchLabels: &fnv1beta1.MatchLabels{Labels: e.MatchLabels},
	}
	return out
}
