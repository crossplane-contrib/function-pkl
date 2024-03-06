package main

import (
	"context"
	"fmt"
	"os"

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
		xr, err := request.GetObservedCompositeResource(req)
		myxr := req.Observed.GetComposite()
	*/
	evaluator, err := pkl.NewEvaluator(ctx, pkl.PreconfiguredOptions,
		WithCrossplane(req, "crossplane"),
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

	rawYaml, err := yaml.Marshal(req.Observed.Composite.GetResource())
	compositionYamlFile := fmt.Sprintf("%s/observed/%s", tempDir, "composition.yaml")
	os.WriteFile(compositionYamlFile, rawYaml, 0666)

	// turn yaml into  pkl file:
	convertEvalManager := pkl.NewEvaluatorManagerWithCommand([]string{"/home/tim/.local/bin/pkl", "-p", "input=" + compositionYamlFile}) // TODO fix path
	convertEval, err := convertEvalManager.NewEvaluator(ctx, pkl.PreconfiguredOptions)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "error creating new evaluater"))
		return rsp, err
	}
	x, err := convertEval.EvaluateOutputFiles(ctx, pkl.UriSource("package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl"))
	if len(x) > 10 {
		response.Fatal(rsp, errors.Wrap(err, "this is a test"))
		return rsp, err
	}
	//os.WriteFile(fmt.Sprintf("%s/observed/%s", tempDir, "composition.pkl"))
	// turn input yaml to pkl file
	// for each observed resource create a file
	//os.WriteFile()

	//req.Observed.Composite

	for name, source := range sources {
		resource, err := parseFile(ctx, evaluator, source)
		if err != nil {
			response.Fatal(rsp, errors.Wrap(err, "error during parsing of file"))
		}
		outResources[name] = resource
	}

	rsp.Desired.Resources = outResources
	return rsp, nil
}

/*

// generates a Pkl File from yaml
func createPklFile(name string) error {
	// to yaml and then to pkl
	rawYaml, err := yaml.Marshal(r.GetResource())
	if err != nil {
		return err
	}

	//string(rawYaml)
	// try to run evaluator with command?
	// pkl eval -p input=resource.yaml -o resource.pkl package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl

	os.WriteFile(name)
}*/

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
