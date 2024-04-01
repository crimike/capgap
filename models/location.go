package models

type Location struct {
	DisplayName string
	ObjectId    string
	PolicyId    string
	IpRange     []string
	IsTrusted   bool
}
