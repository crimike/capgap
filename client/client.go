package client

import (
	"capgap/settings"
	"net/http"
)

type AzureClient struct {
	AccessToken string
	MainUrl     string
	Tenant      string
	ApiVersion  string
	HttpClient  *http.Client
}

func (c *AzureClient) InitializeClient() {
	if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
		c.InitializeAzureADGraphClient()
	}
	//TODO add MS graph
}
