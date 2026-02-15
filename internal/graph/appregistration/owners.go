package appregistration

import (
	"context"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/log"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func (s *Service) GetAppOwners(ctx context.Context, objectID string) ([]string, error) {
	logger := log.FromContext(ctx).WithValues("component", "graph-client", "op", "GetAppOwners", "AppId", objectID)

	logger.V(1).Info("Getting application owners from Microsoft Graph")
	response, err := s.sdk.Applications().ByApplicationId(objectID).Owners().Get(ctx, nil)
	if err != nil {
		logger.Error(err, "failed to get application owners", "objectID", objectID)
		return nil, err
	}

	var owners []string
	if response != nil && response.GetValue() != nil {
		for _, owner := range response.GetValue() {
			if owner.GetId() != nil {
				owners = append(owners, *owner.GetId())
			}
		}
	}

	return owners, nil
}

func (s *Service) AddAppOwners(ctx context.Context, appID string, owners []string) error {
	logger := log.FromContext(ctx).WithValues("component", "graph-client", "op", "AddAppOwners", "appID", appID)
	request := graphmodels.NewReferenceCreate()
	for _, owner := range owners {
		odataID := "https://graph.microsoft.com/v1.0/directoryObjects/" + owner
		request.SetOdataId(&odataID)
		err := s.sdk.Applications().ByApplicationId(appID).Owners().Ref().Post(ctx, request, nil)
		if err != nil {
			if odataError, ok := err.(*odataerrors.ODataError); ok {
				logger.V(1).Error(err, "OData error code", "code", *odataError.GetErrorEscaped().GetMessage())
				// If the object is already added as an owner by IDM. the error will be skipped
				if odataError.GetErrorEscaped().GetMessage() != nil && strings.Contains(*odataError.GetErrorEscaped().GetMessage(), "already exist") {
					logger.V(1).Info("Owner already exists for the application, skipping addition", "appID", appID, "ownerID", owner)
					continue
				}
			}
			return err
		}
	}
	return nil
}

func (s *Service) RemoveAppOwners(ctx context.Context, appID string, owners []string) error {
	logger := log.FromContext(ctx).WithValues("component", "graph-client", "op", "RemoveAppOwners", "appID", appID)
	for _, owner := range owners {
		err := s.sdk.Applications().ByApplicationId(appID).Owners().ByDirectoryObjectId(owner).Ref().Delete(ctx, nil)
		if err != nil {
			logger.V(1).Error(err, "failed to remove owner from application", "appID", appID, "ownerID", owner)
			return err
		}
	}
	return nil
}
