/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package function package handles the gRPC calls
package function

import (
	"context"

	"github.com/apple/pkl-go/pkl"
	"google.golang.org/protobuf/types/known/durationpb"
	"sigs.k8s.io/yaml"

	"github.com/crossplane/crossplane-runtime/pkg/errors"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"

	"github.com/crossplane-contrib/function-pkl/input/v1beta1"
	"github.com/crossplane-contrib/function-pkl/internal/helper"
	"github.com/crossplane-contrib/function-pkl/internal/pkl/reader"
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

	moduleSource, err := getModuleSource(in.Spec)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "invalid composition function input"))
		return rsp, nil
	}

	var evaluator pkl.Evaluator

	switch in.Spec.Type {
	case "local":
		evaluator, err = f.EvaluatorManager.NewProjectEvaluator(ctx, in.Spec.Local.ProjectDir,
			pkl.PreconfiguredOptions,
			pkl.WithDefaultCacheDir,
			reader.WithCrossplane(&reader.CrossplaneReader{
				ReaderScheme: "crossplane",
				Request: &helper.CompositionRequest{
					RunFunctionRequest: req,
					ExtraResources:     req.GetExtraResources(),
				},
				Log: f.Log,
				Ctx: ctx,
			}))
	default:
		evaluator, err = f.EvaluatorManager.NewEvaluator(ctx,
			pkl.PreconfiguredOptions,
			pkl.WithDefaultCacheDir,
			reader.WithCrossplane(&reader.CrossplaneReader{
				ReaderScheme: "crossplane",
				Request: &helper.CompositionRequest{
					RunFunctionRequest: req,
					ExtraResources:     req.GetExtraResources(),
				},
				Log: f.Log,
				Ctx: ctx,
			}))
	}
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "could not create Pkl Evaluator"))
		return rsp, nil
	}
	defer func(evaluator pkl.Evaluator) {
		if err := evaluator.Close(); err != nil {
			f.Log.Info("evaluator could not be closed correctly:", err)
		}
	}(evaluator)

	renderedManifest, err := evaluator.EvaluateOutputText(ctx, moduleSource)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "error while parsing the Pkl file"))
		return rsp, nil
	}

	rsp, err = toResponse(renderedManifest, req.GetMeta().GetTag())
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "error while converting Pkl file"))
		return rsp, nil
	}

	return rsp, nil
}

func toResponse(renderedManifest, tag string) (*fnv1beta1.RunFunctionResponse, error) {
	rspHelper := &helper.CompositionResponse{}
	err := yaml.Unmarshal([]byte(renderedManifest), rspHelper)
	if err != nil {
		return nil, errors.Wrapf(err, "rendered Pkl file was not in expected format. did you amend @crossplane/CompositionResponse.pkl?")
	}

	responseMeta := &fnv1beta1.ResponseMeta{
		Tag: tag,
		Ttl: durationpb.New(response.DefaultTTL),
	}
	if ttl := rspHelper.GetMeta().GetTtl(); ttl != nil {
		responseMeta.Ttl = ttl
	}

	// Note: consider not overwriting rsp and whether it makes a difference.
	rsp := &fnv1beta1.RunFunctionResponse{
		Meta:    responseMeta,
		Desired: rspHelper.Desired,
		Results: rspHelper.Results,
		Context: rspHelper.Context,
	}

	if rspHelper.Requirements != nil && rspHelper.Requirements.ExtraResources != nil {
		rsp.Requirements = &fnv1beta1.Requirements{
			ExtraResources: convertExtraResources(rspHelper.Requirements.ExtraResources),
		}
	}
	return rsp, nil
}

func convertExtraResources(extraResources map[string]*helper.ResourceSelector) map[string]*fnv1beta1.ResourceSelector {
	out := make(map[string]*fnv1beta1.ResourceSelector)
	for name, fixedrs := range extraResources {
		rs := &fnv1beta1.ResourceSelector{
			ApiVersion: fixedrs.APIVersion,
			Kind:       fixedrs.Kind,
		}
		if fixedrs.Match.MatchLabels != nil && len(fixedrs.Match.MatchLabels.GetLabels()) > 0 {
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
		if pklSpec.URI == "" {
			return nil, errors.New("manifest type is uri but uri is empty")
		}
		return pkl.UriSource(pklSpec.URI), nil

	case "inline":
		if pklSpec.Inline == "" {
			return nil, errors.New("manifest type is inline but inline is empty")
		}
		return pkl.TextSource(pklSpec.Inline), nil
	case "local":
		if pklSpec.Local == nil {
			return nil, errors.New("manifest type is file but uri is empty")
		}
		return pkl.FileSource(pklSpec.Local.File), nil
	default:
		return nil, errors.New("unknown pklSpec type")
	}
}
