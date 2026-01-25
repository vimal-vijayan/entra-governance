package client

import (
	"context"
	"fmt"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	v1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

type EntraGroupClient interface {
	CreateEntraGroup(entraGroup v1alpha1.EntraSecurityGroupSpec) (*GroupCreateResponse, error)
	GetEntraGroupNameByID(groupID string) (string, error)
	DeleteEntraGroupByID(groupID string) error
	// AddMemberToGroup(groupID, userID string) error
	// RemoveMemberFromGroup(groupID, userID string) error
}

type GroupCreateResponse struct {
	DisplayName string `json:"displayName"`
	ID          string `json:"id"`
}

func (c *GraphClient) CreateEntraGroup(entraGroup v1alpha1.EntraSecurityGroupSpec) (*GroupCreateResponse, error) {

	if err := c.ensureClient(); err != nil {
		return nil, err
	}

	// create the group using the Microsoft Graph SDK

	requestBody := graphmodels.NewGroup()
	description := entraGroup.Description
	displayName := entraGroup.Name
	mailEnabled := entraGroup.MailEnabled
	mailNickname := entraGroup.MailNickname
	securityEnabled := entraGroup.SecurityEnabled
	groupTypes := entraGroup.GroupTypes

	requestBody.SetGroupTypes(groupTypes)
	requestBody.SetDescription(&description)
	requestBody.SetDisplayName(&displayName)
	requestBody.SetMailEnabled(&mailEnabled)
	requestBody.SetMailNickname(&mailNickname)
	requestBody.SetSecurityEnabled(&securityEnabled)

	// Call the SDK to create the group
	groups, err := c.sdk.Groups().Post(context.Background(), requestBody, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return &GroupCreateResponse{
		DisplayName: *groups.GetDisplayName(),
		ID:          *groups.GetId(),
	}, nil
}

func (c *GraphClient) GetEntraGroupNameByID(groupID string) (string, error) {
	if err := c.ensureClient(); err != nil {

		return "Dummy Group Name", nil
	}

	return "", nil
}

func (c *GraphClient) DeleteEntraGroupByID(groupID string) error {
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
