/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=clusterclasses,shortName=cc,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of ClusterClass"

// ClusterClass is a template which can be used to create managed topologies.
type ClusterClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ClusterClassSpec `json:"spec,omitempty"`
}

// ClusterClassSpec describes the desired state of the ClusterClass.
type ClusterClassSpec struct {
	// Infrastructure is a reference to a provider-specific template that holds
	// the details for provisioning infrastructure specific cluster
	// for the underlying provider.
	// The underlying provider is responsible for the implementation
	// of the template to an infrastructure cluster.
	// +optional
	Infrastructure LocalObjectTemplate `json:"infrastructure,omitempty"`

	// ControlPlane is a reference to a local struct that holds the details
	// for provisioning the Control Plane for the Cluster.
	// +optional
	ControlPlane ControlPlaneClass `json:"controlPlane,omitempty"`

	// Workers describes the worker nodes for the cluster.
	// It is a collection of node types which can be used to create
	// the worker nodes of the cluster.
	// +optional
	Workers WorkersClass `json:"workers,omitempty"`

	// Variables defines the variables which can be configured
	// in the Cluster topology and are then used in patches.
	// +optional
	Variables []ClusterClassVariable `json:"variables,omitempty"`

	// Patches defines the patches which are applied to customize
	// referenced templates of a ClusterClass.
	// Note: Patches will be applied in the order of the array.
	// +optional
	Patches []ClusterClassPatch `json:"patches,omitempty"`
}

// ControlPlaneClass defines the class for the control plane.
type ControlPlaneClass struct {
	// Metadata is the metadata applied to the machines of the ControlPlane.
	// At runtime this metadata is merged with the corresponding metadata from the topology.
	//
	// This field is supported if and only if the control plane provider template
	// referenced is Machine based.
	// +optional
	Metadata ObjectMeta `json:"metadata,omitempty"`

	// LocalObjectTemplate contains the reference to the control plane provider.
	LocalObjectTemplate `json:",inline"`

	// MachineTemplate defines the metadata and infrastructure information
	// for control plane machines.
	//
	// This field is supported if and only if the control plane provider template
	// referenced above is Machine based and supports setting replicas.
	//
	// +optional
	MachineInfrastructure *LocalObjectTemplate `json:"machineInfrastructure,omitempty"`
}

// WorkersClass is a collection of deployment classes.
type WorkersClass struct {
	// MachineDeployments is a list of machine deployment classes that can be used to create
	// a set of worker nodes.
	// +optional
	MachineDeployments []MachineDeploymentClass `json:"machineDeployments,omitempty"`
}

// MachineDeploymentClass serves as a template to define a set of worker nodes of the cluster
// provisioned using the `ClusterClass`.
type MachineDeploymentClass struct {
	// Class denotes a type of worker node present in the cluster,
	// this name MUST be unique within a ClusterClass and can be referenced
	// in the Cluster to create a managed MachineDeployment.
	Class string `json:"class"`

	// Template is a local struct containing a collection of templates for creation of
	// MachineDeployment objects representing a set of worker nodes.
	Template MachineDeploymentClassTemplate `json:"template"`
}

// MachineDeploymentClassTemplate defines how a MachineDeployment generated from a MachineDeploymentClass
// should look like.
type MachineDeploymentClassTemplate struct {
	// Metadata is the metadata applied to the machines of the MachineDeployment.
	// At runtime this metadata is merged with the corresponding metadata from the topology.
	// +optional
	Metadata ObjectMeta `json:"metadata,omitempty"`

	// Bootstrap contains the bootstrap template reference to be used
	// for the creation of worker Machines.
	Bootstrap LocalObjectTemplate `json:"bootstrap"`

	// Infrastructure contains the infrastructure template reference to be used
	// for the creation of worker Machines.
	Infrastructure LocalObjectTemplate `json:"infrastructure"`
}

// ClusterClassVariable defines a variable which can
// be configured in the Cluster topology and used in patches.
type ClusterClassVariable struct {
	// Name of the variable.
	Name string `json:"name"`

	// Required specifies if the variable is required.
	// Note: this applies to the variable as a whole and thus the
	// top-level object defined in the schema. If nested fields are
	// required, this will be specified inside the schema.
	Required bool `json:"required"`

	// Schema defines the schema of the variable.
	Schema VariableSchema `json:"schema"`
}

// VariableSchema defines the schema of a variable.
type VariableSchema struct {
	// OpenAPIV3Schema defines the schema of a variable via OpenAPI v3
	// schema. The schema is a subset of the schema used in
	// Kubernetes CRDs.
	OpenAPIV3Schema JSONSchemaProps `json:"openAPIV3Schema"`
}

// JSONSchemaProps is a JSON-Schema following Specification Draft 4 (http://json-schema.org/).
// This struct has been initially copied from apiextensionsv1.JSONSchemaProps, but all fields
// which are not supported in CAPI have been removed.
type JSONSchemaProps struct {
	// Type is the type of the variable.
	// Valid values are: string, integer, number or boolean.
	Type string `json:"type"`

	// Nullable specifies if the variable can be set to null.
	// +optional
	Nullable bool `json:"nullable,omitempty"`

	// Format is an OpenAPI v3 format string. Unknown formats are ignored.
	// For a list of supported formats please see: (of the k8s.io/apiextensions-apiserver version we're currently using)
	// https://github.com/kubernetes/apiextensions-apiserver/blob/master/pkg/apiserver/validation/formats.go
	// +optional
	Format string `json:"format,omitempty"`

	// MaxLength is the max length of a string variable.
	// +optional
	MaxLength *int64 `json:"maxLength,omitempty"`

	// MinLength is the min length of a string variable.
	// +optional
	MinLength *int64 `json:"minLength,omitempty"`

	// Pattern is the regex which a string variable must match.
	// +optional
	Pattern string `json:"pattern,omitempty"`

	// Maximum is the maximum of an integer or number variable.
	// If ExclusiveMaximum is false, the variable is valid if it is lower than, or equal to, the value of Maximum.
	// If ExclusiveMaximum is true, the variable is valid if it is strictly lower than the value of Maximum.
	// +optional
	Maximum *json.Number `json:"maximum,omitempty"`

	// ExclusiveMaximum specifies if the Maximum is exclusive.
	// +optional
	ExclusiveMaximum bool `json:"exclusiveMaximum,omitempty"`

	// Minimum is the minimum of an integer or number variable.
	// If ExclusiveMinimum is false, the variable is valid if it is greater than, or equal to, the value of Minimum.
	// If ExclusiveMinimum is true, the variable is valid if it is strictly greater than the value of Minimum.
	// +optional
	Minimum *json.Number `json:"minimum,omitempty"`

	// ExclusiveMinimum specifies if the Minimum is exclusive.
	// +optional
	ExclusiveMinimum bool `json:"exclusiveMinimum,omitempty"`

	// MultipleOf specifies the number of which the variable must be a multiple of.
	// +optional
	MultipleOf *json.Number `json:"multipleOf,omitempty"`

	// Enum is the list of valid values of the variable.
	// +optional
	Enum []apiextensionsv1.JSON `json:"enum,omitempty"`

	// Default is the default value of the variable.
	// +optional
	Default *apiextensionsv1.JSON `json:"default,omitempty"`
}

// ClusterClassPatch defines a patch which is applied to customize the referenced templates.
type ClusterClassPatch struct {
	// Name of the patch.
	Name string `json:"name"`

	// Definitions define the patches inline.
	// Note: Patches will be applied in the order of the array.
	Definitions []PatchDefinition `json:"definitions"`
}

// PatchDefinition defines a patch which is applied to customize the referenced templates.
type PatchDefinition struct {
	// Selector defines on which templates the patch should be applied.
	Selector PatchSelector `json:"selector"`

	// JSONPatches defines the patches which should be applied on the templates
	// matching the selector.
	// Note: Patches will be applied in the order of the array.
	JSONPatches []JSONPatch `json:"jsonPatches"`
}

// PatchSelector defines on which templates the patch should be applied.
// Note: Matching on APIVersion and Kind is mandatory, to enforce that the patches are
// written for the correct version. The version of the references in the ClusterClass may
// be automatically updated during reconciliation if there is a newer version for the same contract.
// Note: The results of selection based on the individual fields are ANDed.
type PatchSelector struct {
	// APIVersion filters templates by apiVersion.
	APIVersion string `json:"apiVersion"`

	// Kind filters templates by kind.
	Kind string `json:"kind"`

	// MatchResources selects templates based on where they are referenced.
	MatchResources PatchSelectorMatch `json:"matchResources"`
}

// PatchSelectorMatch selects templates based on where they are referenced.
// Note: At least one of the fields must be set.
// Note: The results of selection based on the individual fields are ORed.
type PatchSelectorMatch struct {
	// ControlPlane selects templates referenced in .spec.ControlPlane.
	// Note: this will match the controlPlane and also the controlPlane
	// machineInfrastructure (depending on the kind and apiVersion).
	// +optional
	ControlPlane *bool `json:"controlPlane,omitempty"`

	// InfrastructureCluster selects templates referenced in .spec.infrastructure.
	// +optional
	InfrastructureCluster *bool `json:"infrastructureCluster,omitempty"`

	// MachineDeploymentClass selects templates referenced in specific MachineDeploymentClasses in
	// .spec.workers.machineDeployments.
	// +optional
	MachineDeploymentClass *PatchSelectorMatchMachineDeploymentClass `json:"machineDeploymentClass,omitempty"`
}

// PatchSelectorMatchMachineDeploymentClass selects templates referenced
// in specific MachineDeploymentClasses in .spec.workers.machineDeployments.
type PatchSelectorMatchMachineDeploymentClass struct {
	// Names selects templates by class names.
	Names []string `json:"names"`
}

// JSONPatch defines a JSON patch.
type JSONPatch struct {
	// Op defines the operation of the patch.
	// Note: Only `add`, `replace` and `remove` are supported.
	Op string `json:"op"`

	// Path defines the path of the patch.
	// Note: Only the spec of a template can be patched, thus the path has to start with /spec/.
	// Note: For now the only allowed array modifications are `append` and `prepend`, i.e.:
	// * for op: `add`: only index 0 (prepend) and - (append) are allowed
	// * for op: `replace` or `remove`: no indexes are allowed
	Path string `json:"path"`

	// Value defines the value of the patch.
	// Note: Either Value or ValueFrom is required for add and replace
	// operations. Only one of them is allowed to be set at the same time.
	// Note: We have to use apiextensionsv1.JSON instead of our JSON type,
	// because controller-tools has a hard-coded schema for apiextensionsv1.JSON
	// which cannot be produced by another type (unset type field).
	// Ref: https://github.com/kubernetes-sigs/controller-tools/blob/d0e03a142d0ecdd5491593e941ee1d6b5d91dba6/pkg/crd/known_types.go#L106-L111
	// +optional
	Value *apiextensionsv1.JSON `json:"value,omitempty"`

	// ValueFrom defines the value of the patch.
	// Note: Either Value or ValueFrom is required for add and replace
	// operations. Only one of them is allowed to be set at the same time.
	// +optional
	ValueFrom *JSONPatchValue `json:"valueFrom,omitempty"`
}

// JSONPatchValue defines the value of a patch.
// Note: Only one of the fields is allowed to be set at the same time.
type JSONPatchValue struct {
	// Variable is the variable to be used as value.
	// Variable can be one of the variables defined in .spec.variables or a builtin variable.
	// +optional
	Variable *string `json:"variable,omitempty"`

	// Template is the Go template to be used to calculate the value.
	// A template can reference variables defined in .spec.variables and builtin variables.
	// Note: The template must evaluate to a valid YAML or JSON value.
	// +optional
	Template *string `json:"template,omitempty"`
}

// LocalObjectTemplate defines a template for a topology Class.
type LocalObjectTemplate struct {
	// Ref is a required reference to a custom resource
	// offered by a provider.
	Ref *corev1.ObjectReference `json:"ref"`
}

// +kubebuilder:object:root=true

// ClusterClassList contains a list of Cluster.
type ClusterClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterClass{}, &ClusterClassList{})
}