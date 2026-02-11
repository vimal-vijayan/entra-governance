package serviceprincipal

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
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
	Create(ctx context.Context, req CreateRequest) (*ServicePrincipalCreateResponse, error)
	Delete(ctx context.Context, appObjectID string) error
	Update(ctx context.Context, appObjectID string) error
	Upsert(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error)
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}

func (s *Service) Get(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, appObjectID string) error {
	if err := s.sdk.ServicePrincipals().ByServicePrincipalId(appObjectID).Delete(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, appObjectID string) error {
	return nil
}

func (s *Service) Upsert(ctx context.Context, appObjectID string) (*ServicePrincipalGetResponse, error) {
	return nil, nil
}
