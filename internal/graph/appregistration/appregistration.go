package appregistration

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Service struct {
	sdk *msgraphsdk.GraphServiceClient
}

type AppRegistrationGetResponse struct {
	ID             string
	DisplayName    string
	HttpStatusCode string
}

type AppRegistrationCreateRequest struct {
	AppClientID string
	AppObjectID string
}

type API interface {
	Get(ctx context.Context, appID string) (*AppRegistrationGetResponse, error)
	Create(ctx context.Context, app appregistration.EntraAppRegistrationSpec) (*AppRegistrationCreateRequest, error)
	Delete(ctx context.Context, appID string) error
}

func (s *Service) Get(ctx context.Context, appID string) (*AppRegistrationGetResponse, error) {
	return nil, nil
}

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
		AppClientID: *client.GetId(),
		AppObjectID: *client.GetAppId(),
	}, nil
}

func (s *Service) Delete(ctx context.Context, appID string) error {
	logger := log.FromContext(ctx)
	err := s.sdk.Applications().ByApplicationId(appID).Delete(ctx, nil)
	if err != nil {
		logger.Error(err, "failed to delete application", "applicationID", appID)
		return err
	}

	logger.Info("application deleted successfully", "applicationID", appID)
	return nil
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}
