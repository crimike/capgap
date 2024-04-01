package adgraph

import "capgap/enums"

//ApplicationList
type ApplicationList struct {
	Applications []string               `json:"Applications"`
	Acrs         []enums.ApplicationACR `json:"Acrs"`
}

type ConditionalAccessApplications struct {
	IncludedApplications []ApplicationList `json:"Include"`
	ExcludedApplications []ApplicationList `json:"Exclude"`
}

type ExternalUsers struct {
	GuestOrExternalUserTypes enums.GuestOrExternalUserType `json:"GuestOrExternalUserTypes"`
}

//Users
type Users struct {
	Users                 []string      `json:"Users"`
	Groups                []string      `json:"Groups"`
	Roles                 []string      `json:"Roles"`
	GuestsOrExternalUsers ExternalUsers `json:"GuestsOrExternalUsers"`
}

type ConditionalAccessUsers struct {
	IncludedUsers []Users `json:"Include"`
	ExcludedUsers []Users `json:"Exclude"`
}

//Locations
type Locations struct {
	Locations []string `json:"Locations"`
}

type ConditionalAccessLocations struct {
	IncludedLocations []Locations `json:"Include"`
	ExcludedLocations []Locations `json:"Exclude"`
}

type ClientTypes struct {
	ClientTypes []enums.ClientType `json:"ClientTypes"`
}

type ConditionalAccessClientTypes struct {
	IncludedClientTypes []ClientTypes `json:"Include"`
}

type DevicePlatforms struct {
	DevicePlatforms []enums.DevicePlatform `json:"DevicePlatforms"`
}

type ConditionalAccessDevicePlatforms struct {
	IncludedDevicePlatforms []DevicePlatforms `json:"Include"`
	ExcludedDevicePlatforms []DevicePlatforms `json:"Exclude"`
}

type Devices struct {
	DeviceRule string `json:"DeviceRule"`
}

type ConditionalAccessDevices struct {
	IncludedDevices []Devices `json:"Include"`
	ExcludedDevices []Devices `json:"Exclude"`
}

type PolicyCondition struct {
	ApplicationList    ConditionalAccessApplications    `json:"Applications"`
	LocationList       ConditionalAccessLocations       `json:"Locations"`
	UserList           ConditionalAccessUsers           `json:"Users"`
	ClientTypeList     ConditionalAccessClientTypes     `json:"ClientTypes"`
	DevicePlatformList ConditionalAccessDevicePlatforms `json:"DevicePlatforms"`
	DeviceList         ConditionalAccessDevices         `json:"Devices"`
}

type PolicyControls struct {
	Controls        []enums.ConditionalAccessControl `json:"Control"`
	AuthStrengthIds []string                         `json:"AuthStrengthIds"`
}

type PolicyDetail struct {
	State                                     enums.ConditionalAccessPolicyState `json:"State"`
	Conditions                                PolicyCondition                    `json:"Conditions"`
	Controls                                  []PolicyControls                   `json:"Controls"`
	SessionControls                           []string                           `json:"SessionControls"`
	EnforceAllPoliciesForEas                  bool                               `json:"EnforceAllPoliciesForEas"`
	IncludeOtherLegacyClientTypeForEvaluation bool                               `json:"IncludeOtherLegacyClientTypeForEvaluation"`
	PersistentBrowserSessionMode              string                             `json:"PersistentBrowserSessionMode"`
	CasSessionControlType                     int                                `json:"CasSessionControlType"`
	ContinuousAccessEvaluationMode            string                             `json:"ContinuousAccessEvaluationMode"`
}

type ConditionalAccessPolicy struct {
	DisplayName  string   `json:"displayName"`
	ObjectId     string   `json:"objectId"`
	PolicyDetail []string `json:"policyDetail"`
}

type ConditionalAccessPolicyResponse struct {
	ConditionalAccessPolicies []ConditionalAccessPolicy `json:"value"`
	NextLink                  string                    `json:"odata.nextLink,omitempty"`
}

type KnownNetworkPolicy struct {
	NetworkName  string   `json:"NetworkName"`
	NetworkId    string   `json:"NetworkId"`
	CidrIpRanges []string `json:"CidrIpRanges"`
	Categories   []string `json:"Categories"`
}

type LocationPolicyDetail struct {
	CompressedCidrIpRanges string             `json:"CompressedCidrIpRanges"`
	Categories             []string           `json:"Categories"`
	KnownNetworkPolicies   KnownNetworkPolicy `json:"KnownNetworkPolicies"`
}

type LocationPolicy struct {
	DisplayName  string   `json:"displayName"`
	ObjectId     string   `json:"objectId"`
	PolicyId     string   `json:"policyIdentifier"`
	PolicyDetail []string `json:"policyDetail"`
}

type LocationPolicyResponse struct {
	LocationPolicies []LocationPolicy `json:"value"`
	NextLink         string           `json:"odata.nextLink,omitempty"`
}
