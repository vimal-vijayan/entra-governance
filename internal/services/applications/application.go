package applications

import (
	"context"
	"fmt"
	"sort"

	entrav1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/client"
	graphappregistration "github.com/vimal-vijayan/entra-governance/internal/graph/appregistration"
)

type Service struct {
	factory *client.ClientFactory
}

func NewService(factory *client.ClientFactory) *Service {
	return &Service{factory: factory}
}

func (s *Service) getGraphClient(ctx context.Context, entraApp entrav1alpha1.EntraAppRegistration) (*client.GraphClient, error) {
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

func (s *Service) Create(ctx context.Context, entraApp entrav1alpha1.EntraAppRegistration) (string, string, error) {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return "", "", err
	}

	response, err := graphClient.AppRegistration.Create(ctx, graphappregistration.CreateRequest{
		DisplayName:               entraApp.Spec.Name,
		Description:               entraApp.Spec.Description,
		Tags:                      append([]string(nil), entraApp.Spec.Tags...),
		SignInAudience:            entraApp.Spec.SignInAudience,
		SamlMetadataURL:           entraApp.Spec.SamlMetadataUrl,
		IsFallbackPublicClient:    entraApp.Spec.IsFallbackPublicClient,
		IsDeviceOnlyAuthSupported: entraApp.Spec.IsDeviceOnlyAuthSupported,
		GroupMembershipClaims:     entraApp.Spec.GroupMembershipClaims,
	})
	if err != nil {
		return "", "", err
	}

	return response.AppClientID, response.AppObjectID, nil
}

func (s *Service) Delete(ctx context.Context, appID string, entraApp entrav1alpha1.EntraAppRegistration) error {
	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return err
	}

	return graphClient.AppRegistration.Delete(ctx, appID)
}

func (s *Service) Update(ctx context.Context, entraApp entrav1alpha1.EntraAppRegistration) error {
	_, err := s.GetAndPatch(ctx, entraApp)
	return err
}

type desiredApplication struct {
	ObjectID                  string
	DisplayName               string
	Description               string
	Tags                      []string
	SignInAudience            string
	SamlMetadataURL           string
	IsFallbackPublicClient    bool
	IsDeviceOnlyAuthSupported bool
	GroupMembershipClaims     string
}

func (s *Service) GetAndPatch(ctx context.Context, entraApp entrav1alpha1.EntraAppRegistration) (bool, error) {

	if entraApp.Status.AppRegistrationObjID == "" {
		return false, fmt.Errorf("appRegistrationObjID is empty in status")
	}

	graphClient, err := s.getGraphClient(ctx, entraApp)
	if err != nil {
		return false, err
	}

	resp, err := graphClient.AppRegistration.Get(ctx, entraApp.Status.AppRegistrationObjID)

	if err != nil {
		return false, err
	}

	desired := desiredApplication{
		ObjectID:                  entraApp.Status.AppRegistrationObjID,
		DisplayName:               entraApp.Spec.Name,
		Description:               entraApp.Spec.Description,
		Tags:                      append([]string(nil), entraApp.Spec.Tags...),
		SignInAudience:            entraApp.Spec.SignInAudience,
		SamlMetadataURL:           entraApp.Spec.SamlMetadataUrl,
		IsFallbackPublicClient:    entraApp.Spec.IsFallbackPublicClient,
		IsDeviceOnlyAuthSupported: entraApp.Spec.IsDeviceOnlyAuthSupported,
		GroupMembershipClaims:     entraApp.Spec.GroupMembershipClaims,
	}

	patchRequest, hasChanges := buildPatchRequest(resp, desired)
	if !hasChanges {
		return false, nil
	}

	if err := graphClient.AppRegistration.Patch(ctx, patchRequest); err != nil {
		return false, err
	}

	return true, nil
}

func buildPatchRequest(current *graphappregistration.Application, desired desiredApplication) (graphappregistration.PatchRequest, bool) {
	request := graphappregistration.PatchRequest{
		ObjectID: desired.ObjectID,
	}
	hasChanges := false

	if current.DisplayName != desired.DisplayName {
		value := desired.DisplayName
		request.DisplayName = &value
		hasChanges = true
	}
	if current.Description != desired.Description {
		value := desired.Description
		request.Description = &value
		hasChanges = true
	}
	if !equalStringSets(current.Tags, desired.Tags) {
		value := append([]string(nil), desired.Tags...)
		request.Tags = &value
		hasChanges = true
	}
	if current.SignInAudience != desired.SignInAudience {
		value := desired.SignInAudience
		request.SignInAudience = &value
		hasChanges = true
	}
	if current.SamlMetadataURL != desired.SamlMetadataURL {
		value := desired.SamlMetadataURL
		request.SamlMetadataURL = &value
		hasChanges = true
	}
	if current.IsFallbackPublicClient != desired.IsFallbackPublicClient {
		value := desired.IsFallbackPublicClient
		request.IsFallbackPublicClient = &value
		hasChanges = true
	}
	if current.IsDeviceOnlyAuthSupported != desired.IsDeviceOnlyAuthSupported {
		value := desired.IsDeviceOnlyAuthSupported
		request.IsDeviceOnlyAuthSupported = &value
		hasChanges = true
	}
	if current.GroupMembershipClaims != desired.GroupMembershipClaims {
		value := desired.GroupMembershipClaims
		request.GroupMembershipClaims = &value
		hasChanges = true
	}

	return request, hasChanges
}

func equalStringSets(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	normalizedLeft := append([]string(nil), left...)
	normalizedRight := append([]string(nil), right...)
	sort.Strings(normalizedLeft)
	sort.Strings(normalizedRight)

	for i := range normalizedLeft {
		if normalizedLeft[i] != normalizedRight[i] {
			return false
		}
	}

	return true
}
