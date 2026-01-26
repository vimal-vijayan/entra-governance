/*
Copyright 2026.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EntraSecurityGroupSpec defines the desired state of EntraSecurityGroup
type EntraSecurityGroupSpec struct {
	ForProvider *ProviderSpec `json:"forProvider,omitempty"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Pattern=`^[^<>%&:\\?\/\*]+$`
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// TODO: add a validation webhook for allowed values ("Unified", "DynamicMembership")
	// +kubebuilder:validation:Optional
	GroupTypes []string `json:"groupTypes,omitempty"`
	// +kubebuilder:validation:Optional
	MailNickname string `json:"mailNickname,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	MailEnabled bool `json:"mailEnabled,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	SecurityEnabled bool `json:"securityEnabled,omitempty"`
	// +kubebuilder:validation:Optional
	Owners []string `json:"owners,omitempty"`
	// +kubebuilder:validation:Optional
	Members []string `json:"members,omitempty"`
}

type ProviderSpec struct {
	CredentialSecretRef string `json:"credentialSecretRef,omitempty"`
	ServiceAccountRef   string `json:"serviceAccountRef,omitempty"`
}

// EntraSecurityGroupStatus defines the observed state of EntraSecurityGroup
type EntraSecurityGroupStatus struct {
	// ObservedGeneration is the latest observed generation of the resource.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions of the EntraSecurityGroup.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Phase represents the current phase of the EntraSecurityGroup.
	Phase string `json:"phase,omitempty"`
	// ID is the unique identifier of the EntraSecurityGroup in the external system.
	ID string `json:"id,omitempty"`
	// DisplayName is the display name of the EntraSecurityGroup.
	DisplayName string `json:"displayName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="The current phase of the EntraSecurityGroup"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of the EntraSecurityGroup"
// +kubebuilder:printcolumn:name="ID",type="string",JSONPath=".status.id",description="The ID of the EntraSecurityGroup in Entra"

// EntraSecurityGroup is the Schema for the entrasecuritygroups API
type EntraSecurityGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntraSecurityGroupSpec   `json:"spec,omitempty"`
	Status EntraSecurityGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// EntraSecurityGroupList contains a list of EntraSecurityGroup
type EntraSecurityGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntraSecurityGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EntraSecurityGroup{}, &EntraSecurityGroupList{})
}
