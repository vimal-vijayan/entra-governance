package groups

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	entraGroup "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

func (s *Service) Create(ctx context.Context, groupSpec entraGroup.EntraSecurityGroupSpec) (*GroupCreateResponse, error) {

	group := models.NewGroup()
	group.SetDisplayName(&groupSpec.Name)
	group.SetDescription(&groupSpec.Description)
	group.SetMailEnabled(&groupSpec.MailEnabled)
	group.SetMailNickname(&groupSpec.MailNickname)
	group.SetSecurityEnabled(&groupSpec.SecurityEnabled)
	group.SetGroupTypes(groupSpec.GroupTypes)

	resp, err := s.sdk.Groups().Post(ctx, group, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	return &GroupCreateResponse{
		DisplayName: *resp.GetDisplayName(),
		ID:          *resp.GetId(),
	}, nil
}
