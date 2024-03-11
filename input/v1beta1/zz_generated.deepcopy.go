//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMapRef) DeepCopyInto(out *ConfigMapRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMapRef.
func (in *ConfigMapRef) DeepCopy() *ConfigMapRef {
	if in == nil {
		return nil
	}
	out := new(ConfigMapRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pkl) DeepCopyInto(out *Pkl) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pkl.
func (in *Pkl) DeepCopy() *Pkl {
	if in == nil {
		return nil
	}
	out := new(Pkl)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Pkl) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PklCrdRef) DeepCopyInto(out *PklCrdRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PklCrdRef.
func (in *PklCrdRef) DeepCopy() *PklCrdRef {
	if in == nil {
		return nil
	}
	out := new(PklCrdRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PklFileRef) DeepCopyInto(out *PklFileRef) {
	*out = *in
	out.ConfigMapRef = in.ConfigMapRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PklFileRef.
func (in *PklFileRef) DeepCopy() *PklFileRef {
	if in == nil {
		return nil
	}
	out := new(PklFileRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PklSpec) DeepCopyInto(out *PklSpec) {
	*out = *in
	if in.PklManifests != nil {
		in, out := &in.PklManifests, &out.PklManifests
		*out = make([]PklFileRef, len(*in))
		copy(*out, *in)
	}
	if in.PklCRDs != nil {
		in, out := &in.PklCRDs, &out.PklCRDs
		*out = make([]PklCrdRef, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PklSpec.
func (in *PklSpec) DeepCopy() *PklSpec {
	if in == nil {
		return nil
	}
	out := new(PklSpec)
	in.DeepCopyInto(out)
	return out
}
