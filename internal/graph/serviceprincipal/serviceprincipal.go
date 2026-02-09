package serviceprincipal

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	graphserviceprincipalswithappid "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipalswithappid"
	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Service struct {
	sdk *msgraphsdk.GraphServiceClient
}

type ServicePrincipalGetResponse struct {
	ID             string
	DisplayName    string
	AppID          string
	HttpStatusCode int
}

type ServicePrincipalCreateResponse struct {
	ServicePrincipalID string
}

type API interface {
	Get(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error)
	Create(ctx context.Context, app appregistration.EntraAppRegistration) (*ServicePrincipalCreateResponse, error)
	Delete(ctx context.Context, appObjectID string) error
	Update(ctx context.Context, appObjectID string) error
	Upsert(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error)
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}

func (s *Service) Create(ctx context.Context, app appregistration.EntraAppRegistration) (*ServicePrincipalCreateResponse, error) {
	logger := log.FromContext(ctx)

	headers := abstractions.NewRequestHeaders()
	headers.Add("Prefer", "create-if-missing")
	configuration := &graphserviceprincipalswithappid.ServicePrincipalsWithAppIdRequestBuilderPatchRequestConfiguration{
		Headers: headers,
	}

	requestBody := graphmodels.NewServicePrincipal()
	displayName := app.Spec.Name
	requestBody.SetDisplayName(&displayName)
	requestBody.SetTags([]string{"HideApp"})
	appId := app.Status.AppRegistrationID
	response, err := s.sdk.ServicePrincipalsWithAppId(&appId).Patch(ctx, requestBody, configuration)

	if err != nil {
		logger.Error(err, "failed to create or update service principal", "applicationID", app.Status.AppRegistrationID)
		return nil, err
	}

	// Check if response or ID is nil before dereferencing
	if response == nil {
		logger.Error(nil, "service principal response is nil", "applicationID", app.Status.AppRegistrationID)
		return nil, fmt.Errorf("service principal response is nil")
	}

	servicePrincipalID := response.GetId()
	if servicePrincipalID == nil {
		logger.Error(nil, "service principal ID is nil", "applicationID", app.Status.AppRegistrationID)
		return nil, fmt.Errorf("service principal ID is nil in response")
	}

	logger.Info("service principal created or updated successfully", "applicationID", app.Status.AppRegistrationID, "servicePrincipalID", *servicePrincipalID)

	return &ServicePrincipalCreateResponse{
		ServicePrincipalID: *servicePrincipalID,
	}, nil
}

func (s *Service) Get(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, appObjectID string) error {
	return nil
}

func (s *Service) Update(ctx context.Context, appObjectID string) error {
	return nil
}

func (s *Service) Upsert(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error) {
	return nil, nil
}
