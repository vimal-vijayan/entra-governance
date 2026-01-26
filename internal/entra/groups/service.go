package groups

import (
	"context"
	"fmt"

	"github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/entra/client"
)

type Service struct {
	factory *client.ClientFactory
}

func NewService(factory *client.ClientFactory) *Service {
	return &Service{factory: factory}
}

func (s *Service) Create(ctx context.Context, groupSpec v1alpha1.EntraSecurityGroup) (string, string, error) {

	if groupSpec.Spec.ForProvider == nil {
		return "", "", fmt.Errorf("forProvider spec is nil")
	}

	secretRef := client.SecretRef{
		Name:      groupSpec.Spec.ForProvider.CredentialSecretRef,
		Namespace: groupSpec.Namespace,
	}

	if groupSpec.Spec.ForProvider.CredentialSecretRef != "" {

		sdk, err := s.factory.ForClientSecret(ctx, secretRef)
		if err != nil {
			return "", "", err
		}

		graphClient := client.NewGraphClient(sdk)
		resp, err := graphClient.CreateEntraGroup(ctx, groupSpec.Spec)
		if err != nil {
			return "", "", err
		}

		return resp.ID, resp.DisplayName, nil
	}

	if groupSpec.Spec.ForProvider.ServiceAccountRef != "" {
		// sa := groupSpec.Spec.ForProvider.ServiceAccountRef
		//TODO: check if service account ref is valid

		return "", "", nil
	}

	return "", "", fmt.Errorf("no valid credential reference found in the EntraSecurityGroup spec")
}

func (s *Service) Get(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroup, groupID string) (string, error) {

	if entraGroup.Spec.ForProvider == nil {
		return "", fmt.Errorf("forProvider spec is nil")
	}

	secretRef := client.SecretRef{
		Name:      entraGroup.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraGroup.Namespace,
	}

	sdk, err := s.factory.ForClientSecret(ctx, secretRef)
	if err != nil {
		return "", fmt.Errorf("failed to create SDK client: %v", err)
	}

	graphClient := client.NewGraphClient(sdk)
	statusCode, err := graphClient.GetEntraGroupByID(ctx, groupID)
	if err != nil {
		return statusCode, err
	}
	// You might want to do something with the status here
	return statusCode, nil
}

func (s *Service) Delete(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroup, groupID string) error {

	if entraGroup.Spec.ForProvider == nil {
		return fmt.Errorf("forProvider spec is nil")
	}

	secretRef := client.SecretRef{
		Name:      entraGroup.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraGroup.Namespace,
	}

	sdk, err := s.factory.ForClientSecret(ctx, secretRef)
	if err != nil {
		return fmt.Errorf("failed to create SDK client: %v", err)
	}

	graphClient := client.NewGraphClient(sdk)
	return graphClient.DeleteEntraGroupByID(ctx, groupID)
}
