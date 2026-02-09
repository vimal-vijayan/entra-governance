package serviceprincipal

import (
	"context"
	"fmt"

	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/client"
	"github.com/vimal-vijayan/entra-governance/internal/graph/serviceprincipal"
)

type Service struct {
	factory *client.ClientFactory
}

func NewService(factory *client.ClientFactory) *Service {
	return &Service{factory: factory}
}

func (s *Service) getGraphClient(ctx context.Context, entraApp appregistration.EntraAppRegistration) (*client.GraphClient, error) {
	if entraApp.Spec.ForProvider == nil {
		return nil, fmt.Errorf("forProvider spec is nil")
	}

	if entraApp.Spec.ForProvider.CredentialSecretRef == "" {
		return nil, fmt.Errorf("credential secret reference is empty in forProvider spec")
	}

	secretRef := client.SecretRef{
		Name:      entraApp.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraApp.Namespace,
	}

	sdk, err := s.factory.ForClientSecret(ctx, secretRef)
	if err != nil {
		return nil, err
	}

	return client.NewGraphClient(sdk), nil
}

func (s *Service) Create(ctx context.Context, entraApp appregistration.EntraAppRegistration) (string, error) {

	displayName := entraApp.Spec.Name
	spnParameters := entraApp.Spec.ServicePrincipal

	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return "", err
	}

	resp, err := graphClient.ServicePrincipals.Create(ctx, serviceprincipal.CreateRequest{
		DisplayName:                displayName,
		AppID:                      entraApp.Status.AppRegistrationID,
		DisableVisibilityForGuests: spnParameters.DisableVisibilityForGuests,
		AccountEnabled:             spnParameters.AccountEnabled,
	})

	if err != nil {
		return "", err
	}

	return resp.ServicePrincipalID, nil
}

func (s *Service) Delete(ctx context.Context, appID string, entraApp appregistration.EntraAppRegistration) error {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return err
	}

	return graphClient.AppRegistration.Delete(ctx, appID)
}

type appRegistrationResponse struct {
	AppClientID string
	AppObjectID string
}

func (s *Service) Update(ctx context.Context, entraApp appregistration.EntraAppRegistration) (*appRegistrationResponse, error) {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return nil, err
	}

	err = graphClient.ServicePrincipals.Update(ctx, "objectid")
	if err != nil {
		return nil, err
	}

	// ensure tags
	if err = s.upsertTags(ctx, entraApp); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Service) upsertTags(ctx context.Context, entraApp appregistration.EntraAppRegistration) error {
	return nil
}
