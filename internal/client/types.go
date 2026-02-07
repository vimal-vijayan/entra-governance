package client

import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

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
