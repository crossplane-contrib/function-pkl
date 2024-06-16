package function

import (
	"context"

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

type CompositionResponse struct {
	fnv1beta1.RunFunctionResponse

	Requirements *Requirements `json:"requirements,omitempty"`
}

type Requirements struct {
	//fnv1beta1.Requirements

	ExtraResources map[string]*ResourceSelector `protobuf:"bytes,1,rep,name=extra_resources,json=extraResources,proto3" json:"extraResources,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type ResourceSelector struct {
	ApiVersion string `protobuf:"bytes,1,opt,name=api_version,json=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Match      Match  `json:"match,omitempty"`
}

type Match struct {
	MatchName   string                 `protobuf:"bytes,3,opt,name=match_name,json=matchName,proto3,oneof" json:"matchName,omitempty"`
	MatchLabels *fnv1beta1.MatchLabels `protobuf:"bytes,4,opt,name=match_labels,json=matchLabels,proto3,oneof" json:"matchLabels,omitempty"`
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
		return rsp, nil
	}
	defer evaluator.Close()

	moduleSource, err := getModuleSource(in.Spec)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not evaluate fileRef"))
		return rsp, nil
	}

	renderedManifest, err := evaluator.EvaluateOutputText(ctx, moduleSource)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse Pkl file")
	}

	helper := &CompositionResponse{}
	err = yaml.Unmarshal([]byte(renderedManifest), helper)
	if err != nil {
		return nil, errors.Wrapf(err, "could not unmarshal pkl result")
	}

	fixedRequirements := &fnv1beta1.Requirements{
		ExtraResources: convertExtraResources(helper.Requirements.ExtraResources),
	}

	// Note: consider not overwriting rsp and whether it makes a difference.
	rsp = &fnv1beta1.RunFunctionResponse{
		Meta:         helper.Meta,
		Desired:      helper.Desired,
		Results:      helper.Results,
		Context:      helper.Context,
		Requirements: fixedRequirements,
	}

	return rsp, nil
}

func convertExtraResources(extraResources map[string]*ResourceSelector) map[string]*fnv1beta1.ResourceSelector {
	out := make(map[string]*fnv1beta1.ResourceSelector)
	for name, fixedrs := range extraResources {
		rs := &fnv1beta1.ResourceSelector{
			ApiVersion: fixedrs.ApiVersion,
			Kind:       fixedrs.Kind,
		}
		if fixedrs.Match.MatchLabels != nil && len(fixedrs.Match.MatchLabels.Labels) > 0 {
			rs.Match = &fnv1beta1.ResourceSelector_MatchLabels{
				MatchLabels: &fnv1beta1.MatchLabels{
					Labels: fixedrs.Match.MatchLabels.GetLabels(),
				},
			}
		} else {
			rs.Match = &fnv1beta1.ResourceSelector_MatchName{
				MatchName: fixedrs.Match.MatchName,
			}
		}
		out[name] = rs
	}
	return out
}

func getModuleSource(pklSpec v1beta1.PklSpec) (*pkl.ModuleSource, error) {
	switch pklSpec.Type {
	case "uri":
		if pklSpec.Uri == "" {
			return nil, errors.New("manifest type is uri but uri is empty")
		}
		return pkl.UriSource(pklSpec.Uri), nil

	case "inline":
		if pklSpec.Inline == "" {
			return nil, errors.New("manifest type is inline but inline is empty")
		}
		return pkl.TextSource(pklSpec.Inline), nil
	default:
		return nil, errors.New("unknown pklSpec type")
	}
}
