package appregistration

import (
	"context"
	"fmt"
)

func (s *Service) Get(ctx context.Context, objectID string) (*Application, error) {
	resp, err := s.sdk.Applications().ByApplicationId(objectID).Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("application response is nil for objectID %q", objectID)
	}

	state := &Application{}
	if value := resp.GetId(); value != nil {
		state.ObjectID = *value
	}
	if value := resp.GetAppId(); value != nil {
		state.AppID = *value
	}
	if value := resp.GetDisplayName(); value != nil {
		state.DisplayName = *value
	}
	if value := resp.GetDescription(); value != nil {
		state.Description = *value
	}
	if value := resp.GetTags(); value != nil {
		state.Tags = append([]string(nil), value...)
	}
	if value := resp.GetSignInAudience(); value != nil {
		state.SignInAudience = *value
	}
	if value := resp.GetSamlMetadataUrl(); value != nil {
		state.SamlMetadataURL = *value
	}
	if value := resp.GetIsFallbackPublicClient(); value != nil {
		state.IsFallbackPublicClient = *value
	}
	if value := resp.GetIsDeviceOnlyAuthSupported(); value != nil {
		state.IsDeviceOnlyAuthSupported = *value
	}
	if value := resp.GetGroupMembershipClaims(); value != nil {
		state.GroupMembershipClaims = *value
	}

	return state, nil
}
