package serviceprincipal

import (
	"context"
	"fmt"

	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (s *Service) Create(ctx context.Context, req CreateRequest) (*ServicePrincipalCreateResponse, error) {
	logger := log.FromContext(ctx)

	requestBody := graphmodels.NewServicePrincipal()
	appId := req.AppID
	displayName := req.DisplayName
	requestBody.SetDisplayName(&displayName)
	if req.DisableVisibilityForGuests {
		requestBody.SetTags([]string{"HideApp"})
	}
	requestBody.SetAppId(&appId)
	response, err := s.sdk.ServicePrincipals().Post(ctx, requestBody, nil)

	if err != nil {
		logger.Error(err, "failed to create service principal", "applicationID", appId)
		return nil, err
	}

	// Check if response or ID is nil before dereferencing
	if response == nil {
		logger.Error(nil, "service principal response is nil", "applicationID", appId)
		return nil, fmt.Errorf("service principal response is nil")
	}

	servicePrincipalID := response.GetId()
	if servicePrincipalID == nil {
		logger.Error(nil, "service principal ID is nil", "applicationID", appId)
		return nil, fmt.Errorf("service principal ID is nil in response")
	}

	logger.Info("service principal created successfully", "applicationID", appId, "servicePrincipalID", *servicePrincipalID)

	return &ServicePrincipalCreateResponse{
		ServicePrincipalID: *servicePrincipalID,
	}, nil
}
