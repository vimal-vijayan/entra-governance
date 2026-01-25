package client

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	v1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

type EntraGroupClient interface {
	CreateEntraGroup(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroupSpec) (*GroupCreateResponse, error)
	GetEntraGroupByID(ctx context.Context, groupID string) (string, error)
	DeleteEntraGroupByID(ctx context.Context, groupID string) error
	// AddMemberToGroup(groupID, userID string) error
	// RemoveMemberFromGroup(groupID, userID string) error
}

type GraphClient struct {
	sdk *msgraphsdk.GraphServiceClient
}

func NewGraphClient(sdk *msgraphsdk.GraphServiceClient) *GraphClient {
	return &GraphClient{sdk: sdk}
}
