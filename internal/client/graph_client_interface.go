package client

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	v1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

type GroupCreateResponse struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

type ApplicationCreateResponse struct {
	ClientID    string `json:"clientId"`
	PrincipalID string `json:"principalId"`
}

type EntraGroupClient interface {

	// Group operations
	CreateEntraGroup(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroupSpec) (*GroupCreateResponse, error)
	GetEntraGroupByID(ctx context.Context, groupID string) (string, error)
	DeleteEntraGroupByID(ctx context.Context, groupID string) error
	AddMembersToGroup(ctx context.Context, groupID string, resourceType string, memberRefs []string) error
	CheckGroupMembers(ctx context.Context, groupID string, memberId string) error

	// Application operations
	CreateEntraApplication(ctx context.Context, entraApp v1alpha1.EntraAppRegistrationSpec) (*ApplicationCreateResponse, error)
	DeleteEntraApplicationByID(ctx context.Context, appID string) error
}

type GraphClient struct {
	sdk *msgraphsdk.GraphServiceClient
}

func NewGraphClient(sdk *msgraphsdk.GraphServiceClient) *GraphClient {
	return &GraphClient{sdk: sdk}
}
