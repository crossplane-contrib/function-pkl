package function

import (
	"context"
	"fmt"
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

var (
	pklPackage     = "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.1.2"
	pklCorePackage = "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.19"
	pklK8sPackage  = "package://pkg.pkl-lang.org/pkl-k8s/k8s@1.0.1"
)

func DefaultCRDs() []v1beta1.PklCrdRef {
	return []v1beta1.PklCrdRef{
		{
			ApiVersion: "example.crossplane.io/v1",
			Kind:       "XR",
			Uri:        pklPackage + "#/crds/XR.pkl",
		},
		{
			ApiVersion: "kubernetes.crossplane.io/v1alpha2",
			Kind:       "Object",
			Uri:        pklPackage + "#/crds/Object.pkl",
		},
	}
}

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
		"SingleResource-Uri": {
			reason: "The Function should parse one pkl file",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							CRDs: DefaultCRDs(),
							Resources: []v1beta1.PklFileRef{
								{
									Name: "object-one",
									Type: "uri",
									Uri:  pklPackage + "#/object-one.pkl",
								},
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
						},
					},
				},
			},
		},
		"Composite-Status": {
			reason: "The Function should update the status of the Composite Resource",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							CRDs: DefaultCRDs(),
							Composition: &v1beta1.PklFileRef{
								Name: "xr",
								Type: "uri",
								Uri:  pklPackage + "#/xr.pkl",
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
							"cm-one": {
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
									},
									"status": {
										"atProvider": {
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
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
							"apiVersion": "example.crossplane.io/v1",
							"kind": "XR",
							"metadata": {
								"name": "example-xr"
							},
							"spec": {},
							"status": {
								"someStatus": "I observed cm-one's namespace. it is crossplane-system"
							}
							}`),
						},
						Resources: map[string]*fnv1beta1.Resource{},
					},
				},
			},
		},

		"ExtraResource": {
			reason: "The Function should parse a requested and given ExtraResource",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							CRDs: DefaultCRDs(),
							Composition: &v1beta1.PklFileRef{
								Name: "XR",
								Type: "uri",
								Uri:  pklPackage + "#/crds/XR.pkl",
							},
							Resources: []v1beta1.PklFileRef{
								{
									Name: "object-needs-extra-resource",
									Type: "uri",
									Uri:  pklPackage + "#/object-needs-extra-resource.pkl",
								},
							},
							Requirements: &v1beta1.PklFileRef{
								Name: "extra-resource",
								Type: "uri",
								Uri:  pklPackage + "#/extra-resource.pkl",
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
					},
					ExtraResources: map[string]*fnv1beta1.Resources{
						"myextras": {
							Items: []*fnv1beta1.Resource{
								{
									Resource: resource.MustStructJSON(`{
		   												"apiVersion": "kubernetes.crossplane.io/v1alpha2",
		   												"kind": "Object",
		   												"metadata": {
		   													"name": "iamspecial"
		   												}
		   											}`),
								},
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
							"object-needs-extra-resource": {
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
		   															"bar": "iamspecial"
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
							"myextras": {
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

		"Packages": {
			reason: "The function should correctly replace Package references in it's input",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							CRDs: DefaultCRDs(),
							Composition: &v1beta1.PklFileRef{
								Name: "XR",
								Type: "uri",
								Uri:  "@example/crds/XR.pkl",
							},
							Resources: []v1beta1.PklFileRef{
								{
									Name: "object-one",
									Type: "uri",
									Uri:  "@example/object-one.pkl",
								},
							},
							Requirements: &v1beta1.PklFileRef{
								Name: "extra-resource",
								Type: "uri",
								Uri:  "@example/extra-resource.pkl",
							},
							Packages: []v1beta1.Package{
								{
									Name: "example",
									Core: false,
									Uri:  pklPackage,
								},
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
						},
					},
					Requirements: &fnv1beta1.Requirements{
						ExtraResources: map[string]*fnv1beta1.ResourceSelector{
							"myextras": {
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

		"SingleResource-Inline": {
			reason: "The Function should parse one inline pkl file",
			args: args{
				ctx: context.TODO(),
				req: &fnv1beta1.RunFunctionRequest{
					Meta: &fnv1beta1.RequestMeta{Tag: "extra"},
					Input: resource.MustStructObject(&v1beta1.Pkl{
						Spec: v1beta1.PklSpec{
							CRDs: DefaultCRDs(),
							Composition: &v1beta1.PklFileRef{
								Name: "XR",
								Type: "inline",
								Inline: fmt.Sprintf(`
		   amends "%[1]s#/CrossplaneResource.pkl"
		   import "%[1]s#/CompositionInput.pkl"

		   import "%[2]s#/crds/XR.pkl"
		   import "%[3]s#/api/core/v1/ConfigMap.pkl"

		   local state = import("crossplane:state") as CompositionInput
		   local observedCompositeResource: XR = state.observed.composite.resource as XR
		   local cmOne: ConfigMap? = state.observed.resources.getOrNull("cm-one")?.resource as ConfigMap?

		   resource = (observedCompositeResource) {
		     status {
		       when (cmOne?.metadata?.namespace != null) {
		         someStatus = "I observed cm-one's namespace. it is \(cmOne.metadata.namespace)"
		       }
		     }
		   }
		   connectionDetails {
		     ["test"] = "bar"
		   }
		   `, pklCorePackage, pklPackage, pklK8sPackage),
							},
							Resources: []v1beta1.PklFileRef{
								{
									Name: "object-one",
									Type: "inline",
									Inline: fmt.Sprintf(`
		   amends "%[1]s#/CrossplaneResource.pkl"
		   import "%[1]s#/CompositionInput.pkl"

		   import "%[2]s#/crds/XR.pkl"
		   import "%[2]s#/crds/Object.pkl"
		   import "%[3]s#/api/core/v1/ConfigMap.pkl"

		   local state = import("crossplane:state") as CompositionInput
		   local observedCompositeResource: XR = state.observed.composite.resource as XR

		   resource = (Object) {
		     metadata {
		       name = "cm-one"
		     }

		     spec {
		       forProvider {
		         manifest = (ConfigMap) {
		           metadata {
		             name = "cm-one"
		             namespace = "crossplane-system"
		           }
		           data {
		             ["foo"] = observedCompositeResource.metadata.name ?? throw("Composite could not find observed composite name")
		           }
		         }
		       }
		     }
		   }
		   ready = Ready_READY_FALSE
		   `, pklCorePackage, pklPackage, pklK8sPackage),
								},
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
						Composite: &fnv1beta1.Resource{
							Resource: resource.MustStructJSON(`{
		   							"apiVersion": "example.crossplane.io/v1",
		   							"kind": "XR",
		   							"metadata": {
		   								"name": "example-xr"
		   							},
		   							"spec": {},
		   							"status": {}
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
