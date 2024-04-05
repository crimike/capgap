package parsers

import (
	"bytes"
	"capgap/client"
	"capgap/enums"
	"capgap/models"
	"capgap/models/adgraph"
	"capgap/settings"
	"compress/flate"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"slices"
)

func ParseApplicationsADGraph(c *client.AzureClient) ([]models.Application, error) {
	var response []models.Application

	appList, err := c.GetApplicationsAdGraph()
	if err != nil {
		settings.ErrorLogger.Println(err)
		return response, err
	}

	for _, appEntry := range appList {
		var app models.Application
		app.ApplicationId = appEntry.ApplicationId
		app.DisplayName = appEntry.DisplayName
		app.ObjectId = appEntry.ObjectId

		response = append(response, app)
	}

	return response, nil
}

func ParseUsersADGraph(c *client.AzureClient) ([]models.User, error) {
	var response []models.User

	userList, err := c.GetUsersAdGraph()
	if err != nil {
		settings.ErrorLogger.Println(err)
		return response, err
	}

	for _, userEntry := range userList {
		var user models.User
		user.DisplayName = userEntry.DisplayName
		user.EmailAddress = userEntry.EmailAddress
		user.ObjectId = userEntry.ObjectId
		user.UserPrincipalName = userEntry.UserPrincipalName

		response = append(response, user)
	}

	return response, nil
}

func ParseLocationsADGraph(c *client.AzureClient) ([]models.Location, error) {

	var response []models.Location
	namedLocations, err := c.GetNamedLocationsAdGraph()
	if err != nil {
		settings.ErrorLogger.Println(err)
		return response, err
	}

	for _, namedLocation := range namedLocations {
		if namedLocation.DisplayName == "Known Networks List" {
			for _, policyDetail := range namedLocation.PolicyDetail {

				var (
					loc       models.Location
					polDetail adgraph.LocationPolicyDetail
				)
				err := json.Unmarshal([]byte(policyDetail), &polDetail)
				if err != nil {
					settings.ErrorLogger.Println(err)
					return response, err
				}
				loc.DisplayName = polDetail.KnownNetworkPolicies.NetworkName
				loc.ObjectId = polDetail.KnownNetworkPolicies.NetworkId
				loc.PolicyId = polDetail.KnownNetworkPolicies.NetworkId
				loc.IpRange = polDetail.KnownNetworkPolicies.CidrIpRanges
				if slices.Contains(polDetail.KnownNetworkPolicies.Categories, "trusted") {
					loc.IsTrusted = true
				} else {
					loc.IsTrusted = false
				}
				if loc.DisplayName != "" {
					response = append(response, loc)
				}
			}
		} else {
			var (
				loc       models.Location
				polDetail adgraph.LocationPolicyDetail
			)

			err := json.Unmarshal([]byte(namedLocation.PolicyDetail[0]), &polDetail)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return response, err
			}

			loc.DisplayName = namedLocation.DisplayName
			loc.ObjectId = namedLocation.ObjectId
			loc.PolicyId = namedLocation.PolicyId
			b64decoded, err := base64.StdEncoding.DecodeString(polDetail.CompressedCidrIpRanges)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return response, err
			}
			b := bytes.NewReader(b64decoded)

			z := flate.NewReader(b)
			defer z.Close()
			p, err := io.ReadAll(z)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return response, err
			}
			loc.IpRange = append(loc.IpRange, string(p))
			if slices.Contains(polDetail.Categories, "trusted") {
				loc.IsTrusted = true
			} else {
				loc.IsTrusted = false
			}
			if loc.DisplayName != "" {
				response = append(response, loc)
			}
		}

	}

	return response, nil
}

func ParseGroupADGraph(c *client.AzureClient, groupId string) (models.DynamicGroup, []string, error) {

	var (
		dgroup models.DynamicGroup
	)
	group_ag, members, err := c.GetGroupAndMembersAdGraph(groupId)
	if err != nil {
		settings.ErrorLogger.Println(err)
		return (models.DynamicGroup{}), members, err
	}

	if len(group_ag.GroupType) > 0 && group_ag.GroupType[0] == enums.DynamicGroupType {
		dgroup.DisplayName = group_ag.DisplayName
		dgroup.MembershipRule = group_ag.MembershipRule
		dgroup.ObjectId = group_ag.ObjectId
		return dgroup, members, nil
	} else {
		return (models.DynamicGroup{}), members, nil
	}
}

func ParseConditionalAccessPolicyListADGraph(c *client.AzureClient) ([]models.ConditionalAccessPolicy, error) {
	var (
		response []models.ConditionalAccessPolicy
	)

	adGraphCaps, err := c.GetConditionalAccessPoliciesAdGraph()
	if err != nil {
		settings.ErrorLogger.Println(err)
		return response, nil
	}

	for _, capPolicy := range adGraphCaps {
		result, err := ParseConditionalAccessPolicyADGraph([]byte(capPolicy.PolicyDetail[0]), capPolicy.ObjectId, capPolicy.DisplayName, c)
		if err != nil {
			log.Fatal(err)
		}
		response = append(response, result)
	}

	return response, nil
}

func ParseConditionalAccessPolicyADGraph(policyDetail []byte, objectId string, displayName string, c *client.AzureClient) (models.ConditionalAccessPolicy, error) {
	var (
		polDetail adgraph.PolicyDetail
		cap       models.ConditionalAccessPolicy
	)

	settings.InfoLogger.Println("Parsing " + displayName)

	err := json.Unmarshal(policyDetail, &polDetail)
	if err != nil {
		settings.ErrorLogger.Println(err)
		return cap, err
	}
	cap.DisplayName = displayName
	cap.ObjectId = objectId
	cap.State = polDetail.State

	if cap.State != enums.StateEnabled {
		// Skip the rest of the parsing for now if state is not Enabled
		return cap, nil
	}

	//TODO: expand hardcoded applications such as Microsoft 365, or the windows services API into it's subparts
	for _, application := range polDetail.Conditions.ApplicationList.IncludedApplications {
		if len(application.Acrs) > 0 {
			if slices.Contains(application.Acrs, enums.AuthenticationContextPIM) {
				cap.ApplicationMode = enums.ApplicationModeAuthenticationContext
			} else {
				cap.ApplicationMode = enums.ApplicationModeUserActions
			}
			cap.IncludedActions = append(cap.IncludedActions, application.Acrs...)
		} else {
			cap.ApplicationMode = enums.ApplicationModeCloudApps
			cap.IncludedApplications = append(cap.IncludedApplications, application.Applications...)
		}
	}

	for _, application := range polDetail.Conditions.ApplicationList.ExcludedApplications {
		cap.ExcludedApplications = append(cap.ExcludedApplications, application.Applications...)
	}

	//include users, parse groups and roles
	for _, userGrouping := range polDetail.Conditions.UserList.IncludedUsers {
		if userGrouping.GuestsOrExternalUsers.GuestOrExternalUserTypes != "" {
			cap.IncludedUsers = append(cap.IncludedUsers, string(userGrouping.GuestsOrExternalUsers.GuestOrExternalUserTypes))
		}
		cap.IncludedUsers = append(cap.IncludedUsers, userGrouping.Users...)
		if len(cap.IncludedUsers) == 0 || (len(cap.IncludedUsers) > 0 && cap.IncludedUsers[0] != enums.AllUsers) {
			for _, groupId := range userGrouping.Groups {
				dgroup, members, err := ParseGroupADGraph(c, groupId)
				if err != nil {
					settings.ErrorLogger.Println(err)
					return cap, err
				}
				cap.IncludedUsers = append(cap.IncludedUsers, members...)
				if dgroup != (models.DynamicGroup{}) {
					cap.IncludedDynamicGroups = append(cap.IncludedDynamicGroups, dgroup)
				}
			}
			for _, roleId := range userGrouping.Roles {
				roleMembers, err := c.GetRoleMembersAdGraph(roleId)
				if err != nil {
					settings.ErrorLogger.Println(err)
					return cap, err
				}
				cap.IncludedUsers = append(cap.IncludedUsers, roleMembers...)
			}
		}
	}
	//exclude as needed
	for _, userGrouping := range polDetail.Conditions.UserList.ExcludedUsers {
		if userGrouping.GuestsOrExternalUsers.GuestOrExternalUserTypes != "" {
			cap.ExcludedUsers = append(cap.ExcludedUsers, string(userGrouping.GuestsOrExternalUsers.GuestOrExternalUserTypes))
		}
		cap.ExcludedUsers = append(cap.ExcludedUsers, userGrouping.Users...)

		for _, groupId := range userGrouping.Groups {
			dgroup, members, err := ParseGroupADGraph(c, groupId)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return cap, err
			}
			cap.ExcludedUsers = append(cap.ExcludedUsers, members...)
			if dgroup != (models.DynamicGroup{}) {
				cap.ExcludedDynamicGroups = append(cap.ExcludedDynamicGroups, dgroup)
			}
		}
		for _, roleId := range userGrouping.Roles {
			roleMembers, err := c.GetRoleMembersAdGraph(roleId)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return cap, err
			}
			cap.ExcludedUsers = append(cap.ExcludedUsers, roleMembers...)
		}

	}

	locations, err := ParseLocationsADGraph(c)
	if err != nil {
		settings.ErrorLogger.Println(err)
		return cap, err
	}

	// include locations
	for _, locationGrouping := range polDetail.Conditions.LocationList.IncludedLocations {
		for _, locationEntry := range locationGrouping.Locations {
			if locationEntry == string(enums.AllLocations) {
				for _, location := range locations {
					cap.IncludedLocations = append(cap.IncludedLocations, location.PolicyId)
				}
			} else if locationEntry == string(enums.AllTrustedLocations) {
				for _, location := range locations {
					if location.IsTrusted {
						cap.IncludedLocations = append(cap.IncludedLocations, location.PolicyId)
					}
				}
			} else {
				cap.IncludedLocations = append(cap.IncludedLocations, locationEntry)
			}
		}
	}

	for _, locationGrouping := range polDetail.Conditions.LocationList.ExcludedLocations {
		for _, locationEntry := range locationGrouping.Locations {
			if locationEntry == string(enums.AllTrustedLocations) {
				for _, location := range locations {
					if location.IsTrusted {
						cap.ExcludedLocations = append(cap.ExcludedLocations, location.PolicyId)
					}
				}
			} else {
				cap.ExcludedLocations = append(cap.ExcludedLocations, locationEntry)
			}
		}
	}

	//included client types
	for _, clientType := range polDetail.Conditions.ClientTypeList.IncludedClientTypes {
		cap.IncludedClientTypes = append(cap.IncludedClientTypes, clientType.ClientTypes...)
	}

	// Parsing included device platforms
	for _, devicePlatform := range polDetail.Conditions.DevicePlatformList.IncludedDevicePlatforms {
		cap.IncludedDevicePlatforms = make(map[enums.DevicePlatform]struct{})
		for _, devicePlatformString := range devicePlatform.DevicePlatforms {
			if devicePlatformString == enums.DevicePlatformAll {
				cap.IncludedDevicePlatforms[enums.DevicePlatformAndroid] = struct{}{}
				cap.IncludedDevicePlatforms[enums.DevicePlatformIOS] = struct{}{}
				cap.IncludedDevicePlatforms[enums.DevicePlatformLinux] = struct{}{}
				cap.IncludedDevicePlatforms[enums.DevicePlatformMacOS] = struct{}{}
				cap.IncludedDevicePlatforms[enums.DevicePlatformWindows] = struct{}{}
				cap.IncludedDevicePlatforms[enums.DevicePlatformWindowsPhone] = struct{}{}
			} else {
				cap.IncludedDevicePlatforms[devicePlatformString] = struct{}{}
			}
		}
	}

	//Excluding device platforms
	for _, devicePlatform := range polDetail.Conditions.DevicePlatformList.ExcludedDevicePlatforms {
		for _, devicePlatformString := range devicePlatform.DevicePlatforms {
			delete(cap.IncludedDevicePlatforms, devicePlatformString)
		}
	}

	// included, excluded devices
	if len(polDetail.Conditions.DeviceList.IncludedDevices) > 0 {
		cap.IncludedDevices = polDetail.Conditions.DeviceList.IncludedDevices[0].DeviceRule
	}
	if len(polDetail.Conditions.DeviceList.ExcludedDevices) > 0 {
		cap.ExcludedDevices = polDetail.Conditions.DeviceList.ExcludedDevices[0].DeviceRule
	}

	// Parse controls
	for _, control := range polDetail.Controls {
		var ca_control models.ConditionalAccessControl
		for _, ctrl := range control.Controls {
			ca_control.EnforcedControl = append(ca_control.EnforcedControl, string(ctrl))
		}
		ca_control.EnforcedControl = append(ca_control.EnforcedControl, control.AuthStrengthIds...)
		cap.Controls = append(cap.Controls, ca_control)
	}

	return cap, nil
}
