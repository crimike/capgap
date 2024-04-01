package capgap

import (
	"capgap/enums"
	"capgap/models"
	"slices"
)

func CapAppliesToUser(cap models.ConditionalAccessPolicy, userId string) bool {
	if len(cap.IncludedUsers) == 0 {
		return false
	}
	if (cap.IncludedUsers[0] == enums.AllUsers || slices.Contains(cap.IncludedUsers, userId)) && (!slices.Contains(cap.ExcludedUsers, userId)) {
		return true
	}
	return false
}

func CapAppliesToApplication(cap models.ConditionalAccessPolicy, appId string) bool {

	if (slices.Contains(cap.IncludedApplications, enums.AllApplications) || slices.Contains(cap.IncludedApplications, appId)) && !slices.Contains(cap.ExcludedApplications, appId) {
		return true
	}
	return false
}

func CapAppliesToLocation(cap models.ConditionalAccessPolicy, location models.Location) bool {

	if len(cap.IncludedLocations) == 0 {
		return true
	}
	if slices.Contains(cap.IncludedLocations, location.PolicyId) && !slices.Contains(cap.ExcludedLocations, location.PolicyId) {
		return true
	}
	return false
}

func CapAppliesToClientType(cap models.ConditionalAccessPolicy, clientType enums.ClientType) bool {
	if len(cap.IncludedClientTypes) == 0 || slices.Contains(cap.IncludedClientTypes, clientType) {
		return true
	}
	return false
}

func CapAppliesToDevicePlatform(cap models.ConditionalAccessPolicy, devicePlatform enums.DevicePlatform) bool {
	if len(cap.IncludedDevicePlatforms) == 0 {
		return true
	}
	_, ok := cap.IncludedDevicePlatforms[devicePlatform]
	if ok {
		return true
	}
	return false
}
