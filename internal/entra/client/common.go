package client

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (c *GraphClient) getUsers(ctx context.Context, user string) (string, error) {
	logger := log.FromContext(ctx)
	resp, err := c.sdk.Users().ByUserId(user).Get(ctx, nil)

	if err != nil {
		logger.Error(err, "failed to get user", "user", user)
		return "", fmt.Errorf("failed to get user %w", err)
	}

	logger.Info("successfully fetched user", "userID", *resp.GetId())
	return *resp.GetId(), nil
}
