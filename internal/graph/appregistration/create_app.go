package appregistration

import (
	"context"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Create(ctx context.Context, req CreateRequest) (*CreateResponse, error) {
	logger := log.FromContext(ctx)

	body := graphmodels.NewApplication()
	body.SetDisplayName(&req.DisplayName)
	body.SetDescription(req.Description)
	body.SetTags(req.Tags)
	body.SetNotes(req.Notes)

	body.SetSignInAudience(req.SignInAudience)
	body.SetIdentifierUris(req.IdentifierURIs)
	body.SetSamlMetadataUrl(req.SamlMetadataURL)
	body.SetIsFallbackPublicClient(req.IsFallbackPublicClient)
	body.SetGroupMembershipClaims(req.GroupMembershipClaims)
	body.SetIsDeviceOnlyAuthSupported(req.IsDeviceOnlyAuthSupported)
	// FIXME: setting OAuth2RequiredPostResponse is causing creation issues, need to investigate further ( unknown error )
	// body.SetOauth2RequirePostResponse(req.OAuth2RequiredPostResponse)

	// TODO: token encryption key ID is expected to be a UUID, Put it in a function
	// if req.TokenEncryptionKeyID != nil {
	// 	parsedUUID, err := uuid.Parse(*req.TokenEncryptionKeyID)
	// 	if err != nil {
	// 		logger.Error(err, "failed to parse TokenEncryptionKeyID", "value", *req.TokenEncryptionKeyID)
	// 		return nil, err
	// 	}
	// 	body.SetTokenEncryptionKeyId(&parsedUUID)
	// }

	//TODO: Set Optional Claims
	//TODO: API, Web, SPA, PublicClient settings
	//TODO: Required Resource Access and App Roles

	//TODO: Key and Password Credentials
	//TODO: Passwrod Credentials

	//TODO: Informational URLs

	//TODO: Parental Control Settings

	body.SetServiceManagementReference(req.ServiceManagementReference)

	app, err := s.sdk.Applications().Post(ctx, body, nil)
	if err != nil {
		logger.Error(err, "failed to create application", "applicationName", req.DisplayName)
		return nil, err
	}

	appID := ""
	if value := app.GetAppId(); value != nil {
		appID = *value
	}
	objectID := ""
	if value := app.GetId(); value != nil {
		objectID = *value
	}

	logger.Info("application created successfully", "applicationName", req.DisplayName, "clientID", appID, "objectID", objectID)

	return &CreateResponse{
		AppClientID: appID,
		AppObjectID: objectID,
	}, nil
}
