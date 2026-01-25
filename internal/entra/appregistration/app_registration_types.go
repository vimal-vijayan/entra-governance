package appregistration

import (
	"time"

	uuid "github.com/google/uuid"
)

type AppRegistrationRequest struct {
	DisplayName            string
	description            string
	Tags                   []string
	identifierUris         []string
	Api                    *ApiPermissionRequest
	PasswordCredentials    []PasswordCredentialRequest
	ApplicationTemplatedID string
	AppRoles               []AppRoleRequest
	SignInAudience         string
	Web                    *WebApplicationRequest
	UniqueName             string
	RequiredResourceAccess []RequiredResourceAccessRequest
	PublisherDomain        string
}

type RequiredResourceAccessRequest struct {
	ReourceAccess []ResourceAccessRequest
	ResourceAppId string
}

type ResourceAccessRequest struct {
	Id   uuid.UUID
	Type string
}

type WebApplicationRequest struct {
	HomePageUrl           string
	LogoutUrl             string
	RedirectUris          []string
	ImplicitGrantSettings *ImplicitGrantSettingsRequest
}
type ImplicitGrantSettingsRequest struct {
	enableAccessTokenIssuance bool
	enableIdTokenIssuance     bool
}

type PasswordCredentialRequest struct {
	DisplayName         string
	EndDateTime         time.Time
	CustomKeyIdentifier []byte
	Hint                string
	KeyId               uuid.UUID
	SecretText          string
	StartDateTime       time.Time
}

type AppRoleRequest struct{}

type ApiPermissionRequest struct{}
