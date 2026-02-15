package appregistration

import (
	"context"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/log"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

func (s *Service) GetAppOwners(ctx context.Context, objectID string) ([]string, error) {
	response, err := s.sdk.Applications().ByApplicationId(objectID).Owners().Get(ctx, nil)
	if err != nil {
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

func (s *Service) AddAppOwners(ctx context.Context, appID string, owners []string) ([]string, error) {
	logger := log.FromContext(ctx)
	request := graphmodels.NewReferenceCreate()
	for _, owner := range owners {
		odataID := "https://graph.microsoft.com/v1.0/directoryObjects/" + owner
		request.SetOdataId(&odataID)
		err := s.sdk.Applications().ByApplicationId(appID).Owners().Ref().Post(ctx, request, nil)
		if err != nil {
			if odataError, ok := err.(*odataerrors.ODataError); ok {
				logger.Error(err, "failed to add owner to application", "appID", appID, "ownerID", owner)
				logger.Error(err, "OData error code", "code", *odataError.GetErrorEscaped().GetMessage())
				// If the object is already added as an owner by IDM. the error will be skipped
				if odataError.GetErrorEscaped().GetMessage() != nil && strings.Contains(*odataError.GetErrorEscaped().GetMessage(), "already exist") {
					logger.Info("Owner already exists for the application, skipping addition", "appID", appID, "ownerID", owner)
					continue
				}
			}
			return nil, err
		}
	}
	return owners, nil
}

func (s *Service) RemoveAppOwners(ctx context.Context, appID string, owners []string) error {
	for _, owner := range owners {
		err := s.sdk.Applications().ByApplicationId(appID).Owners().ByDirectoryObjectId(owner).Ref().Delete(ctx, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
