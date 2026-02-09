package appregistration

import (
	"context"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Update(ctx context.Context, app appregistration.EntraAppRegistration) error {
	logger := log.FromContext(ctx)
	appSpec := app.Spec

	requestbody := generateBasicRequestBody(appSpec)

	_, err := s.sdk.Applications().ByApplicationId(app.Status.AppRegistrationObjID).Patch(ctx, requestbody, nil)

	if err != nil {
		logger.Error(err, "failed to update application", "applicationName", app.Spec.Name, "applicationID", app.Status.AppRegistrationID)
		return err
	}

	logger.Info("application updated successfully", "applicationName", app.Spec.Name, "applicationID", app.Status.AppRegistrationID)
	return nil
}

func generateBasicRequestBody(app appregistration.EntraAppRegistrationSpec) *graphmodels.Application {
	requestbody := graphmodels.NewApplication()
	requestbody.SetDisplayName(&app.Name)
	requestbody.SetDescription(&app.Description)
	requestbody.SetTags(app.Tags)
	requestbody.SetSignInAudience(&app.SignInAudience)
	requestbody.SetSamlMetadataUrl(&app.SamlMetadataUrl)
	requestbody.SetIsFallbackPublicClient(&app.IsFallbackPublicClient)
	requestbody.SetIsDeviceOnlyAuthSupported(&app.IsDeviceOnlyAuthSupported)
	requestbody.SetGroupMembershipClaims(&app.GroupMembershipClaims)
	return requestbody
}
