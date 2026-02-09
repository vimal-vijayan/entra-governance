package applications

import (
	"context"
	"fmt"

	appregistration "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/client"
)

type Service struct {
	factory *client.ClientFactory
}

func NewService(factory *client.ClientFactory) *Service {
	return &Service{factory: factory}
}

func (s *Service) getGraphClient(ctx context.Context, entraApp appregistration.EntraAppRegistration) (*client.GraphClient, error) {
	if entraApp.Spec.ForProvider == nil {
		return nil, fmt.Errorf("forProvider spec is nil")
	}

	if entraApp.Spec.ForProvider.CredentialSecretRef == "" {
		return nil, fmt.Errorf("credential secret reference is empty in forProvider spec")
	}

	secretRef := client.SecretRef{
		Name:      entraApp.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraApp.Namespace,
	}

	sdk, err := s.factory.ForClientSecret(ctx, secretRef)
	if err != nil {
		return nil, err
	}

	return client.NewGraphClient(sdk), nil
}

func (s *Service) Create(ctx context.Context, entraApp appregistration.EntraAppRegistration) (string, string, error) {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return "", "", err
	}

	response, err := graphClient.AppRegistration.Create(ctx, entraApp.Spec)
	if err != nil {
		return "", "", err
	}
	return response.AppClientID, response.AppObjectID, nil
}

func (s *Service) Delete(ctx context.Context, appID string, entraApp appregistration.EntraAppRegistration) error {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return err
	}

	return graphClient.AppRegistration.Delete(ctx, appID)
}

type servicePrincipalCreateResponse struct {
	ServicePrincipalID string
}

type appRegistrationResponse struct {
	AppClientID string
	AppObjectID string
}

func (s *Service) Update(ctx context.Context, entraApp appregistration.EntraAppRegistration) error {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return err
	}

	return graphClient.AppRegistration.Update(ctx, entraApp)
}

type desiredApplication struct {
	AppId           string
	AppObjectID     string
	Tags            []string
	IdentifierURI   []string
	WebRedirectURIs []string
}

func (s *Service) GetAndPatch(ctx context.Context, entraApp appregistration.EntraAppRegistration) (*desiredApplication, error) {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	
	if err != nil {
		return nil, err
	}

	resp, err := graphClient.AppRegistration.Get(ctx, entraApp.Status.AppRegistrationObjID)

	if err != nil {
		return nil, err
	}

	if entraApp.Spec.Tags != nil {
		tags := compareTags(resp.GetTags(), entraApp.Spec.Tags)
		resp.SetTags(tags)
	}
	

	return &desiredApplication{
		AppId:           *resp.GetAppId(),
		AppObjectID:     *resp.GetId(),
		Tags:            resp.GetTags(),
		IdentifierURI:   resp.GetIdentifierUris(),
		WebRedirectURIs: resp.GetWeb().GetRedirectUris(),
	}, nil

}

func compareTags(existingTags, desiredTags []string) []string {
	for _, existingTag := range existingTags {
		found := false
		for _, desiredTag := range desiredTags {
			if existingTag == desiredTag {
				found = true
				break
			}
		}
		if !found {
			desiredTags = append(desiredTags, existingTag)
		}
	}
	return desiredTags
}