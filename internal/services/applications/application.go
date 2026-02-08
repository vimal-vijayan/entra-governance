package applications

import (
	"context"
	"fmt"

	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/client"
)

type Service struct {
	factory *client.ClientFactory
}

func NewService(factory *client.ClientFactory) *Service {
	return &Service{factory: factory}
}

func (s *Service) Create(ctx context.Context, entraApp appregistration.EntraAppRegistration) (string, string, error) {

	if entraApp.Spec.ForProvider == nil {
		return "", "", fmt.Errorf("forProvider spec is nil")
	}

	secretRef := client.SecretRef{
		Name:      entraApp.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraApp.Namespace,
	}

	if entraApp.Spec.ForProvider.CredentialSecretRef != "" {
		sdk, err := s.factory.ForClientSecret(ctx, secretRef)
		if err != nil {
			return "", "", err
		}

		graphClient := client.NewGraphClient(sdk)
		// response, err := graphClient.CreateEntraApplication(ctx, entraApp.Spec)
		response, err := graphClient.AppRegistration.Create(ctx, entraApp.Spec)
		if err != nil {
			return "", "", err
		}
		return response.AppClientID, response.AppObjectID, nil
	}

	return "", "", fmt.Errorf("credential secret reference is empty in forProvider spec")
}

func (s *Service) Delete(ctx context.Context, appID string, entraApp appregistration.EntraAppRegistration) error {

	if entraApp.Spec.ForProvider == nil {
		return fmt.Errorf("forProvider spec is nil")
	}

	secretRef := client.SecretRef{
		Name:      entraApp.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraApp.Namespace,
	}

	if entraApp.Spec.ForProvider.CredentialSecretRef != "" {
		sdk, err := s.factory.ForClientSecret(ctx, secretRef)
		if err != nil {
			return err
		}

		graphClient := client.NewGraphClient(sdk)
		return graphClient.AppRegistration.Delete(ctx, appID)
	}

	return fmt.Errorf("credential secret reference is empty in forProvider spec")
}
