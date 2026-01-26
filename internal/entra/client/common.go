package client

import (
	"context"
	"fmt"

	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (c *GraphClient) getUsers(ctx context.Context, user string) (string, error) {
	logger := log.FromContext(ctx)
	resp, err := c.sdk.Users().ByUserId(user).Get(ctx, nil)

	if err != nil {
		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			logger.Error(err, "failed to get user", "user", user, "statusCode", odataErr.GetStatusCode())
		}
		logger.Error(err, "failed to get user", "user", user)
		return "", fmt.Errorf("failed to get user %w", err)
	}

	logger.Info("successfully fetched user", "userID", *resp.GetId())
	return *resp.GetId(), nil
}

func (c *GraphClient) getGroups(ctx context.Context, group string) (string, error) {
	logger := log.FromContext(ctx)
	resp, err := c.sdk.Groups().ByGroupId(group).Get(ctx, nil)

	if err != nil {
		if odataErr, ok := err.(*odataerrors.ODataError); ok {
			logger.Error(err, "failed to get group", "group", group, "statusCode", odataErr.GetStatusCode())
		}
		logger.Error(err, "failed to get group", "group", group)
		return "", err
	}

	logger.Info("successfully fetched group", "groupID", *resp.GetId())
	return *resp.GetId(), nil
}
