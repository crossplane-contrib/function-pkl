package function

import (
	"context"
	"fmt"
	"os"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/avarei/function-pkl/internal/pkl/reader"
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

	Log logging.Logger
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
	/*
		xr, err := request.GetObservedCompositeResource(req)
		myxr := req.Observed.GetComposite()
	*/
	evaluatorManager := pkl.NewEvaluatorManager()
	defer evaluatorManager.Close()
	evaluator, err := evaluatorManager.NewEvaluator(ctx, pkl.PreconfiguredOptions,
		reader.WithCrossplane(req, "crossplane"),
	) // TODO disallow FS access
	if err != nil {
		evaluator.Close()
		response.Fatal(rsp, errors.Wrap(err, "could not create Pkl Evaluater"))
	}
	defer evaluator.Close()

	var outResources map[string]*fnv1beta1.Resource = make(map[string]*fnv1beta1.Resource)

	var sources map[string]*pkl.ModuleSource = make(map[string]*pkl.ModuleSource)

	for fileName, fileContent := range in.Spec.Files {
		sources[fileName] = pkl.TextSource(fileContent)
	}

	for i, uri := range in.Spec.Uris {
		sources[fmt.Sprintf("uri-%d", i+1)] = pkl.UriSource(uri)
	}

	// TODO configMap

	tempDir, err := os.MkdirTemp("", "pkl-run-")
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "error creating temporary directory"))
		return rsp, err
	}
	defer os.RemoveAll(tempDir)
	err = os.Mkdir(fmt.Sprintf("%s/%s", tempDir, "observed"), 0777)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "error creating observed directory"))
		return rsp, err
	}

	for name, source := range sources {
		resource, err := parseFile(ctx, evaluator, source)
		if err != nil {
			fmt.Print(err)
			response.Fatal(rsp, errors.Wrap(err, "error during parsing of file"))
		}
		outResources[name] = resource
	}

	rsp.Desired.Resources = outResources
	return rsp, nil
}

func parseFile(ctx context.Context, evaluator pkl.Evaluator, source *pkl.ModuleSource) (*fnv1beta1.Resource, error) {
	// TODO request a new Function to EvaluateOutputValue which does not require a Struct Tag
	out, err := evaluator.EvaluateOutputText(ctx, source)
	if err != nil {
		return nil, errors.Wrap(err, "could not evaluate Pkl file")
	}
	var x map[string]any
	if err := yaml.Unmarshal([]byte(out), &x); err != nil {
		return nil, errors.Wrap(err, "could not parse yaml to map[string]any")
	}

	st, err := structpb.NewStruct(x)
	if err != nil {
		return nil, errors.Wrap(err, "could not evaluate Pkl output as map with keys and values")
	}

	return &fnv1beta1.Resource{
		Resource: st,
	}, nil

}
