package appregistration

import (
	"context"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Create(ctx context.Context, app appregistration.EntraAppRegistrationSpec) (*AppRegistrationCreateRequest, error) {
	logger := log.FromContext(ctx)
	entraApp := graphmodels.NewApplication()
	entraApp.SetDisplayName(&app.Name)

	client, err := s.sdk.Applications().Post(ctx, entraApp, nil)
	if err != nil {
		logger.Error(err, "failed to create application", "applicationName", app.Name)
		return nil, err
	}

	logger.Info("application created successfully", "applicationName", app.Name, "clientID", *client.GetId(), "principalID", *client.GetAppId())

	return &AppRegistrationCreateRequest{
		AppClientID: *client.GetAppId(),
		AppObjectID: *client.GetId(),
	}, nil
}
