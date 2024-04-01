package adgraph

type Group struct {
	DisplayName    string   `json:"displayName"`
	ObjectId       string   `json:"objectId"`
	GroupType      []string `json:"groupTypes"`
	MembershipRule string   `json:"membershipRule"`
}

type GroupResponse struct {
	Groups   []Group `json:"value"`
	NextLink string  `json:"odata.nextLink,omitempty"`
}
