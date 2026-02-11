package serviceprincipal

import (
	"context"
	"fmt"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CreateRequest struct {
	// Required
	AppID string

	// Basic properties
	DisplayName                string
	AccountEnabled             bool
	Description                string
	Tags                       []string
	DisableVisibilityForGuests bool

	// URLs and endpoints
	HomePage  string
	LoginUrl  string
	LogoutUrl string
	ReplyUrls []string

	// Names and identifiers
	AlternativeNames       []string
	ServicePrincipalNames  []string
	AppDescription         string
	AppOwnerOrganizationId string

	// Configuration
	AppRoleAssignmentRequired          bool
	SignInAudience                     string
	ServicePrincipalType               string
	PreferredSingleSignOnMode          string
	PreferredTokenSigningKeyThumbprint string
	Notes                              string
	NotificationEmailAddresses         []string
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*ServicePrincipalCreateResponse, error) {
	logger := log.FromContext(ctx)

	requestBody := graphmodels.NewServicePrincipal()

	// Required field
	requestBody.SetAppId(&req.AppID)

	// Basic properties
	requestBody.SetDisplayName(&req.DisplayName)
	requestBody.SetAccountEnabled(&req.AccountEnabled)
	requestBody.SetDescription(&req.Description)

	// Tags handling
	tags := req.Tags
	if req.DisableVisibilityForGuests {
		tags = append(tags, "HideApp")
	}
	requestBody.SetTags(tags)

	// URLs and endpoints
	requestBody.SetHomepage(&req.HomePage)
	requestBody.SetLoginUrl(&req.LoginUrl)
	requestBody.SetLogoutUrl(&req.LogoutUrl)
	requestBody.SetReplyUrls(req.ReplyUrls)

	// Names and identifiers
	requestBody.SetAlternativeNames(req.AlternativeNames)
	requestBody.SetServicePrincipalNames(req.ServicePrincipalNames)
	requestBody.SetAppDescription(&req.AppDescription)

	// Configuration
	requestBody.SetAppRoleAssignmentRequired(&req.AppRoleAssignmentRequired)
	requestBody.SetSignInAudience(&req.SignInAudience)
	requestBody.SetServicePrincipalType(&req.ServicePrincipalType)
	requestBody.SetPreferredSingleSignOnMode(&req.PreferredSingleSignOnMode)
	requestBody.SetPreferredTokenSigningKeyThumbprint(&req.PreferredTokenSigningKeyThumbprint)
	requestBody.SetNotes(&req.Notes)
	requestBody.SetNotificationEmailAddresses(req.NotificationEmailAddresses)

	response, err := s.sdk.ServicePrincipals().Post(ctx, requestBody, nil)

	if err != nil {
		logger.Error(err, "failed to create service principal", "applicationID", req.AppID)
		return nil, err
	}

	// Check if response or ID is nil before dereferencing
	if response == nil {
		logger.Error(nil, "service principal response is nil", "applicationID", req.AppID)
		return nil, fmt.Errorf("service principal response is nil")
	}

	servicePrincipalID := response.GetId()
	if servicePrincipalID == nil {
		logger.Error(nil, "service principal ID is nil", "applicationID", req.AppID)
		return nil, fmt.Errorf("service principal ID is nil in response")
	}

	logger.Info("service principal created successfully", "applicationID", req.AppID, "servicePrincipalID", *servicePrincipalID)

	return &ServicePrincipalCreateResponse{
		ServicePrincipalID: *servicePrincipalID,
	}, nil
}
