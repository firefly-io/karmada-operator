//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package search

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendStoreConfig) DeepCopyInto(out *BackendStoreConfig) {
	*out = *in
	if in.OpenSearch != nil {
		in, out := &in.OpenSearch, &out.OpenSearch
		*out = new(OpenSearchConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendStoreConfig.
func (in *BackendStoreConfig) DeepCopy() *BackendStoreConfig {
	if in == nil {
		return nil
	}
	out := new(BackendStoreConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenSearchConfig) DeepCopyInto(out *OpenSearchConfig) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.SecretRef = in.SecretRef
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenSearchConfig.
func (in *OpenSearchConfig) DeepCopy() *OpenSearchConfig {
	if in == nil {
		return nil
	}
	out := new(OpenSearchConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Proxying) DeepCopyInto(out *Proxying) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Proxying.
func (in *Proxying) DeepCopy() *Proxying {
	if in == nil {
		return nil
	}
	out := new(Proxying)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Proxying) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRegistry) DeepCopyInto(out *ResourceRegistry) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRegistry.
func (in *ResourceRegistry) DeepCopy() *ResourceRegistry {
	if in == nil {
		return nil
	}
	out := new(ResourceRegistry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceRegistry) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRegistryList) DeepCopyInto(out *ResourceRegistryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceRegistry, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRegistryList.
func (in *ResourceRegistryList) DeepCopy() *ResourceRegistryList {
	if in == nil {
		return nil
	}
	out := new(ResourceRegistryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceRegistryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRegistrySpec) DeepCopyInto(out *ResourceRegistrySpec) {
	*out = *in
	in.TargetCluster.DeepCopyInto(&out.TargetCluster)
	if in.ResourceSelectors != nil {
		in, out := &in.ResourceSelectors, &out.ResourceSelectors
		*out = make([]ResourceSelector, len(*in))
		copy(*out, *in)
	}
	if in.BackendStore != nil {
		in, out := &in.BackendStore, &out.BackendStore
		*out = new(BackendStoreConfig)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRegistrySpec.
func (in *ResourceRegistrySpec) DeepCopy() *ResourceRegistrySpec {
	if in == nil {
		return nil
	}
	out := new(ResourceRegistrySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRegistryStatus) DeepCopyInto(out *ResourceRegistryStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRegistryStatus.
func (in *ResourceRegistryStatus) DeepCopy() *ResourceRegistryStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceRegistryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSelector) DeepCopyInto(out *ResourceSelector) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSelector.
func (in *ResourceSelector) DeepCopy() *ResourceSelector {
	if in == nil {
		return nil
	}
	out := new(ResourceSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Search) DeepCopyInto(out *Search) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Search.
func (in *Search) DeepCopy() *Search {
	if in == nil {
		return nil
	}
	out := new(Search)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Search) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
