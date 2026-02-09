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
	Get(ctx context.Context, appID string) (*graphmodels.Application, error)
	Create(ctx context.Context, app appregistration.EntraAppRegistrationSpec) (*AppRegistrationCreateRequest, error)
	Delete(ctx context.Context, appID string) error
	Update(ctx context.Context, app appregistration.EntraAppRegistration) error
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}

func (s *Service) Get(ctx context.Context, appID string) (*graphmodels.Application, error) {
	resp, err := s.sdk.Applications().ByApplicationId(appID).Get(ctx, nil)	
	if err != nil {
		return nil, err
	}
	return resp.(*graphmodels.Application), nil
}

func (s *Service) Create(ctx context.Context, app appregistration.EntraAppRegistrationSpec) (*AppRegistrationCreateRequest, error) {
	logger := log.FromContext(ctx)
	entraApp := graphmodels.NewApplication()
	entraApp.SetDisplayName(&app.Name)
	entraApp.SetDescription(&app.Description)
	entraApp.SetTags(app.Tags)

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


func (s *Service) Update(ctx context.Context, app appregistration.EntraAppRegistration) error {
	logger := log.FromContext(ctx)
	appSpec := app.Spec

	requestbody := generateBasicRequestBody(appSpec)

	_, err := s.sdk.Applications().ByApplicationId(app.Status.AppRegistrationObjID).Patch(ctx, requestbody, nil)

	if err != nil {
		logger.Error(err, "failed to update application", "applicationName", app.Spec.Name, "applicationID", app.Status.AppRegistrationID)
		return err
	}

	logger.Info("application updated successfully", "applicationName", app.Spec.Name, "applicationID", app.Status.AppRegistrationID)
	return nil
}



func generateBasicRequestBody(app appregistration.EntraAppRegistrationSpec) *graphmodels.Application {
	requestbody := graphmodels.NewApplication()
	requestbody.SetDisplayName(&app.Name)
	requestbody.SetDescription(&app.Description)
	requestbody.SetTags(app.Tags)
	requestbody.SetSignInAudience(&app.SignInAudience)
	requestbody.SetSamlMetadataUrl(&app.SamlMetadataUrl)
	requestbody.SetIsFallbackPublicClient(&app.IsFallbackPublicClient)
	requestbody.SetIsDeviceOnlyAuthSupported(&app.IsDeviceOnlyAuthSupported)
	requestbody.SetGroupMembershipClaims(&app.GroupMembershipClaims)
	return requestbody
}
