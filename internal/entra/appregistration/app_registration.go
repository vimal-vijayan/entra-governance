package appregistration

// Code snippets are only available for the latest major version. Current major version is $v1.*

// Dependencies
import (
	"context"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	//other-imports
)

func CreateApplication(graphClient *msgraphsdk.GraphServiceClient) (*graphmodels.Application, error) {
	requestBody := graphmodels.NewApplication()
	displayName := "Display name"
	requestBody.SetDisplayName(&displayName)

	// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
	// applications, err := graphClient.Applications().Post(context.Background(), requestBody, nil)
	_, err := graphClient.Applications().Post(context.Background(), requestBody, nil)
	if err != nil {
		return nil, err
	}
	// return applications, nil
	return &graphmodels.Application{}, nil
}
