package reader

import (
	"context"
	"fmt"
	"net/url"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/crossplane/function-sdk-go/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"sigs.k8s.io/yaml"
)

type CrossplaneReader struct {
	Request      *fnv1beta1.RunFunctionRequest
	ReaderScheme string
	Log          logging.Logger
	Ctx          context.Context
}

func (f *CrossplaneReader) Scheme() string {
	return f.ReaderScheme
}

func (f *CrossplaneReader) IsGlobbable() bool {
	return false
}

func (f *CrossplaneReader) HasHierarchicalUris() bool {
	return false
}

// ListElements returns the list of elements at a specified path.
// If HasHierarchicalUris is false, path will be empty and ListElements should return all
// available values.
//
// This method is only called if it is hierarchical and local, or if it is globbable.
func (f *CrossplaneReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	out := []pkl.PathElement{
		pkl.NewPathElement("state", false),
		pkl.NewPathElement("input", false),
		pkl.NewPathElement("crds", false),
	}
	return out, nil
}

var _ pkl.Reader = (*CrossplaneReader)(nil)

type crossplaneModuleReader struct {
	*CrossplaneReader
}

func (f crossplaneModuleReader) IsLocal() bool {
	return true
}

var evaluatorManager pkl.EvaluatorManager = pkl.NewEvaluatorManager()

// TODO find better solution
func Close() error {
	return evaluatorManager.Close()
}

var WithCrossplane = func(crossplaneReader *CrossplaneReader) func(opts *pkl.EvaluatorOptions) {
	return func(opts *pkl.EvaluatorOptions) {
		reader := crossplaneReader
		pkl.WithModuleReader(&crossplaneModuleReader{reader})(opts)
		pkl.WithResourceReader(&crossplaneResourceReader{reader})(opts)
	}
}

func (f CrossplaneReader) BaseRead(url url.URL) ([]byte, error) {
	switch url.Opaque {
	case "state":
		evaluator, err := evaluatorManager.NewEvaluator(
			f.Ctx,
			pkl.PreconfiguredOptions,
			WithCrossplane(&CrossplaneReader{
				Request:      f.Request,
				ReaderScheme: "crossplane",
				Log:          nil,
			}), // TODO: This should be a seperate reader Implementation, as calling crossplane:state within crossplane:state would softlock and should not be allowed
		)
		if err != nil {
			return nil, err
		}
		defer evaluator.Close()

		out, err := evaluator.EvaluateOutputText(f.Ctx, pkl.UriSource("package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.6#/convert.pkl"))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		fmt.Println(out)
		return []byte(out), err
	case "input":
		requestYaml, err := yaml.Marshal(f.Request)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(requestYaml))

		return requestYaml, nil
	case "crds":
		in := &v1beta1.Pkl{}
		if err := request.GetInput(f.Request, in); err != nil {
			return nil, err
		}

		resourceTemplates := make(map[string]map[string]string)
		for _, crd := range in.Spec.PklCRDs {
			if resourceTemplates[crd.Kind] == nil {
				resourceTemplates[crd.Kind] = map[string]string{
					crd.ApiVersion: crd.Uri,
				}
			} else {
				resourceTemplates[crd.Kind][crd.ApiVersion] = crd.Uri
			}
		}

		message := buildResourceTemplatesModule(resourceTemplates)
		fmt.Println(message)
		return []byte(message), nil
	default:
		return nil, fmt.Errorf("unsupported path")
	}
}

// generates a resourceTemplate similar to https://github.com/apple/pkl-k8s/blob/main/generated-package/k8sSchema.pkl but for the custom Resources
func buildResourceTemplatesModule(resourceTemplates map[string]map[string]string) string {
	message := "resourceTemplates: Mapping<String, Mapping<String, unknown>> = new {\n"
	for kind, versionUris := range resourceTemplates {
		message += fmt.Sprintf("  [\"%s\"] {\n", kind)
		for version, uri := range versionUris {
			message += fmt.Sprintf("    [\"%s\"] = import(\"%s\")\n", version, uri)
		}
		message += "  }\n"
	}
	message += "}\n"
	return message
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
	*CrossplaneReader
}

func (f crossplaneResourceReader) Read(url url.URL) ([]byte, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _ pkl.ResourceReader = (*crossplaneResourceReader)(nil)
