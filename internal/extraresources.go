package internal

import (
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

type ExtraResourceSelectors struct {
	ExtraResources map[string]*ExtraResourceSelector `json:"extraResourceSelectors,omitempty"`
}

func (e *ExtraResourceSelectors) ToResourceSelectors() map[string]*fnv1beta1.ResourceSelector {
	out := map[string]*fnv1beta1.ResourceSelector{}
	for name, extraResourceSelector := range e.ExtraResources {
		out[name] = extraResourceSelector.ToResourceSelector()
	}
	return out
}

type ExtraResourceSelector struct {
	ApiVersion  string            `json:"apiVersion"`
	Kind        string            `json:"kind"`
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
	Name        string            `json:"name,omitempty"`
}

func (e *ExtraResourceSelector) ToResourceSelector() *fnv1beta1.ResourceSelector {
	out := &fnv1beta1.ResourceSelector{
		ApiVersion: e.ApiVersion,
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
