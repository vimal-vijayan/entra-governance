package client

import (
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/vimal-vijayan/entra-governance/internal/graph/appregistration"
	"github.com/vimal-vijayan/entra-governance/internal/graph/groups"
	"github.com/vimal-vijayan/entra-governance/internal/graph/serviceprincipal"
)

type GraphClient struct {
	sdk               *msgraphsdk.GraphServiceClient
	Groups            groups.API
	AppRegistration   appregistration.API
	ServicePrincipals serviceprincipal.API
}

func NewGraphClient(sdk *msgraphsdk.GraphServiceClient) *GraphClient {
	return &GraphClient{
		sdk:               sdk,
		Groups:            groups.NewAPI(sdk),
		AppRegistration:   appregistration.NewAPI(sdk),
		ServicePrincipals: serviceprincipal.NewAPI(sdk),
	}
}
