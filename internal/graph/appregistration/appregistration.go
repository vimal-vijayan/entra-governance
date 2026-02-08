package appregistration

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
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
	Update(ctx context.Context, app appregistration.EntraAppRegistration) error
}

func (s *Service) Get(ctx context.Context, appID string) (*AppRegistrationGetResponse, error) {
	return nil, nil
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}
