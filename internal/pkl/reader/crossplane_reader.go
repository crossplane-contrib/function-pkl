package reader

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/apple/pkl-go/pkl"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
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

// e.g. crossplane:/observed/composition/resource
func (f *crossplaneReader) HasHierarchicalUris() bool {
	return true
}

// e.g. crossplane:/observed/composition/
func (f *crossplaneReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	path := strings.TrimSuffix(strings.TrimPrefix(url.Path, "/"), "/")
	pathElements := strings.Split(path, "/")

	var state *fnv1beta1.State
	if len(pathElements) > 0 {
		switch pathElements[0] {
		case "context":
			return nil, fmt.Errorf("context is not yet implemented")
		case "observed":
			state = f.request.GetObserved()
		case "desired":
			state = f.request.GetDesired()
		default:
			return nil, fmt.Errorf("unexpected state type: %s", pathElements[0])
		}
	} else {
		return []pkl.PathElement{
			pkl.NewPathElement("context", true),
			pkl.NewPathElement("observed", true),
			pkl.NewPathElement("desired", true),
		}, nil
	}

	pathElements = pathElements[1:]

	var resource *fnv1beta1.Resource
	// var isComposition = false
	if len(pathElements) > 0 {
		switch pathElements[0] {
		case "composition":
			//isComposition = true
			resource = state.GetComposite()
		case "resources":
			if len(pathElements) > 1 {
				resource = state.GetResources()[pathElements[1]]
				pathElements = pathElements[1:]
			} else {
				var out []pkl.PathElement
				for name, _ := range state.GetResources() {
					out = append(out, pkl.NewPathElement(name, true))
				}
				return out, nil
			}
		default:
			return nil, fmt.Errorf("unexpected resource type: %s", pathElements[0])
		}
	} else {
		return []pkl.PathElement{
			pkl.NewPathElement("composition", true),
			pkl.NewPathElement("resources", true),
		}, nil
	}

	pathElements = pathElements[1:]

	if len(pathElements) > 0 {
		switch pathElements[0] {
		case "resource":
			return nil, fmt.Errorf("resource is not a directory")
		case "connectionDetails":
			return nil, fmt.Errorf("connectionDetails is not a directory")
		case "ready":
			return nil, fmt.Errorf("ready is not a directory")
		default:
			return nil, fmt.Errorf("unexpected resource type: %s", pathElements[0])
		}
	} else {
		out := []pkl.PathElement{}

		if resource.GetResource() != nil {
			out = append(out, pkl.NewPathElement("resource", false))
		}
		if resource.GetConnectionDetails() != nil {
			out = append(out, pkl.NewPathElement("connectionDetails", false))
		}
		if resource.GetReady() != fnv1beta1.Ready_READY_UNSPECIFIED {
			out = append(out, pkl.NewPathElement("ready", false))
		}
		return out, nil
	}
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
	path := strings.TrimSuffix(strings.TrimPrefix(url.Path, "/"), "/")
	pathElements := strings.Split(path, "/")

	var state *fnv1beta1.State
	switch pathElements[0] {
	case "context":
		return nil, fmt.Errorf("context is not yet implemented")
	case "observed":
		state = f.request.GetObserved()
	case "desired":
		state = f.request.GetDesired()
	default:
		return nil, fmt.Errorf("unexpected state type: %s", pathElements[0])
	}

	pathElements = pathElements[1:]

	var resource *fnv1beta1.Resource
	// var isComposition = false
	switch pathElements[0] {
	case "composition":
		//isComposition = true
		resource = state.GetComposite()
	case "resources":
		resource = state.GetResources()[pathElements[1]]
		pathElements = pathElements[1:]
	default:
		return nil, fmt.Errorf("unexpected resource type: %s", pathElements[0])
	}

	pathElements = pathElements[1:]
	switch pathElements[0] {
	case "resource":
		subResource := resource.GetResource().AsMap()
		// Convert subResource to Yaml

		yaml, err := yaml.Marshal(subResource)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		out, err := NewYamlManifestToPklFile(string(yaml), nil).SillyHack()
		if err != nil {
			return nil, err
		}
		fmt.Println(out)

		// Eval the pkl file // TODO use NewEvalutorManager
		evaluator, err := pkl.NewEvaluator(context.TODO(), pkl.PreconfiguredOptions)
		if err != nil {
			return nil, err
		}

		outy, err := evaluator.EvaluateOutputText(context.TODO(), pkl.TextSource(out))
		if err != nil {
			return nil, err
		}

		fmt.Println(outy)

		return []byte(outy), nil
	case "connectionDetails":
		// resource.GetConnectionDetails()
		return nil, fmt.Errorf("not implemented")
	case "ready":
		// resource.GetReady()
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("unexpected resource type: %s", pathElements[0])
	}
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

// TODO Implement
func (f crossplaneResourceReader) Read(url url.URL) ([]byte, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _ pkl.ResourceReader = (*crossplaneResourceReader)(nil)
