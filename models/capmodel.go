package models

import "capgap/enums"

type ConditionalAccessControl struct {
	EnforcedControl []string // each entry is OR-ed together
}

type DynamicGroup struct {
	DisplayName    string
	MembershipRule string
	ObjectId       string
}

type ConditionalAccessPolicy struct {
	DisplayName             string                             `json:"DisplayName"`
	ObjectId                string                             `json:"ObjectId"`
	State                   enums.ConditionalAccessPolicyState `json:"State"`
	IncludedApplications    []string                           `json:"IncludedApplications"`
	ExcludedApplications    []string                           `json:"ExcludedApplications"`
	IncludedActions         []enums.ApplicationACR             `json:"IncludedActions"`
	ApplicationMode         enums.ApplicationMode              `json:"ApplicationMode"`
	IncludedDynamicGroups   []DynamicGroup                     `json:"IncludedDynamicGroups"`
	ExcludedDynamicGroups   []DynamicGroup                     `json:"ExcludedDynamicGroups"`
	IncludedUsers           []string                           `json:"IncludedUsers"`
	ExcludedUsers           []string                           `json:"ExcludedUsers"`
	IncludedLocations       []string                           `json:"IncludedLocations"`
	ExcludedLocations       []string                           `json:"ExcludedLocations"`
	IncludedClientTypes     []enums.ClientType                 `json:"IncludedClientTypes"`
	IncludedDevicePlatforms map[enums.DevicePlatform]struct{}  `json:"IncludedDevicePlatforms"`
	IncludedDevices         string                             `json:"IncludedDevices"`
	ExcludedDevices         string                             `json:"ExcludedDevices"`
	Controls                []ConditionalAccessControl         `json:"Controls"` // each entry is AND-ed together
}
