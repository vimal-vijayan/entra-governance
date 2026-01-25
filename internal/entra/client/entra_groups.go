package client

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	v1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type GroupCreateResponse struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

func (c *GraphClient) CreateEntraGroup(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroupSpec) (*GroupCreateResponse, error) {

	group := models.NewGroup()
	group.SetDisplayName(&entraGroup.Name)
	group.SetDescription(&entraGroup.Description)
	group.SetMailEnabled(&entraGroup.MailEnabled)
	group.SetMailNickname(&entraGroup.MailNickname)
	group.SetSecurityEnabled(&entraGroup.SecurityEnabled)
	group.SetGroupTypes(entraGroup.GroupTypes)

	// Call the SDK to create the group
	groups, err := c.sdk.Groups().Post(ctx, group, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return &GroupCreateResponse{
		DisplayName: *groups.GetDisplayName(),
		ID:          *groups.GetId(),
	}, nil
}

func (c *GraphClient) GetEntraGroupByID(ctx context.Context, groupID string) error {
	logger := log.FromContext(ctx)

	if groupID == "" {
		logger.Error(fmt.Errorf("groupID cannot be empty"), "invalid groupID")
		return fmt.Errorf("group id is empty")
	}

	groups, err := c.sdk.Groups().ByGroupId(groupID).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get group by ID: %v", err)
	}

	if groups.GetId() == nil {
		return fmt.Errorf("group with ID %s not found", groupID)
	}

	logger.Info("successfully fetched group", "groupID", *groups.GetId())
	return nil
}

func (c *GraphClient) DeleteEntraGroupByID(ctx context.Context, groupID string) error {

	if err := c.ensureClient(); err != nil {
		return err
	}

	// call the function to delete group by ID
	return nil
}

func (c *GraphClient) ensureClient() error {
	if c == nil || c.sdk == nil {
		return fmt.Errorf("Graph client is not initialized")
	}
	return nil
}
