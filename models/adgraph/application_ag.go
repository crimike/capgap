package adgraph

type Application struct {
	DisplayName   string `json:"displayName"`
	ObjectId      string `json:"objectId"`
	ApplicationId string `json:"appId"`
}

type ApplicationResponse struct {
	Applications []Application `json:"value"`
	NextLink     string        `json:"odata.nextLink,omitempty"`
}
