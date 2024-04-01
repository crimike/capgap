package adgraph

type RoleAssignment struct {
	Id          string `json:"id"`
	PrincipalId string `json:"principalId"`
}

type RoleAssignmentResponse struct {
	RoleAssignments []RoleAssignment `json:"value"`
	NextLink        string           `json:"odata.nextLink,omitempty"`
}
