package reader

import (
	"fmt"
	"net/url"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
)

type crdReader struct {
	request *fnv1beta1.RunFunctionRequest
	scheme  string
}

func (f *crdReader) Scheme() string {
	return f.scheme
}

func (f *crdReader) IsGlobbable() bool {
	return false
}

// e.g. crd:example
func (f *crdReader) HasHierarchicalUris() bool {
	return false
}

// e.g. crd:
func (f *crdReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	in := &v1beta1.Pkl{}
	if err := request.GetInput(f.request, in); err != nil {
		return nil, err
	}

	// create list of pathElement
	var ret []pkl.PathElement
	for _, pklCrdRef := range in.Spec.PklCrds {
		ret = append(ret, pkl.NewPathElement(pklCrdRef.Name, false))
	}
	return ret, nil
}

var _ pkl.Reader = (*crdReader)(nil)

type crdModuleReader struct {
	*crdReader
}

// triple dot notation does not make sense here eventhough this reader does not use remote resources
func (f crdModuleReader) IsLocal() bool {
	return false
}

var WithCrd = func(req *fnv1beta1.RunFunctionRequest, scheme string) func(opts *pkl.EvaluatorOptions) {
	return func(opts *pkl.EvaluatorOptions) {
		reader := &crdReader{request: req, scheme: scheme}
		pkl.WithModuleReader(&crdModuleReader{reader})(opts)
		pkl.WithResourceReader(&crdResourceReader{reader})(opts)
	}
}

func (f crdReader) BaseRead(url url.URL) (string, error) {
	in := &v1beta1.Pkl{}
	if err := request.GetInput(f.request, in); err != nil {
		return "", err
	}

	for _, pklCrdRef := range in.Spec.PklCrds {
		if pklCrdRef.Name != url.Path {
			continue
		}

		switch pklCrdRef.Type {
		case "inline":
			return pklCrdRef.Inline, nil
		default:
			return "", fmt.Errorf("unknown PklCrdRef type")
		}
	}
	return "", fmt.Errorf("PklCrdRef not found")
}

// Expects an URL like /observed/composition/resource and evaluates the RunFunctionRequest for the state of the desired field and returns it as a pkl file
func (f crdModuleReader) Read(url url.URL) (string, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return "", err
	}
	return out, nil
}

var _ pkl.ModuleReader = (*crdModuleReader)(nil)

type crdResourceReader struct {
	*crdReader
}

// TODO Implement
func (f crdResourceReader) Read(url url.URL) ([]byte, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return nil, err
	}
	return []byte(out), nil
}

var _ pkl.ResourceReader = (*crdResourceReader)(nil)
