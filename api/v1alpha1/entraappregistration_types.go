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
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".status.appRegistrationName"
// +kubebuilder:printcolumn:name="ClientID",type="string",JSONPath=".status.appRegistrationID"
// +kubebuilder:printcolumn:name="ObjectID",type="string",JSONPath=".status.appRegistrationObjID"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="LastRun",type="date",JSONPath=".status.lastRun"

// EntraAppRegistration is the Schema for the entraappregistrations API
type EntraAppRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntraAppRegistrationSpec   `json:"spec,omitempty"`
	Status EntraAppRegistrationStatus `json:"status,omitempty"`
}

type OptionalClaim struct {
	// +kubebuilder:validation:Optional
	AdditionalProperties []string `json:"additionalProperties,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Essential *bool `json:"essential,omitempty"`
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Optional
	Source *string `json:"source,omitempty"`
}

type OptionalClaims struct {
	// +kubebuilder:validation:Optional
	AccessToken []OptionalClaim `json:"accessToken,omitempty"`
	// +kubebuilder:validation:Optional
	IDToken []OptionalClaim `json:"idToken,omitempty"`
	// +kubebuilder:validation:Optional
	SAML2Token []OptionalClaim `json:"saml2Token,omitempty"`
}

// WebApplication specifies settings for a web application
type WebApplication struct {
	// +kubebuilder:validation:Optional
	HomePageURL *string `json:"homePageUrl,omitempty"`
	// +kubebuilder:validation:Optional
	ImplicitGrantSettings *ImplicitGrantSettings `json:"implicitGrantSettings,omitempty"`
	// +kubebuilder:validation:Optional
	LogoutURL *string `json:"logoutUrl,omitempty"`
	// +kubebuilder:validation:Optional
	RedirectURIs []string `json:"redirectUris,omitempty"`
}

// ImplicitGrantSettings specifies whether this web app can request tokens using OAuth 2.0 implicit flow
type ImplicitGrantSettings struct {
	// +kubebuilder:validation:Optional
	EnableAccessTokenIssuance *bool `json:"enableAccessTokenIssuance,omitempty"`
	// +kubebuilder:validation:Optional
	EnableIDTokenIssuance *bool `json:"enableIdTokenIssuance,omitempty"`
}

// SPAApplication specifies settings for a single-page application
type SPAApplication struct {
	// +kubebuilder:validation:Optional
	RedirectURIs []string `json:"redirectUris,omitempty"`
}

// PublicClientApplication specifies settings for installed clients (mobile/desktop)
type PublicClientApplication struct {
	// +kubebuilder:validation:Optional
	RedirectURIs []string `json:"redirectUris,omitempty"`
}

// PreAuthorizedApplication lists applications pre-authorized with specified permissions
type PreAuthorizedApplication struct {
	// +kubebuilder:validation:Required
	AppID string `json:"appId,omitempty"`
	// +kubebuilder:validation:Required
	PermissionIds []string `json:"delegatedPermissionIds,omitempty"`
}

// PermissionScope represents a delegated permission exposed by a web API
type PermissionScope struct {
	// +kubebuilder:validation:Optional
	AdminConsentDescription *string `json:"adminConsentDescription,omitempty"`
	// +kubebuilder:validation:Optional
	AdminConsentDisplayName *string `json:"adminConsentDisplayName,omitempty"`
	// +kubebuilder:validation:Required
	ID string `json:"id"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:default=true
	IsEnabled bool `json:"isEnabled"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=User;Admin
	Type string `json:"type,omitempty"`
	// +kubebuilder:validation:Optional
	UserConsentDescription *string `json:"userConsentDescription,omitempty"`
	// +kubebuilder:validation:Optional
	UserConsentDisplayName *string `json:"userConsentDisplayName,omitempty"`
	// +kubebuilder:validation:Optional
	Value *string `json:"value,omitempty"`
}

// APIApplication specifies settings for an application that implements a web API
type APIApplication struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	AcceptMappedClaims *bool `json:"acceptMappedClaims,omitempty"`
	// +kubebuilder:validation:Optional
	KnownClientApplications []string `json:"knownClientApplications,omitempty"`
	// +kubebuilder:validation:Optional
	OAuth2PermissionScopes []PermissionScope `json:"oauth2PermissionScopes,omitempty"`
	// +kubebuilder:validation:Optional
	PreAuthorizedApplications []PreAuthorizedApplication `json:"preAuthorizedApplications,omitempty"`
	// +kubebuilder:validation:Optional
	AppRoleAllowedMemberTypes []string `json:"appRoleAllowedMemberTypes,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=1;2
	RequestedAccessTokenVersion *int `json:"requestedAccessTokenVersion,omitempty"`
}

type InformationUrl struct {
	// +kubebuilder:validation:Optional
	LogoURL string `json:"logoUrl,omitempty"`
	// +kubebuilder:validation:Optional
	MarketingURL string `json:"marketingUrl,omitempty"`
	// +kubebuilder:validation:Optional
	SupportURL string `json:"supportUrl,omitempty"`
	// +kubebuilder:validation:Optional
	TermsOfServiceURL string `json:"termsOfServiceUrl,omitempty"`
	// +kubebuilder:validation:Optional
	PrivacyStatementURL string `json:"privacyStatementUrl,omitempty"`
}

type RequiredResourceAccess struct {
	// +kubebuilder:validation:Required
	ResourceAppID string `json:"resourceAppId,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems=1
	ResourceAccess []ResourceAccess `json:"resourceAccess,omitempty"`
}

type ResourceAccess struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=uuid
	ID string `json:"id,omitempty"`
	// +kubebuilder:validation:Required
	Type string `json:"type,omitempty"`
}

type EntraAppRegistrationSpec struct {
	// +kubebuilder:validation:Required
	ForProvider *AppRegCredConfig `json:"forProvider"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=120
	// +kubebuilder:validation:Pattern=`^[^<>%&:\\?\/\*]+$`
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags,omitempty"`
	// +kubebuilder:validation:Optional
	Notes string `json:"notes,omitempty"`
	// +kubebuilder:validation:Optional
	InformationUrl *InformationUrl `json:"info,omitempty"`
	// +kubebuilder:validation:Optional
	Owners *[]string `json:"owners,omitempty"`

	// sign-in configuration
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=AzureADMyOrg;AzureADMultipleOrgs;AzureADandPersonalMicrosoftAccount;PersonalMicrosoftAccount
	// +kubebuilder:default=AzureADMyOrg
	SignInAudience string `json:"signInAudience,omitempty"`
	// +kubebuilder:validation:Optional
	IdentifierUris []string `json:"identifierUris,omitempty"`

	// authentication settings
	// +kubebuilder:validation:Optional
	IsFallbackPublicClient bool `json:"isFallbackPublicClient,omitempty"`
	// +kubebuilder:validation:Optional
	IsDeviceOnlyAuthSupported bool `json:"isDeviceOnlyAuthSupported,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	OAuth2RequiredPostResponse bool `json:"oauth2RequirePostResponse,omitempty"`

	// saml configuration
	// +kubebuilder:validation:Optional
	SamlMetadataUrl string `json:"samlMetadataUrl,omitempty"`

	// token configuration
	// +kubebuilder:validation:Optional
	OptionalClaims *OptionalClaims `json:"optionalClaims,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=None;SecurityGroup;All
	// +kubebuilder:default=None
	GroupMembershipClaims string `json:"groupMembershipClaims,omitempty"`

	// Platform-specific settings
	// +kubebuilder:validation:Optional
	API *APIApplication `json:"api,omitempty"`
	// +kubebuilder:validation:Optional
	Web *WebApplication `json:"web,omitempty"`
	// +kubebuilder:validation:Optional
	SPA *SPAApplication `json:"spa,omitempty"`
	// +kubebuilder:validation:Optional
	PublicClient *PublicClientApplication `json:"publicClient,omitempty"`

	// +kubebuilder:validation:Optional
	ServicePrincipal *ServicePrincipalParams `json:"servicePrincipal,omitempty"`

	// Permissions and roles
	// +kubebuilder:validation:Optional
	RequiredResourceAccess []RequiredResourceAccess `json:"requiredResourceAccess,omitempty"`
}

type ServicePrincipalParams struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	AccountEnabled bool `json:"accountEnabled,omitempty"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	DisableVisibilityForGuests bool `json:"disableVisibilityForGuests,omitempty"`
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
	// ServicePrincipalID is the ID of the created Service Principal in Entra ID corresponding to the App Registration
	ServicePrincipalID string `json:"servicePrincipalID,omitempty"`
	// ServicePrincipalEnabled indicates whether the service principal was successfully created and is enabled
	ServicePrincipal string `json:"servicePrincipalEnabled,omitempty"`
	// Owners lists the Object IDs of the owners of the App Registration in Entra ID
	Owners []string `json:"owners,omitempty"`
	// LastRun indicates the last time the controller attempted to reconcile this resource
	LastRun metav1.Time `json:"lastRun,omitempty"`
	// Message provides additional information about the current status or any errors encountered
	Message string `json:"message,omitempty"`
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
