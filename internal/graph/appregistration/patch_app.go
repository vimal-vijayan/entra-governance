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

	hasChanges = applyString(req.DisplayName, body.SetDisplayName) || hasChanges
	hasChanges = applyString(req.Description, body.SetDescription) || hasChanges
	hasChanges = applyStringSlice(req.Tags, body.SetTags) || hasChanges
	hasChanges = applyString(req.SignInAudience, body.SetSignInAudience) || hasChanges
	hasChanges = applyString(req.SamlMetadataURL, body.SetSamlMetadataUrl) || hasChanges
	hasChanges = applyBool(req.IsFallbackPublicClient, body.SetIsFallbackPublicClient) || hasChanges
	hasChanges = applyBool(req.IsDeviceOnlyAuthSupported, body.SetIsDeviceOnlyAuthSupported) || hasChanges
	hasChanges = applyString(req.GroupMembershipClaims, body.SetGroupMembershipClaims) || hasChanges

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
