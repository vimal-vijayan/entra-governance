package groups

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Get(ctx context.Context, groupID string) (*GroupGetResponse, error) {
	logger := log.FromContext(ctx)

	var HttpStatusCode string

	if groupID == "" {
		logger.Error(fmt.Errorf("groupID cannot be empty"), "invalid groupID")
		return nil, fmt.Errorf("group id is empty")
	}

	resp, err := s.sdk.Groups().ByGroupId(groupID).Get(ctx, nil)
	if err != nil {
		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			HttpStatusCode = fmt.Sprintf("%d", odataErr.GetStatusCode())
			logger.Error(err, "failed to get group", "groupID", groupID, "statusCode", odataErr.GetStatusCode())
		} else {
			logger.Error(err, "failed to get group", "groupID", groupID)
		}
		return &GroupGetResponse{
			ID:             "",
			DisplayName:    "",
			HttpStatusCode: HttpStatusCode,
		}, fmt.Errorf("failed to get group %w", err)
	}

	logger.Info("successfully fetched group", "groupID", *resp.GetId())
	return &GroupGetResponse{
		ID:             *resp.GetId(),
		DisplayName:    *resp.GetDisplayName(),
		HttpStatusCode: "200",
	}, nil
}
