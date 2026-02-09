package appregistration

import (
	"context"
	"fmt"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Patch(ctx context.Context, req PatchRequest) error {
	logger := log.FromContext(ctx)

	if req.ObjectID == "" {
		return fmt.Errorf("objectID is required for patch")
	}

	body := graphmodels.NewApplication()
	hasChanges := false

	if req.DisplayName != nil {
		body.SetDisplayName(req.DisplayName)
		hasChanges = true
	}
	if req.Description != nil {
		body.SetDescription(req.Description)
		hasChanges = true
	}
	if req.Tags != nil {
		body.SetTags(*req.Tags)
		hasChanges = true
	}
	if req.SignInAudience != nil {
		body.SetSignInAudience(req.SignInAudience)
		hasChanges = true
	}
	if req.SamlMetadataURL != nil {
		body.SetSamlMetadataUrl(req.SamlMetadataURL)
		hasChanges = true
	}
	if req.IsFallbackPublicClient != nil {
		body.SetIsFallbackPublicClient(req.IsFallbackPublicClient)
		hasChanges = true
	}
	if req.IsDeviceOnlyAuthSupported != nil {
		body.SetIsDeviceOnlyAuthSupported(req.IsDeviceOnlyAuthSupported)
		hasChanges = true
	}
	if req.GroupMembershipClaims != nil {
		body.SetGroupMembershipClaims(req.GroupMembershipClaims)
		hasChanges = true
	}

	if !hasChanges {
		return nil
	}

	_, err := s.sdk.Applications().ByApplicationId(req.ObjectID).Patch(ctx, body, nil)
	if err != nil {
		logger.Error(err, "failed to patch application", "objectID", req.ObjectID)
		return err
	}

	logger.Info("application patched successfully", "objectID", req.ObjectID)
	return nil
}
