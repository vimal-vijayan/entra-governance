package client

import (
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)


type GraphClient struct {
	sdk *msgraphsdk.GraphServiceClient
}

func NewGraphClient(sdk *msgraphsdk.GraphServiceClient) *GraphClient {
	return &GraphClient{sdk: sdk}
}
