package graph

import (
	"context"
	"fmt"
	"strings"

	graphdirectoryobjects "github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	v1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

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

func (c *GraphClient) GetEntraGroupByID(ctx context.Context, groupID string) (string, string, error) {
	logger := log.FromContext(ctx)
	var HttpStatusCode string

	if groupID == "" {
		logger.Error(fmt.Errorf("groupID cannot be empty"), "invalid groupID")
		return "", "", fmt.Errorf("group id is empty")
	}

	id, err := c.getGroups(ctx, groupID)
	if err != nil {
		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			HttpStatusCode = fmt.Sprintf("%d", odataErr.GetStatusCode())
			logger.Error(err, "failed to get group", "groupID", groupID, "statusCode", odataErr.GetStatusCode())
		}
		logger.Error(err, "failed to get group", "groupID", groupID)
		return "", HttpStatusCode, err
	}

	return id, HttpStatusCode, nil
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

// api doc: https://learn.microsoft.com/en-us/graph/api/directoryobject-checkmembergroups?view=graph-rest-1.0&tabs=go
func (c *GraphClient) CheckGroupMembers(ctx context.Context, groupID string, memberId string) error {
	logger := log.FromContext(ctx)
	requestBody := graphdirectoryobjects.NewItemCheckMemberGroupsPostRequestBody()
	requestBody.SetGroupIds([]string{groupID})
	_, err := c.sdk.DirectoryObjects().ByDirectoryObjectId(memberId).CheckMemberGroups().PostAsCheckMemberGroupsPostResponse(ctx, requestBody, nil)

	if err != nil {
		logger.Info("member is not part of the group", "memberId", memberId, "groupID", groupID)
		return err
	}

	return nil
}

func (c *GraphClient) resolveMemberReference(ctx context.Context, memberID string) (string, error) {
	if memberID == "" {
		return "", fmt.Errorf("member id is empty")
	}

	if strings.HasPrefix(memberID, "https://graph.microsoft.com/") {
		return memberID, nil
	}

	userID, userErr := c.getUsers(ctx, memberID)
	if userErr == nil {
		return fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s", userID), nil
	}

	groupID, groupErr := c.getGroups(ctx, memberID)
	if groupErr == nil {
		return fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s", groupID), nil
	}

	return "", fmt.Errorf("member with ID %s not found as user or group: user error: %v; group error: %v", memberID, userErr, groupErr)
}

func (c *GraphClient) AddOwnersToGroup(ctx context.Context, groupID string) ([]string, error) {
	return nil, nil
}

func (c *GraphClient) RemoveOwnersFromGroup(ctx context.Context, groupID string) error {
	return nil
}

func (c *GraphClient) AddMembersToGroup(ctx context.Context, groupID string, resourceType string, memberRefs []string) error {

	if len(memberRefs) == 0 {
		return nil
	}

	bindRefs := make([]string, 0, len(memberRefs))
	for _, memberId := range memberRefs {
		if strings.HasPrefix(memberId, "https://graph.microsoft.com/") {
			bindRefs = append(bindRefs, memberId)
			continue
		}
		bindRefs = append(bindRefs, fmt.Sprintf("https://graph.microsoft.com/v1.0/%s/%s", resourceType, memberId))
	}

	groups := models.NewGroup()
	additionalData := map[string]any{
		"members@odata.bind": bindRefs,
	}
	groups.SetAdditionalData(additionalData)
	_, err := c.sdk.Groups().ByGroupId(groupID).Patch(ctx, groups, nil)
	if err != nil {
		// batch addtion can fail if any of the members already exist in the group or anything invalid
		// try adding members one by one
		for _, memberId := range memberRefs {
			member, er := c.resolveMemberReference(ctx, memberId)
			if er != nil {
				return fmt.Errorf("failed to resolve member reference: %v", er)
			}
			additionalData := map[string]any{
				"members@odata.bind": []string{member},
			}
			groups.SetAdditionalData(additionalData)
			_, err := c.sdk.Groups().ByGroupId(groupID).Patch(ctx, groups, nil)
			if err != nil {
				odataErr, ok := err.(*odataerrors.ODataError)
				if ok && (odataErr.GetStatusCode() == 400 || odataErr.GetStatusCode() == 409) {
					logger := log.FromContext(ctx)
					logger.Info("member already exists in group or bad request", "memberId", memberId, "statusCode", odataErr.GetStatusCode())
					continue
				}
				return fmt.Errorf("failed to add member %s to group: %v", memberId, err)
			}
		}
	}
	return nil
}

func (c *GraphClient) RemoveMembersFromGroup(ctx context.Context, groupID string) error {
	return nil
}
