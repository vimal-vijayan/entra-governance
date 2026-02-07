package graph

import (
	"context"

	v1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

type GroupCreateResponse struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

type EntraGroupClient interface {
	CreateEntraGroup(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroupSpec) (*GroupCreateResponse, error)
	GetEntraGroupByID(ctx context.Context, groupID string) (string, error)
	DeleteEntraGroupByID(ctx context.Context, groupID string) error
	AddMembersToGroup(ctx context.Context, groupID string, resourceType string, memberRefs []string) error
	CheckGroupMembers(ctx context.Context, groupID string, memberId string) error
}
