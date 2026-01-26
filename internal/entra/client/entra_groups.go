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

	owners := []string{}
	if entraGroup.Owners != nil {
		for _, userObjId := range entraGroup.Owners {
			userID, err := c.getUsers(ctx, userObjId)
			if err != nil {
				return nil, fmt.Errorf("failed to get owner user ID: check if user exists or provide valid user object ID instead of user principal name: %v", err)
			}
			owners = append(owners, fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userID))
		}
	}

	members := []string{}
	if entraGroup.Members != nil {
		for _, memberID := range entraGroup.Members {
			memberRef, err := c.resolveMemberReference(ctx, memberID)
			if err != nil {
				return nil, fmt.Errorf("failed to get member user/group ID: check if user/group exists or provide valid object ID instead of user principal name: %v", err)
			}
			members = append(members, memberRef)
		}
	}

	additionalData := map[string]any{
		"owners@odata.bind":  owners,
		"members@odata.bind": members,
	}

	group.SetAdditionalData(additionalData)

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

	_, err := c.getGroups(ctx, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (c *GraphClient) DeleteEntraGroupByID(ctx context.Context, groupID string) error {

	if err := c.ensureClient(); err != nil {
		return err
	}

	if groupID == "" {
		return fmt.Errorf("group id is empty")
	}

	err := c.sdk.Groups().ByGroupId(groupID).Delete(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete group by ID: %v", err)
	}

	return nil
}

func (c *GraphClient) ensureClient() error {
	if c == nil || c.sdk == nil {
		return fmt.Errorf("Graph client is not initialized")
	}
	return nil
}

func (c *GraphClient) resolveMemberReference(ctx context.Context, memberID string) (string, error) {

	if userID, err := c.getUsers(ctx, memberID); err == nil {
		return fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userID), nil
	}

	if groupID, err := c.getGroups(ctx, memberID); err == nil {
		return fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s", groupID), nil
	}

	return "", fmt.Errorf("could not resolve member '%s' as user or group", memberID)
}
