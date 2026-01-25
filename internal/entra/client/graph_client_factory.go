package client

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ClientFactory struct {
	k8s client.Client
}

type SecretRef struct {
	Name      string
	Namespace string
}

type ServiceAccountRef struct {
	Name      string
	Namespace string
}

type GraphClientInterface interface {
	ForClientSecret(ctx context.Context, ref SecretRef) (*msgraphsdk.GraphServiceClient, error)
	ForWorkloadIdentity(ctx context.Context, ref ServiceAccountRef) (*msgraphsdk.GraphServiceClient, error)
}

func (cf *ClientFactory) ForClientSecret(ctx context.Context, ref SecretRef) (*msgraphsdk.GraphServiceClient, error) {
	logger := log.FromContext(ctx)

	if ref.Namespace == "" && ref.Name == "" {
		logger.Error(fmt.Errorf("invalid secret reference: namespace and name cannot both be empty"), "ref", ref)
		return nil, fmt.Errorf("invalid secret reference: namespace and name cannot both be empty")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: ref.Namespace, Name: ref.Name}
	if err := cf.k8s.Get(ctx, key, secret); err != nil {
		logger.Error(err, "failed to get secret", "secret", ref.Name, "namespace", ref.Namespace)
		return nil, err
	}

	tenantId, err := getSecretData(secret, "tenantId")
	if err != nil {
		logger.Error(err, "failed to get tenantId from secret", "secret", ref.Name, "namespace", ref.Namespace)
		return nil, err
	}
	clientId, err := getSecretData(secret, "clientId")
	if err != nil {
		logger.Error(err, "failed to get clientId from secret", "secret", ref.Name, "namespace", ref.Namespace)
		return nil, err
	}
	clientSecret, err := getSecretData(secret, "clientSecret")
	if err != nil {
		logger.Error(err, "failed to get clientSecret from secret", "secret", ref.Name, "namespace", ref.Namespace)
		return nil, err
	}

	logger.Info("Successfully retrieved client credentials from secret", "secret", ref.Name, "namespace", ref.Namespace)
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		logger.Error(err, "failed to create client credentials", "secret", ref.Name, "namespace", ref.Namespace)
		return nil, err
	}

	return cf.setupGraphClient(cred)

}

func (cf *ClientFactory) ForWorkloadIdentity(ctx context.Context, ref ServiceAccountRef) (*msgraphsdk.GraphServiceClient, error) {
	//TODO: Implement workload identity

	return nil, fmt.Errorf("workload identity is not implemented yet")
}

func (cf *ClientFactory) setupGraphClient(cred azcore.TokenCredential) (*msgraphsdk.GraphServiceClient, error) {
	scope := []string{"https://graph.microsoft.com/.default"}
	sdk, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scope)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func getSecretData(secret *corev1.Secret, key string) (string, error) {

	if secret == nil {
		return "", fmt.Errorf("secret is nil")
	}

	value, exists := secret.Data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in secret %s/%s", key, secret.Namespace, secret.Name)
	}

	return string(value), nil
}

func NewClientFactory(k8s client.Client) *ClientFactory {
	return &ClientFactory{k8s: k8s}
}
