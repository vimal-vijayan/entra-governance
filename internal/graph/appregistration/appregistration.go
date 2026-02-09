package appregistration

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type Service struct {
	sdk *msgraphsdk.GraphServiceClient
}

type Application struct {
	ObjectID                  string
	AppID                     string
	DisplayName               string
	Description               string
	Tags                      []string
	SignInAudience            string
	SamlMetadataURL           string
	IsFallbackPublicClient    bool
	IsDeviceOnlyAuthSupported bool
	GroupMembershipClaims     string
}

type CreateRequest struct {
	DisplayName string
	Description string
	Tags        []string

	SignInAudience            string
	SamlMetadataURL           string
	IsFallbackPublicClient    bool
	IsDeviceOnlyAuthSupported bool
	GroupMembershipClaims     string
}

type CreateResponse struct {
	AppClientID string
	AppObjectID string
}

type PatchRequest struct {
	ObjectID string

	DisplayName               *string
	Description               *string
	Tags                      *[]string
	SignInAudience            *string
	SamlMetadataURL           *string
	IsFallbackPublicClient    *bool
	IsDeviceOnlyAuthSupported *bool
	GroupMembershipClaims     *string
}

type API interface {
	Get(ctx context.Context, objectID string) (*Application, error)
	Create(ctx context.Context, req CreateRequest) (*CreateResponse, error)
	Patch(ctx context.Context, req PatchRequest) error
	Delete(ctx context.Context, objectID string) error
}

func NewAPI(sdk *msgraphsdk.GraphServiceClient) API {
	return &Service{sdk: sdk}
}
