package function

import (
	"context"
	"fmt"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/avarei/function-pkl/internal/pkl/reader"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	"sigs.k8s.io/yaml"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	Log              logging.Logger
	EvaluatorManager pkl.EvaluatorManager
}

// RunFunction runs the Function.
func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.Log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1.Pkl{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}
	evaluator, err := f.EvaluatorManager.NewEvaluator(ctx,
		pkl.PreconfiguredOptions,
		reader.WithCrossplane(&reader.CrossplaneReader{
			ReaderScheme: "crossplane",
			Request:      req,
			Log:          f.Log,
			Ctx:          ctx,
		}),
	)

	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not create Pkl Evaluater"))
	}
	defer evaluator.Close()

	var outResources map[string]*fnv1beta1.Resource = make(map[string]*fnv1beta1.Resource)

	for _, pklFileRef := range in.Spec.PklManifests {

		fileName, moduleSource, err := evalFileRef(pklFileRef)
		if err != nil {
			return nil, err
		}
		resource, err := evalPklFile(ctx, fileName, moduleSource, evaluator)
		if err != nil {
			return nil, err
		}

		outResources[fileName] = resource

	}
	rsp.Desired.Resources = outResources

	if in.Spec.PklComposition != nil {
		fileName, moduleSource, err := evalFileRef(*in.Spec.PklComposition)
		if err != nil {
			return nil, err
		}
		resource, err := evalPklFile(ctx, fileName, moduleSource, evaluator)
		if err != nil {
			return nil, err
		}

		rsp.Desired.Composite = resource
	}

	//response.Fatal(rsp, err)
	// TODO add rsp.Results
	return rsp, nil
}

func evalFileRef(pklFileRef v1beta1.PklFileRef) (string, *pkl.ModuleSource, error) {
	switch pklFileRef.Type {
	case "uri":
		if pklFileRef.Uri == "" {
			return "", nil, fmt.Errorf("manifest type of \"%s\" is uri but uri is empty", pklFileRef.Name)
		}
		return pklFileRef.Name, pkl.UriSource(pklFileRef.Uri), nil

	case "inline":
		if pklFileRef.Inline == "" {
			return "", nil, fmt.Errorf("manifest type of \"%s\" is inline but inline is empty", pklFileRef.Name)
		}
		return pklFileRef.Name, pkl.TextSource(pklFileRef.Inline), nil
	default:
		return "", nil, errors.New("unknown PklFileRef type")
	}
}

func evalPklFile(ctx context.Context, name string, source *pkl.ModuleSource, evaluator pkl.Evaluator) (*fnv1beta1.Resource, error) {
	renderedManifest, err := evaluator.EvaluateOutputText(ctx, source)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse Pkl file \"%s\"", name)
	}

	resource := &fnv1beta1.Resource{}
	if err := yaml.Unmarshal([]byte(renderedManifest), resource); err != nil {
		return nil, errors.Wrap(err, "could not parse yaml to Resource")
	}

	return resource, nil
}
