package models

type Location struct {
	DisplayName string
	//commenting out the objectID since the PolicyID is the one being used
	// ObjectId    string
	PolicyId  string
	IpRange   []string
	IsTrusted bool
}
