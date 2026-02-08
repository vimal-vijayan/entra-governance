package groups

import (
	"context"
	"fmt"
)

func (s *Service) Delete(ctx context.Context, groupID string) error {

	if err := s.sdk.Groups().ByGroupId(groupID).Delete(ctx, nil); err != nil {
		return fmt.Errorf("failed to delete group by ID: %v", err)
	}

	return nil
}
