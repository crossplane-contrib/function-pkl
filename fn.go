package main

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"
	"google.golang.org/protobuf/types/known/structpb"
	"sigs.k8s.io/yaml"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(ctx context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	f.log.Info("Running function", "tag", req.GetMeta().GetTag())

	rsp := response.To(req, response.DefaultTTL)

	in := &v1beta1.Pkl{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Function input from %T", req))
		return rsp, nil
	}
	/*
		// TODO: Add your Function logic here!
		response.Normalf(rsp, "I was run with input %q!", in.Example)
		f.log.Info("I was run!", "input", in.Example)
	*/

	/*
		compositeResource, err := request.GetObservedCompositeResource(req)
		if err != nil {
			response.Fatal(rsp, errors.Wrap(err, "could not get Composite Resource from Request"))
		}
	*/

	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions) // TODO disallow FS access
	if err != nil {
		evaluator.Close()
		response.Fatal(rsp, errors.Wrap(err, "could not create Pkl Evaluater"))
	}
	defer evaluator.Close()

	var source *pkl.ModuleSource

	switch in.Spec.Source {
	case "inline":
		source = pkl.TextSource(in.Spec.Inline)
	case "configMap":
		// TODO get configMap (maybe with informer pattern to cache it) and use it's content
		response.Fatal(rsp, errors.Cause(errors.New("not yet implemented")))
	case "uri":
		source = pkl.UriSource(in.Spec.Uri)
	}

	// TODO request a new Function to EvaluateOutputValue which does not require a Struct Tag
	out, err := evaluator.EvaluateOutputText(ctx, source)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not evaluate Pkl file"))
	}
	var x map[string]any
	if err := yaml.Unmarshal([]byte(out), &x); err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not parse yaml to map[string]any"))
	}

	var outResources = make(map[string]*fnv1beta1.Resource)
	st, err := structpb.NewStruct(x)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not evaluate Pkl output as map with keys and values"))
	}
	outResources["foo"] = &fnv1beta1.Resource{
		Resource: st,
	}

	rsp.Desired.Resources = outResources
	return rsp, nil
}
