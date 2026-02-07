package graph

import (
	"fmt"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphClient struct {
	sdk *msgraphsdk.GraphServiceClient
}

func NewGraphClient(sdk *msgraphsdk.GraphServiceClient) *GraphClient {
	return &GraphClient{sdk: sdk}
}

func (c *GraphClient) ensureClient() error {
	if c == nil || c.sdk == nil {
		return fmt.Errorf("Graph client is not initialized")
	}
	return nil
}
