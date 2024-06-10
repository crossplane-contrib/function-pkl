package function

import (
	"context"
	"fmt"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/avarei/function-pkl/internal"
	"github.com/avarei/function-pkl/internal/helper"
	"github.com/avarei/function-pkl/internal/pkl/reader"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	"go.starlark.net/lib/proto"
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
	packages := helper.ParsePackages(in.Spec.Packages)

	evaluator, err := f.EvaluatorManager.NewEvaluator(ctx,
		pkl.PreconfiguredOptions,
		reader.WithCrossplane(&reader.CrossplaneReader{
			ReaderScheme: "crossplane",
			Request:      req,
			Log:          f.Log,
			Ctx:          ctx,
			Packages:     packages,
		}),
	)

	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not create Pkl Evaluater"))
		return rsp, nil
	}
	defer evaluator.Close()

	if fullRef := in.Spec.Full; fullRef != nil {
		fileName, moduleSource, err := evalFileRef(fullRef, packages)
		if err != nil {
			response.Fatal(rsp, errors.Wrap(err, "could not evaluate fileRef"))
			return rsp, nil
		}

		renderedManifest, err := evaluator.EvaluateOutputText(ctx, moduleSource)
		if err != nil {
			return nil, errors.Wrapf(err, "could not parse Pkl file \"%s\"", fileName)
		}

		out := &fnv1beta1.RunFunctionResponse{}
		err = proto.UnmarshalText(renderedManifest, out)

		return out, nil
	}

	var outResources map[string]*fnv1beta1.Resource = make(map[string]*fnv1beta1.Resource)

	for _, pklFileRef := range in.Spec.Resources {

		fileName, moduleSource, err := evalFileRef(&pklFileRef, packages)
		if err != nil {
			response.Fatal(rsp, errors.Wrap(err, "could not evaluate fileRef"))
			return rsp, nil
		}
		resource, err := evalPklFile(ctx, fileName, moduleSource, evaluator)
		if err != nil {
			response.Fatal(rsp, errors.Wrap(err, "could not evaluate Pkl File"))
			return rsp, nil
		}

		outResources[fileName] = resource
	}
	if rsp.Desired == nil {
		rsp.Desired = &fnv1beta1.State{}
	}
	rsp.Desired.Resources = outResources

	if in.Spec.Requirements != nil {
		fileName, moduleSource, err := evalFileRef(in.Spec.Requirements, packages)
		if err != nil {
			return nil, err
		}
		extraResources, err := evalExtraResources(ctx, fileName, moduleSource, evaluator)
		if err != nil {
			return nil, err
		}
		if len(extraResources) > 0 {
			rsp.Requirements = &fnv1beta1.Requirements{
				ExtraResources: extraResources,
			}
		}
	}

	if in.Spec.Composition != nil {
		fileName, moduleSource, err := evalFileRef(in.Spec.Composition, packages)
		if err != nil {
			return nil, err
		}
		resource, err := evalPklFile(ctx, fileName, moduleSource, evaluator)
		if err != nil {
			return nil, err
		}

		rsp.Desired.Composite = resource
	}

	// TODO add rsp.Results
	return rsp, nil
}

func evalFileRef(pklFileRef *v1beta1.PklFileRef, packages *helper.Packages) (string, *pkl.ModuleSource, error) {
	if pklFileRef == nil {
		return "", nil, errors.New("pklFileRef is nil")
	}
	switch pklFileRef.Type {
	case "uri":
		if pklFileRef.Uri == "" {
			return "", nil, fmt.Errorf("manifest type of \"%s\" is uri but uri is empty", pklFileRef.Name)
		}
		return pklFileRef.Name, pkl.UriSource(packages.ParseUri(pklFileRef.Uri)), nil

	case "inline":
		if pklFileRef.Inline == "" {
			return "", nil, fmt.Errorf("manifest type of \"%s\" is inline but inline is empty", pklFileRef.Name)
		}
		// TODO implement packages support @example -> uri
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

func evalExtraResources(ctx context.Context, name string, source *pkl.ModuleSource, evaluator pkl.Evaluator) (map[string]*fnv1beta1.ResourceSelector, error) {
	renderedManifest, err := evaluator.EvaluateOutputText(ctx, source)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse Pkl file \"%s\"", name)
	}

	resources := &internal.ExtraResourceSelectors{}
	if err := yaml.Unmarshal([]byte(renderedManifest), resources); err != nil {
		return nil, errors.Wrap(err, "could not parse yaml to Resource")
	}

	return resources.ToResourceSelectors(), nil
}
