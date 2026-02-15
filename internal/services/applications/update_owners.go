package applications

import (
	"context"

	entrav1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// UpdateOwners manages the owners of an application registration
func (s *Service) UpdateOwners(ctx context.Context, appId string, entraApp entrav1alpha1.EntraAppRegistration) error {

	logger := log.FromContext(ctx).WithValues("component", "application-service", "op", "UpdateOwners", "appID", appId)
	logger.Info("Updating application owners")

	owners := append([]string(nil), (*entraApp.Spec.Owners)...)
	managedOwners := append([]string(nil), entraApp.Status.Owners...)

	if len(owners) == 0 && len(managedOwners) == 0 {
		return nil
	}

	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		logger.Error(err, "failed to get graph client")
		return err
	}

	if !equalStringSets(owners, managedOwners) {
		ownersToAdd := findMissing(owners, managedOwners)

		if err := transformHelper(ctx, appId, ownersToAdd, graphClient.AppRegistration.AddAppOwners); err != nil {
			logger.Error(err, "failed to add application owners", "ownersToAdd", ownersToAdd)
			return err
		}

		ownersToRemove := findMissing(managedOwners, owners)

		if err := transformHelper(ctx, appId, ownersToRemove, graphClient.AppRegistration.RemoveAppOwners); err != nil {
			logger.Error(err, "failed to remove application owners", "ownersToRemove", ownersToRemove)
			return err
		}

	} else {
		currentOwners, err := graphClient.AppRegistration.GetAppOwners(ctx, appId)

		if err != nil {
			logger.Error(err, "failed to get current application owners")
			return err
		}

		ownersToAdd := findMissing(owners, currentOwners)
		if err := transformHelper(ctx, appId, ownersToAdd, graphClient.AppRegistration.AddAppOwners); err != nil {
			logger.Error(err, "failed to add application owners", "ownersToAdd", ownersToAdd)
			return err
		}

		ownersToRemove := findMissing(managedOwners, owners)

		if err := transformHelper(ctx, appId, ownersToRemove, graphClient.AppRegistration.RemoveAppOwners); err != nil {
			logger.Error(err, "failed to remove application owners", "ownersToRemove", ownersToRemove)
			return err
		}
	}

	return nil
}
