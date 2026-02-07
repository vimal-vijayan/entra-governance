package groups

import (
	"context"
	"fmt"

	"github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/client"
	"github.com/vimal-vijayan/entra-governance/internal/graph"
)

type Service struct {
	factory *client.ClientFactory
}

func NewService(factory *client.ClientFactory) *Service {
	return &Service{factory: factory}
}

func (s *Service) Get(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroup, groupID string) (string, string, error) {

	if entraGroup.Spec.ForProvider == nil {
		return "", "", fmt.Errorf("credential reference in forProvider spec is nil")
	}

	secretRef := client.SecretRef{
		Name:      entraGroup.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraGroup.Namespace,
	}

	sdk, err := s.factory.ForClientSecret(ctx, secretRef)
	if err != nil {
		return "", "", err
	}

	graphClient := graph.NewGraphClient(sdk)
	id, statusCode, err := graphClient.GetEntraGroupByID(ctx, groupID)

	if err != nil {
		return "", statusCode, err
	}

	return id, statusCode, nil
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

		graphClient := graph.NewGraphClient(sdk)
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

	graphClient := graph.NewGraphClient(sdk)
	return graphClient.DeleteEntraGroupByID(ctx, groupID)
}

func (s *Service) AddMembers(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroup) error {

	if entraGroup.Spec.Members == nil || len(*entraGroup.Spec.Members) == 0 {
		return nil
	}

	// get users as members from the spec
	var userIDs = getMemberIDs(entraGroup, "User")
	// Fetch group ID using group name
	var groupIDs = getMemberIDs(entraGroup, "Group")
	// Fetch service principal IDs using spec
	var servicePrincipalIDs = getMemberIDs(entraGroup, "ServicePrincipal")

	sdk, err := s.factory.ForClientSecret(ctx, client.SecretRef{
		Name:      entraGroup.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraGroup.Namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to create SDK client: %v", err)
	}

	graphClient := graph.NewGraphClient(sdk)

	// call the api to add members
	if userIDs != nil {
		err := graphClient.AddMembersToGroup(ctx, entraGroup.Status.ID, "users", userIDs)
		if err != nil {
			return fmt.Errorf("failed to add users as members to group: %v", err)
		}
	}

	if groupIDs != nil {
		// add groups as members
	}

	if servicePrincipalIDs != nil {
		// add service principals as members
	}

	return nil
}

func (s *Service) CheckMemberIds(ctx context.Context, entraGroup v1alpha1.EntraSecurityGroup) ([]string, error) {

	currentManagedMemberIds := entraGroup.Status.ManagedMemberGroups

	sdk, err := s.factory.ForClientSecret(ctx, client.SecretRef{
		Name:      entraGroup.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraGroup.Namespace,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create SDK client: %v", err)
	}

	graphClient := graph.NewGraphClient(sdk)

	unManagedMemberIds := []string{}
	for _, member := range currentManagedMemberIds {
		err := graphClient.CheckGroupMembers(ctx, entraGroup.Status.ID, member)
		if err != nil {
			unManagedMemberIds = append(unManagedMemberIds, member)
		}
	}

	return unManagedMemberIds, nil
}

func getMemberIDs(entraGroup v1alpha1.EntraSecurityGroup, Type string) []string {
	var ids []string
	for _, member := range *entraGroup.Spec.Members {
		if member.Type == Type {
			ids = append(ids, member.Id)
		}
	}
	return ids
}
