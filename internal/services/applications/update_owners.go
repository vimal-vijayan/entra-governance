package applications

import (
	"context"

	entrav1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
)

// UpdateOwners manages the owners of an application registration
func (s *Service) UpdateOwners(ctx context.Context, appId string, entraApp entrav1alpha1.EntraAppRegistration) error {

	owners := append([]string(nil), (*entraApp.Spec.Owners)...)
	managedOwners := append([]string(nil), entraApp.Status.Owners...)

	if len(owners) == 0 && len(managedOwners) == 0 {
		return nil
	}

	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return err
	}

	if !equalStringSets(owners, managedOwners) {
		ownersToAdd := findMissingOwners(owners, managedOwners)
		if len(ownersToAdd) > 0 {
			_, err := graphClient.AppRegistration.AddAppOwners(ctx, appId, ownersToAdd)
			if err != nil {
				return err
			}
		}

		ownersToRemove := findMissingOwners(managedOwners, owners)
		if len(ownersToRemove) > 0 {
			err := graphClient.AppRegistration.RemoveAppOwners(ctx, appId, ownersToRemove)
			if err != nil {
				return err
			}
		}
	} else {
		currentOwners, err := graphClient.AppRegistration.GetAppOwners(ctx, appId)
		if err != nil {
			return err
		}
		ownersToAdd := findMissingOwners(owners, currentOwners)
		if len(ownersToAdd) > 0 {
			_, err := graphClient.AppRegistration.AddAppOwners(ctx, appId, ownersToAdd)
			if err != nil {
				return err
			}
		}

		ownersToRemove := findMissingOwners(managedOwners, owners)
		if len(ownersToRemove) > 0 {
			err := graphClient.AppRegistration.RemoveAppOwners(ctx, appId, ownersToRemove)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findMissingOwners(desired, current []string) []string {
	missing := []string{}
	for _, owner := range desired {
		if !contains(current, owner) {
			missing = append(missing, owner)
		}
	}
	return missing
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
