package function

import (
	"context"
	"testing"
	"time"

	"github.com/apple/pkl-go/pkl"
	"github.com/avarei/function-pkl/input/v1beta1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	//objectv1alpha2 "github.com/crossplane-contrib/provider-kubernetes/apis/object/v1alpha2"
)

func TestRunFunction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"SingleResource": {
			reason: "The function should return that it needs extra resources",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							PklCRDs: []v1beta1.PklCrdRef{
								{
									ApiVersion: "example.crossplane.io/v1",
									Kind:       "XR",
									Uri:        "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/crds/XR.pkl",
								},
								{
									ApiVersion: "kubernetes.crossplane.io/v1alpha2",
									Kind:       "Object",
									Uri:        "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/crds/Object.pkl",
								},
							},
							PklComposition: &v1beta1.PklFileRef{
								Name: "XR",
								Type: "uri",
								Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/crds/XR.pkl",
							},
							PklManifests: []v1beta1.PklFileRef{
								{
									Name: "object-one",
									Type: "uri",
									Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/object-one.pkl",
								},
								{
									Name: "object-two",
									Type: "uri",
									Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/object-two.pkl",
								},
								{
									Name: "object-three",
									Type: "uri",
									Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/object-three.pkl",
								},
							},
							Requirements: &v1beta1.PklFileRef{
								Name: "extra-resource",
								Type: "uri",
								Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/requirement.pkl",
							},
						},
					}),
					Observed: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "example.crossplane.io/v1",
								"kind": "XR",
								"metadata": {
									"name": "example-xr"
								},
								"spec": {}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"object-one": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-one"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-one"
												},
												"data": {
													"foo": "bar"
												}
											}
										}
									}
								}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{
						Tag: "extra",
						Ttl: durationpb.New(time.Second * 60),
					},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{},
						Resources: map[string]*fnv1beta1.Resource{
							"object-one": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-one"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-one",
													"namespace": "crossplane-system"
												},
												"data": {
													"foo": "example-xr"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_FALSE,
							},
							"object-two": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-two"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-two",
													"namespace": "crossplane-system"
												},
												"data": {
													"bar": "alternative"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_FALSE,
							},
							"object-three": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-three"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-three",
													"namespace": "crossplane-system"
												},
												"data": {
													"bar": "alternative"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_FALSE,
							},
						},
					},
					Requirements: &fnv1beta1.Requirements{
						ExtraResources: map[string]*fnv1beta1.ResourceSelector{
							"required": {
								ApiVersion: "kubernetes.crossplane.io/v1alpha2",
								Kind:       "Object",
								Match: &fnv1beta1.ResourceSelector_MatchName{
									MatchName: "required",
								},
							},
						},
					},
				},
			},
		},
		"RequestExtraResources": {
			reason: "The function should return that it needs extra resources",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							PklCRDs: []v1beta1.PklCrdRef{
								{
									ApiVersion: "example.crossplane.io/v1",
									Kind:       "XR",
									Uri:        "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/crds/XR.pkl",
								},
								{
									ApiVersion: "kubernetes.crossplane.io/v1alpha2",
									Kind:       "Object",
									Uri:        "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/crds/Object.pkl",
								},
							},
							PklComposition: &v1beta1.PklFileRef{
								Name: "XR",
								Type: "uri",
								Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/crds/XR.pkl",
							},
							PklManifests: []v1beta1.PklFileRef{
								{
									Name: "object-one",
									Type: "uri",
									Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/object-one.pkl",
								},
								{
									Name: "object-two",
									Type: "uri",
									Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/object-two.pkl",
								},
								{
									Name: "object-three",
									Type: "uri",
									Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/object-three.pkl",
								},
							},
							Requirements: &v1beta1.PklFileRef{
								Name: "extra-resource",
								Type: "uri",
								Uri:  "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.4#/requirement.pkl",
							},
						},
					}),
					Observed: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
								"apiVersion": "example.crossplane.io/v1",
								"kind": "XR",
								"metadata": {
									"name": "example-xr"
								},
								"spec": {}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{
							"object-one": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-one"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-one"
												},
												"data": {
													"foo": "bar"
												}
											}
										}
									}
								}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{
						Tag: "extra",
						Ttl: durationpb.New(time.Second * 60),
					},
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{},
						Resources: map[string]*fnv1beta1.Resource{
							"object-one": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-one"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-one",
													"namespace": "crossplane-system"
												},
												"data": {
													"foo": "example-xr"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_FALSE,
							},
							"object-two": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-two"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-two",
													"namespace": "crossplane-system"
												},
												"data": {
													"bar": "alternative"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_FALSE,
							},
							"object-three": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "kubernetes.crossplane.io/v1alpha2",
									"kind": "Object",
									"metadata": {
										"name": "cm-three"
									},
									"spec": {
										"forProvider": {
											"manifest": {
												"apiVersion": "v1",
												"kind": "ConfigMap",
												"metadata": {
													"name": "cm-three",
													"namespace": "crossplane-system"
												},
												"data": {
													"bar": "alternative"
												}
											}
										}
									}
								}`),
								Ready: fnv1beta1.Ready_READY_FALSE,
							},
						},
					},
					Requirements: &fnv1beta1.Requirements{
						ExtraResources: map[string]*fnv1beta1.ResourceSelector{
							"required": {
								ApiVersion: "kubernetes.crossplane.io/v1alpha2",
								Kind:       "Object",
								Match: &fnv1beta1.ResourceSelector_MatchName{
									MatchName: "required",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			evaluatorManager := pkl.NewEvaluatorManager()
			defer evaluatorManager.Close()
			f := &Function{Log: logging.NewNopLogger(), EvaluatorManager: evaluatorManager}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}
