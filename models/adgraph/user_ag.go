package adgraph

type User struct {
	DisplayName       string `json:"displayName"`
	EmailAddress      string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	ObjectId          string `json:"objectId"`
}

type UserResponse struct {
	Users    []User `json:"value"`
	NextLink string `json:"odata.nextLink,omitempty"`
}
