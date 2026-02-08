package appregistration

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Delete(ctx context.Context, appID string) error {
	logger := log.FromContext(ctx)
	err := s.sdk.Applications().ByApplicationId(appID).Delete(ctx, nil)
	if err != nil {
		logger.Error(err, "failed to delete application", "applicationID", appID)
		return err
	}

	logger.Info("application deleted successfully", "applicationID", appID)
	return nil
}
