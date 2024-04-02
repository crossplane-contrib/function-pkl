package reader

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"sigs.k8s.io/yaml"
)

type crossplaneReader struct {
	request *fnv1beta1.RunFunctionRequest
	scheme  string
}

func (f *crossplaneReader) Scheme() string {
	return f.scheme
}

func (f *crossplaneReader) IsGlobbable() bool {
	return true
}

func (f *crossplaneReader) HasHierarchicalUris() bool {
	return false
}

// ListElements returns the list of elements at a specified path.
// If HasHierarchicalUris is false, path will be empty and ListElements should return all
// available values.
//
// This method is only called if it is hierarchical and local, or if it is globbable.
func (f *crossplaneReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	out := []pkl.PathElement{}

	if strings.HasPrefix(url.Opaque, "crds/") {
		selector := strings.TrimPrefix(url.Opaque, "crds/")
		in := &v1beta1.Pkl{}
		if err := request.GetInput(f.request, in); err != nil {
			return nil, err
		}

		if selector != "*" {
			return nil, errors.New("only crds/* is implemented as of now. please open an Issue of you need additional implementation")
		}

		for _, crd := range in.Spec.PklCRDs {
			out = append(out, pkl.NewPathElement(fmt.Sprintf("crds/%s", crd.Name), false))
		}

		return out, nil
	}

	out = []pkl.PathElement{
		pkl.NewPathElement("state", false),
		pkl.NewPathElement("input", false),
		pkl.NewPathElement("crds", true),
	}
	return out, nil
}

var _ pkl.Reader = (*crossplaneReader)(nil)

type crossplaneModuleReader struct {
	*crossplaneReader
}

func (f crossplaneModuleReader) IsLocal() bool {
	return true
}

var WithCrossplane = func(req *fnv1beta1.RunFunctionRequest, scheme string) func(opts *pkl.EvaluatorOptions) {
	return func(opts *pkl.EvaluatorOptions) {
		reader := &crossplaneReader{request: req, scheme: scheme}
		pkl.WithModuleReader(&crossplaneModuleReader{reader})(opts)
		pkl.WithResourceReader(&crossplaneResourceReader{reader})(opts)
	}
}

func (f crossplaneReader) BaseRead(url url.URL) ([]byte, error) {
	path := strings.TrimSuffix(strings.TrimPrefix(url.Opaque, "/"), "/")
	pathElements := strings.Split(path, "/")
	switch pathElements[0] {
	case "state":
		evaluatorManager := pkl.NewEvaluatorManager()
		defer evaluatorManager.Close()
		evaluator, err := evaluatorManager.NewEvaluator(
			context.TODO(),
			pkl.PreconfiguredOptions,
			WithCrossplane(f.request, "crossplane"),
		)
		if err != nil {
			return nil, err
		}

		out, err := evaluator.EvaluateOutputText(context.TODO(), pkl.UriSource("https://raw.githubusercontent.com/Avarei/function-pkl/main/pkl/convert.pkl")) // TODO find better solution
		fmt.Println(err)
		return []byte(out), err
	case "input":
		requestYaml, err := yaml.Marshal(f.request)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(requestYaml))
		return requestYaml, nil
	case "crds":
		if len(pathElements) != 2 {
			return nil, fmt.Errorf("expected exactly one CRD name")
		}
		crdName := pathElements[1]
		in := &v1beta1.Pkl{}
		if err := request.GetInput(f.request, in); err != nil {
			return nil, err
		}
		for _, crd := range in.Spec.PklCRDs {
			if crd.Name == crdName {
				return []byte(crd.Inline), nil
			}
		}
	}
	return nil, fmt.Errorf("path not found")
}

// Expects an URL like /observed/composition/resource and evaluates the RunFunctionRequest for the state of the desired field and returns it as a pkl file
func (f crossplaneModuleReader) Read(url url.URL) (string, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

var _ pkl.ModuleReader = (*crossplaneModuleReader)(nil)

type crossplaneResourceReader struct {
	*crossplaneReader
}

func (f crossplaneResourceReader) Read(url url.URL) ([]byte, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _ pkl.ResourceReader = (*crossplaneResourceReader)(nil)
