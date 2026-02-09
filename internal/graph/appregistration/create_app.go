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
	body.SetDescription(&req.Description)
	body.SetTags(req.Tags)
	body.SetSignInAudience(&req.SignInAudience)
	body.SetSamlMetadataUrl(&req.SamlMetadataURL)
	body.SetIsFallbackPublicClient(&req.IsFallbackPublicClient)
	body.SetIsDeviceOnlyAuthSupported(&req.IsDeviceOnlyAuthSupported)
	body.SetGroupMembershipClaims(&req.GroupMembershipClaims)

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
