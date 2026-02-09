package serviceprincipal

import (
	"context"
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

type CreateRequest struct {
	DisplayName                string
	AppID                      string
	DisableVisibilityForGuests bool
}

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
	Create(ctx context.Context, req CreateRequest) (*ServicePrincipalCreateResponse, error)
	Delete(ctx context.Context, appObjectID string) error
	Update(ctx context.Context, appObjectID string) error
	Upsert(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error)
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*ServicePrincipalCreateResponse, error) {
	logger := log.FromContext(ctx)

	requestBody := graphmodels.NewServicePrincipal()
	appId := req.AppID
	displayName := req.DisplayName
	requestBody.SetDisplayName(&displayName)
	if req.DisableVisibilityForGuests {
		requestBody.SetTags([]string{"HideApp"})
	}
	requestBody.SetAppId(&appId)
	response, err := s.sdk.ServicePrincipals().Post(ctx, requestBody, nil)

	if err != nil {
		logger.Error(err, "failed to create service principal", "applicationID", appId)
		return nil, err
	}

	// Check if response or ID is nil before dereferencing
	if response == nil {
		logger.Error(nil, "service principal response is nil", "applicationID", appId)
		return nil, fmt.Errorf("service principal response is nil")
	}

	servicePrincipalID := response.GetId()
	if servicePrincipalID == nil {
		logger.Error(nil, "service principal ID is nil", "applicationID", appId)
		return nil, fmt.Errorf("service principal ID is nil in response")
	}

	logger.Info("service principal created successfully", "applicationID", appId, "servicePrincipalID", *servicePrincipalID)

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
