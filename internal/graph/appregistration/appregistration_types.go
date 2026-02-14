package appregistration

type Application struct {
	ObjectID                  string
	AppID                     string
	DisplayName               string
	Description               string
	Tags                      []string
	SignInAudience            string
	SamlMetadataURL           string
	IsFallbackPublicClient    bool
	IsDeviceOnlyAuthSupported bool
	GroupMembershipClaims     string
}

// APIApplication specifies settings for an application that implements a web API
type APIApplication struct {
	AcceptMappedClaims          *bool                      `json:"acceptMappedClaims,omitempty"`
	KnownClientApplications     []string                   `json:"knownClientApplications,omitempty"`
	OAuth2PermissionScopes      []PermissionScope          `json:"oauth2PermissionScopes,omitempty"`
	PreAuthorizedApplications   []PreAuthorizedApplication `json:"preAuthorizedApplications,omitempty"`
	RequestedAccessTokenVersion *int                       `json:"requestedAccessTokenVersion,omitempty"`
}

// PermissionScope represents a delegated permission exposed by a web API
type PermissionScope struct {
	AdminConsentDescription *string `json:"adminConsentDescription,omitempty"`
	AdminConsentDisplayName *string `json:"adminConsentDisplayName,omitempty"`
	ID                      string  `json:"id"`
	IsEnabled               bool    `json:"isEnabled"`
	Type                    string  `json:"type,omitempty"`
	UserConsentDescription  *string `json:"userConsentDescription,omitempty"`
	UserConsentDisplayName  *string `json:"userConsentDisplayName,omitempty"`
	Value                   *string `json:"value,omitempty"`
}

// PreAuthorizedApplication lists applications pre-authorized with specified permissions
type PreAuthorizedApplication struct {
	AppID         string   `json:"appId,omitempty"`
	PermissionIds []string `json:"delegatedPermissionIds,omitempty"`
}

// WebApplication specifies settings for a web application
type WebApplication struct {
	HomePageURL           *string                `json:"homePageUrl,omitempty"`
	ImplicitGrantSettings *ImplicitGrantSettings `json:"implicitGrantSettings,omitempty"`
	LogoutURL             *string                `json:"logoutUrl,omitempty"`
	RedirectURIs          []string               `json:"redirectUris,omitempty"`
}

// ImplicitGrantSettings specifies whether this web app can request tokens using OAuth 2.0 implicit flow
type ImplicitGrantSettings struct {
	EnableAccessTokenIssuance *bool `json:"enableAccessTokenIssuance,omitempty"`
	EnableIDTokenIssuance     *bool `json:"enableIdTokenIssuance,omitempty"`
}

// SPAApplication specifies settings for a single-page application
type SPAApplication struct {
	RedirectURIs []string `json:"redirectUris,omitempty"`
}

// PublicClientApplication specifies settings for installed clients (mobile/desktop)
type PublicClientApplication struct {
	RedirectURIs []string `json:"redirectUris,omitempty"`
}

// InformationalURL contains basic profile information URLs
type InformationalURL struct {
	LogoURL             *string `json:"logoUrl,omitempty"`
	MarketingURL        *string `json:"marketingUrl,omitempty"`
	PrivacyStatementURL *string `json:"privacyStatementUrl,omitempty"`
	SupportURL          *string `json:"supportUrl,omitempty"`
	TermsOfServiceURL   *string `json:"termsOfServiceUrl,omitempty"`
}

// RequiredResourceAccess specifies OAuth 2.0 permission scopes and app roles the app requires
type RequiredResourceAccess struct {
	ResourceAppID  string           `json:"resourceAppId"`
	ResourceAccess []ResourceAccess `json:"resourceAccess"`
}

// ResourceAccess identifies an OAuth 2.0 permission scope or app role
type ResourceAccess struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// AppRole represents an application role
type AppRole struct {
	AllowedMemberTypes []string `json:"allowedMemberTypes"`
	Description        string   `json:"description"`
	DisplayName        string   `json:"displayName"`
	ID                 string   `json:"id"`
	IsEnabled          bool     `json:"isEnabled"`
	Value              *string  `json:"value,omitempty"`
}

// KeyCredential represents a key credential associated with an application
type KeyCredential struct {
	CustomKeyIdentifier *string `json:"customKeyIdentifier,omitempty"`
	DisplayName         *string `json:"displayName,omitempty"`
	EndDateTime         *string `json:"endDateTime,omitempty"`
	Key                 *string `json:"key,omitempty"`
	KeyID               *string `json:"keyId,omitempty"`
	StartDateTime       *string `json:"startDateTime,omitempty"`
	Type                *string `json:"type,omitempty"`
	Usage               *string `json:"usage,omitempty"`
}

// PasswordCredential represents a password credential associated with an application
type PasswordCredential struct {
	CustomKeyIdentifier *string `json:"customKeyIdentifier,omitempty"`
	DisplayName         *string `json:"displayName,omitempty"`
	EndDateTime         *string `json:"endDateTime,omitempty"`
	Hint                *string `json:"hint,omitempty"`
	KeyID               *string `json:"keyId,omitempty"`
	SecretText          *string `json:"secretText,omitempty"`
	StartDateTime       *string `json:"startDateTime,omitempty"`
}

// OptionalClaims specifies optional claims requested by the application
type OptionalClaims struct {
	AccessToken []OptionalClaim `json:"accessToken,omitempty"`
	IDToken     []OptionalClaim `json:"idToken,omitempty"`
	SAML2Token  []OptionalClaim `json:"saml2Token,omitempty"`
}

// OptionalClaim represents a single optional claim
type OptionalClaim struct {
	AdditionalProperties []string `json:"additionalProperties,omitempty"`
	Essential            *bool    `json:"essential,omitempty"`
	Name                 string   `json:"name"`
	Source               *string  `json:"source,omitempty"`
}

// ParentalControlSettings specifies parental control settings for an application
type ParentalControlSettings struct {
	CountriesBlockedForMinors []string `json:"countriesBlockedForMinors,omitempty"`
	LegalAgeGroupRule         *string  `json:"legalAgeGroupRule,omitempty"`
}

// CreateRequest represents the request body for creating a new application
type CreateRequest struct {
	// Required field
	DisplayName string `json:"displayName"`

	// Basic information
	Description *string  `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Notes       *string  `json:"notes,omitempty"`
	Owners      []string `json:"owners,omitempty"`

	// Sign-in configuration
	SignInAudience *string  `json:"signInAudience,omitempty"`
	IdentifierURIs []string `json:"identifierUris,omitempty"`

	// Authentication settings
	IsFallbackPublicClient     *bool `json:"isFallbackPublicClient,omitempty"`
	IsDeviceOnlyAuthSupported  *bool `json:"isDeviceOnlyAuthSupported,omitempty"`
	OAuth2RequiredPostResponse *bool `json:"oauth2RequirePostResponse,omitempty"`

	// Token configuration
	GroupMembershipClaims *string         `json:"groupMembershipClaims,omitempty"`
	TokenEncryptionKeyID  *string         `json:"tokenEncryptionKeyId,omitempty"`
	OptionalClaims        *OptionalClaims `json:"optionalClaims,omitempty"`

	// SAML configuration
	SamlMetadataURL *string `json:"samlMetadataUrl,omitempty"`

	// Platform-specific settings
	API          *APIApplication          `json:"api,omitempty"`
	Web          *WebApplication          `json:"web,omitempty"`
	SPA          *SPAApplication          `json:"spa,omitempty"`
	PublicClient *PublicClientApplication `json:"publicClient,omitempty"`

	// Permissions and roles
	RequiredResourceAccess []RequiredResourceAccess `json:"requiredResourceAccess,omitempty"`
	AppRoles               []AppRole                `json:"appRoles,omitempty"`

	// Credentials
	KeyCredentials      []KeyCredential      `json:"keyCredentials,omitempty"`
	PasswordCredentials []PasswordCredential `json:"passwordCredentials,omitempty"`

	// Information URLs
	Info *InformationalURL `json:"info,omitempty"`

	// Additional settings
	ParentalControlSettings    *ParentalControlSettings `json:"parentalControlSettings,omitempty"`
	ServiceManagementReference *string                  `json:"serviceManagementReference,omitempty"`
}

type CreateResponse struct {
	AppClientID string
	AppObjectID string
}

type PatchRequest struct {
	ObjectID                  string
	DisplayName               *string
	Description               *string
	Tags                      *[]string
	SignInAudience            *string
	SamlMetadataURL           *string
	IsFallbackPublicClient    *bool
	IsDeviceOnlyAuthSupported *bool
	GroupMembershipClaims     *string
}
