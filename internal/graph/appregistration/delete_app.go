package appregistration

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Delete(ctx context.Context, objectID string) error {
	logger := log.FromContext(ctx)

	if err := s.sdk.Applications().ByApplicationId(objectID).Delete(ctx, nil); err != nil {
		logger.Error(err, "failed to delete application", "objectID", objectID)
		return err
	}

	logger.Info("application deleted successfully", "objectID", objectID)
	return nil
}
