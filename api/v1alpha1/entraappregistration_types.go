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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EntraAppRegistrationSpec defines the desired state of EntraAppRegistration
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// EntraAppRegistration is the Schema for the entraappregistrations API
type EntraAppRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntraAppRegistrationSpec   `json:"spec,omitempty"`
	Status EntraAppRegistrationStatus `json:"status,omitempty"`
}

type EntraAppRegistrationSpec struct {
	// +kubebuilder:validation:Required
	ForProvider *AppRegCredConfig `json:"forProvider"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=120
	// +kubebuilder:validation:Pattern=`^[^<>%&:\\?\/\*]+$`
	Name string `json:"name,omitempty"`
}

type AppRegCredConfig struct {
	// +kubebuilder:validation:Optional
	ServiceAccountRef string `json:"serviceAccountRef,omitempty"`
	// +kubebuilder:validation:Optional
	CredentialSecretRef string `json:"credentialSecretRef,omitempty"`
}

// EntraAppRegistrationStatus defines the observed state of EntraAppRegistration
type EntraAppRegistrationStatus struct {
	// ObservedGeneration is the latest observed generation of the resource.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Phase represents the current phase of the EntraAppRegistration.
	Phase string `json:"phase,omitempty"`
	// AppRegistrationName is the name of the created App Registration in Entra ID
	AppRegistrationName string `json:"appRegistrationName,omitempty"`
	// AppRegistrationID is the ID of the created App Registration in Entra ID
	AppRegistrationID string `json:"appRegistrationID,omitempty"`
	// AppRegistrationObjID is the Object ID of the created App Registration in Entra ID
	AppRegistrationObjID string `json:"appRegistrationObjID,omitempty"`
}

// +kubebuilder:object:root=true

// EntraAppRegistrationList contains a list of EntraAppRegistration
type EntraAppRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EntraAppRegistration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EntraAppRegistration{}, &EntraAppRegistrationList{})
}
