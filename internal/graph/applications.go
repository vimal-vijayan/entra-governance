package graph

import (
	"context"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Application operations will be implemented here

type ApplicationCreateResponse struct {
	ClientID    string
	PrincipalID string
}

func (c *GraphClient) CreateEntraApplication(ctx context.Context, entraApp appregistration.EntraAppRegistrationSpec) (*ApplicationCreateResponse, error) {

	logger := log.FromContext(ctx)

	app := graphmodels.NewApplication()
	app.SetDisplayName(&entraApp.Name)

	client, err := c.sdk.Applications().Post(ctx, app, nil)
	if err != nil {
		logger.Error(err, "failed to create application", "applicationName", entraApp.Name)
		return nil, err
	}

	logger.Info("application created successfully", "applicationName", entraApp.Name, "clientID", *client.GetId(), "principalID", *client.GetAppId())

	return &ApplicationCreateResponse{
		ClientID:    *client.GetId(),
		PrincipalID: *client.GetAppId(),
	}, nil
}

func (c *GraphClient) DeleteEntraApplicationByID(ctx context.Context, appID string) error {
	logger := log.FromContext(ctx)

	err := c.sdk.Applications().ByApplicationId(appID).Delete(ctx, nil)

	if err != nil {
		logger.Error(err, "failed to delete application", "applicationID", appID)
		return err
	}

	logger.Info("application deleted successfully", "applicationID", appID)
	return nil
}
