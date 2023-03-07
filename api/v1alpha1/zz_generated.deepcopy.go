//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 The Primaza Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationSelector) DeepCopyInto(out *ApplicationSelector) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationSelector.
func (in *ApplicationSelector) DeepCopy() *ApplicationSelector {
	if in == nil {
		return nil
	}
	out := new(ApplicationSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterEnvironment) DeepCopyInto(out *ClusterEnvironment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterEnvironment.
func (in *ClusterEnvironment) DeepCopy() *ClusterEnvironment {
	if in == nil {
		return nil
	}
	out := new(ClusterEnvironment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterEnvironment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterEnvironmentList) DeepCopyInto(out *ClusterEnvironmentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterEnvironment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterEnvironmentList.
func (in *ClusterEnvironmentList) DeepCopy() *ClusterEnvironmentList {
	if in == nil {
		return nil
	}
	out := new(ClusterEnvironmentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterEnvironmentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterEnvironmentSpec) DeepCopyInto(out *ClusterEnvironmentSpec) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ApplicationNamespaces != nil {
		in, out := &in.ApplicationNamespaces, &out.ApplicationNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ServiceNamespaces != nil {
		in, out := &in.ServiceNamespaces, &out.ServiceNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterEnvironmentSpec.
func (in *ClusterEnvironmentSpec) DeepCopy() *ClusterEnvironmentSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterEnvironmentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterEnvironmentStatus) DeepCopyInto(out *ClusterEnvironmentStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterEnvironmentStatus.
func (in *ClusterEnvironmentStatus) DeepCopy() *ClusterEnvironmentStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterEnvironmentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Environment) DeepCopyInto(out *Environment) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Environment.
func (in *Environment) DeepCopy() *Environment {
	if in == nil {
		return nil
	}
	out := new(Environment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentConstraints) DeepCopyInto(out *EnvironmentConstraints) {
	*out = *in
	if in.Environments != nil {
		in, out := &in.Environments, &out.Environments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentConstraints.
func (in *EnvironmentConstraints) DeepCopy() *EnvironmentConstraints {
	if in == nil {
		return nil
	}
	out := new(EnvironmentConstraints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HealthCheckContainer) DeepCopyInto(out *HealthCheckContainer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HealthCheckContainer.
func (in *HealthCheckContainer) DeepCopy() *HealthCheckContainer {
	if in == nil {
		return nil
	}
	out := new(HealthCheckContainer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegisteredService) DeepCopyInto(out *RegisteredService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegisteredService.
func (in *RegisteredService) DeepCopy() *RegisteredService {
	if in == nil {
		return nil
	}
	out := new(RegisteredService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RegisteredService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegisteredServiceConstraints) DeepCopyInto(out *RegisteredServiceConstraints) {
	*out = *in
	if in.Environments != nil {
		in, out := &in.Environments, &out.Environments
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegisteredServiceConstraints.
func (in *RegisteredServiceConstraints) DeepCopy() *RegisteredServiceConstraints {
	if in == nil {
		return nil
	}
	out := new(RegisteredServiceConstraints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegisteredServiceHealthCheck) DeepCopyInto(out *RegisteredServiceHealthCheck) {
	*out = *in
	out.Container = in.Container
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegisteredServiceHealthCheck.
func (in *RegisteredServiceHealthCheck) DeepCopy() *RegisteredServiceHealthCheck {
	if in == nil {
		return nil
	}
	out := new(RegisteredServiceHealthCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegisteredServiceList) DeepCopyInto(out *RegisteredServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RegisteredService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegisteredServiceList.
func (in *RegisteredServiceList) DeepCopy() *RegisteredServiceList {
	if in == nil {
		return nil
	}
	out := new(RegisteredServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RegisteredServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegisteredServiceSpec) DeepCopyInto(out *RegisteredServiceSpec) {
	*out = *in
	in.Constraints.DeepCopyInto(&out.Constraints)
	out.HealthCheck = in.HealthCheck
	if in.ServiceClassIdentity != nil {
		in, out := &in.ServiceClassIdentity, &out.ServiceClassIdentity
		*out = make([]ServiceClassIdentityItem, len(*in))
		copy(*out, *in)
	}
	if in.ServiceEndpointDefinition != nil {
		in, out := &in.ServiceEndpointDefinition, &out.ServiceEndpointDefinition
		*out = make([]ServiceEndpointDefinitionItem, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegisteredServiceSpec.
func (in *RegisteredServiceSpec) DeepCopy() *RegisteredServiceSpec {
	if in == nil {
		return nil
	}
	out := new(RegisteredServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegisteredServiceStatus) DeepCopyInto(out *RegisteredServiceStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegisteredServiceStatus.
func (in *RegisteredServiceStatus) DeepCopy() *RegisteredServiceStatus {
	if in == nil {
		return nil
	}
	out := new(RegisteredServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBinding) DeepCopyInto(out *ServiceBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBinding.
func (in *ServiceBinding) DeepCopy() *ServiceBinding {
	if in == nil {
		return nil
	}
	out := new(ServiceBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBindingList) DeepCopyInto(out *ServiceBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBindingList.
func (in *ServiceBindingList) DeepCopy() *ServiceBindingList {
	if in == nil {
		return nil
	}
	out := new(ServiceBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBindingSpec) DeepCopyInto(out *ServiceBindingSpec) {
	*out = *in
	in.Application.DeepCopyInto(&out.Application)
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]Environment, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBindingSpec.
func (in *ServiceBindingSpec) DeepCopy() *ServiceBindingSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceBindingStatus) DeepCopyInto(out *ServiceBindingStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceBindingStatus.
func (in *ServiceBindingStatus) DeepCopy() *ServiceBindingStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceBindingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceCatalog) DeepCopyInto(out *ServiceCatalog) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceCatalog.
func (in *ServiceCatalog) DeepCopy() *ServiceCatalog {
	if in == nil {
		return nil
	}
	out := new(ServiceCatalog)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceCatalog) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceCatalogList) DeepCopyInto(out *ServiceCatalogList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceCatalog, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceCatalogList.
func (in *ServiceCatalogList) DeepCopy() *ServiceCatalogList {
	if in == nil {
		return nil
	}
	out := new(ServiceCatalogList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceCatalogList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceCatalogService) DeepCopyInto(out *ServiceCatalogService) {
	*out = *in
	if in.ServiceClassIdentity != nil {
		in, out := &in.ServiceClassIdentity, &out.ServiceClassIdentity
		*out = make([]ServiceClassIdentityItem, len(*in))
		copy(*out, *in)
	}
	if in.ServiceEndpointDefinitionKeys != nil {
		in, out := &in.ServiceEndpointDefinitionKeys, &out.ServiceEndpointDefinitionKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceCatalogService.
func (in *ServiceCatalogService) DeepCopy() *ServiceCatalogService {
	if in == nil {
		return nil
	}
	out := new(ServiceCatalogService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceCatalogSpec) DeepCopyInto(out *ServiceCatalogSpec) {
	*out = *in
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = make([]ServiceCatalogService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceCatalogSpec.
func (in *ServiceCatalogSpec) DeepCopy() *ServiceCatalogSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceCatalogSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClaim) DeepCopyInto(out *ServiceClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClaim.
func (in *ServiceClaim) DeepCopy() *ServiceClaim {
	if in == nil {
		return nil
	}
	out := new(ServiceClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClaimApplicationClusterContext) DeepCopyInto(out *ServiceClaimApplicationClusterContext) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClaimApplicationClusterContext.
func (in *ServiceClaimApplicationClusterContext) DeepCopy() *ServiceClaimApplicationClusterContext {
	if in == nil {
		return nil
	}
	out := new(ServiceClaimApplicationClusterContext)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClaimList) DeepCopyInto(out *ServiceClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClaimList.
func (in *ServiceClaimList) DeepCopy() *ServiceClaimList {
	if in == nil {
		return nil
	}
	out := new(ServiceClaimList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClaimSpec) DeepCopyInto(out *ServiceClaimSpec) {
	*out = *in
	if in.ServiceClassIdentity != nil {
		in, out := &in.ServiceClassIdentity, &out.ServiceClassIdentity
		*out = make([]ServiceClassIdentityItem, len(*in))
		copy(*out, *in)
	}
	if in.ServiceEndpointDefinitionKeys != nil {
		in, out := &in.ServiceEndpointDefinitionKeys, &out.ServiceEndpointDefinitionKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Application.DeepCopyInto(&out.Application)
	out.ApplicationClusterContext = in.ApplicationClusterContext
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClaimSpec.
func (in *ServiceClaimSpec) DeepCopy() *ServiceClaimSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClaimStatus) DeepCopyInto(out *ServiceClaimStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClaimStatus.
func (in *ServiceClaimStatus) DeepCopy() *ServiceClaimStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceClaimStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClass) DeepCopyInto(out *ServiceClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClass.
func (in *ServiceClass) DeepCopy() *ServiceClass {
	if in == nil {
		return nil
	}
	out := new(ServiceClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClassIdentityItem) DeepCopyInto(out *ServiceClassIdentityItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClassIdentityItem.
func (in *ServiceClassIdentityItem) DeepCopy() *ServiceClassIdentityItem {
	if in == nil {
		return nil
	}
	out := new(ServiceClassIdentityItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClassList) DeepCopyInto(out *ServiceClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClassList.
func (in *ServiceClassList) DeepCopy() *ServiceClassList {
	if in == nil {
		return nil
	}
	out := new(ServiceClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClassMapping) DeepCopyInto(out *ServiceClassMapping) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClassMapping.
func (in *ServiceClassMapping) DeepCopy() *ServiceClassMapping {
	if in == nil {
		return nil
	}
	out := new(ServiceClassMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClassResource) DeepCopyInto(out *ServiceClassResource) {
	*out = *in
	if in.ServiceEndpointDefinitionMapping != nil {
		in, out := &in.ServiceEndpointDefinitionMapping, &out.ServiceEndpointDefinitionMapping
		*out = make([]ServiceClassMapping, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClassResource.
func (in *ServiceClassResource) DeepCopy() *ServiceClassResource {
	if in == nil {
		return nil
	}
	out := new(ServiceClassResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClassSpec) DeepCopyInto(out *ServiceClassSpec) {
	*out = *in
	if in.Constraints != nil {
		in, out := &in.Constraints, &out.Constraints
		*out = new(EnvironmentConstraints)
		(*in).DeepCopyInto(*out)
	}
	in.Resource.DeepCopyInto(&out.Resource)
	if in.ServiceClassIdentity != nil {
		in, out := &in.ServiceClassIdentity, &out.ServiceClassIdentity
		*out = make([]ServiceClassIdentityItem, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClassSpec.
func (in *ServiceClassSpec) DeepCopy() *ServiceClassSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceClassSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceClassStatus) DeepCopyInto(out *ServiceClassStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceClassStatus.
func (in *ServiceClassStatus) DeepCopy() *ServiceClassStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceClassStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceEndpointDefinitionItem) DeepCopyInto(out *ServiceEndpointDefinitionItem) {
	*out = *in
	out.ValueFromSecret = in.ValueFromSecret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceEndpointDefinitionItem.
func (in *ServiceEndpointDefinitionItem) DeepCopy() *ServiceEndpointDefinitionItem {
	if in == nil {
		return nil
	}
	out := new(ServiceEndpointDefinitionItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceEndpointDefinitionSecretRef) DeepCopyInto(out *ServiceEndpointDefinitionSecretRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceEndpointDefinitionSecretRef.
func (in *ServiceEndpointDefinitionSecretRef) DeepCopy() *ServiceEndpointDefinitionSecretRef {
	if in == nil {
		return nil
	}
	out := new(ServiceEndpointDefinitionSecretRef)
	in.DeepCopyInto(out)
	return out
}
