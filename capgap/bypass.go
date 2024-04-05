package capgap

import (
	"capgap/enums"
	"capgap/models"
	"capgap/parsers"
	"capgap/settings"
	"fmt"
	"slices"
)

func GetCapBypass(cap models.ConditionalAccessPolicy, userId string, appId string, location models.Location, devicePlatform enums.DevicePlatform, clientType enums.ClientType) models.Bypass {
	var bypass models.Bypass

	if cap.State == enums.StateEnabled {
		if len(cap.Controls) == 0 || !CapAppliesToApplication(cap, appId) || !CapAppliesToClientType(cap, clientType) || !CapAppliesToDevicePlatform(cap, devicePlatform) || !CapAppliesToLocation(cap, location) || !CapAppliesToUser(cap, userId) {
			bypass.ApplicationId = appId
			bypass.ClientType = clientType
			bypass.DevicePlatform = devicePlatform
			bypass.LocationId = location.PolicyId
			bypass.UserId = userId
		}
	}
	return bypass

}

func GetBypassesForClientType(cap models.ConditionalAccessPolicy, userId string, appId string, location models.Location, devicePlatform enums.DevicePlatform) []models.Bypass {
	var bypasses []models.Bypass
	for _, clientType := range enums.GetAllClientTypes() {
		b := GetCapBypass(cap, userId, appId, location, devicePlatform, clientType)
		if b != (models.Bypass{}) {
			bypasses = append(bypasses, b)
		}
	}
	return bypasses
}

func GetBypassesForDevicePlatform(cap models.ConditionalAccessPolicy, userId string, appId string, location models.Location, clientType enums.ClientType) []models.Bypass {
	var bypasses []models.Bypass
	for _, devicePlatform := range enums.GetAllDevicePlatforms() {
		b := GetCapBypass(cap, userId, appId, location, devicePlatform, clientType)
		if b != (models.Bypass{}) {
			bypasses = append(bypasses, b)
		}
	}
	return bypasses
}

func GetBypassesForLocation(cap models.ConditionalAccessPolicy, userId string, appId string, locations []models.Location, devicePlatform enums.DevicePlatform, clientType enums.ClientType) []models.Bypass {
	var bypasses []models.Bypass
	for _, location := range locations {
		b := GetCapBypass(cap, userId, appId, location, devicePlatform, clientType)
		if b != (models.Bypass{}) {
			bypasses = append(bypasses, b)
		}
	}
	return bypasses
}

func GetBypassesForApplication(cap models.ConditionalAccessPolicy, userId string, applications []models.Application, location models.Location, devicePlatform enums.DevicePlatform, clientType enums.ClientType) []models.Bypass {
	var bypasses []models.Bypass
	for _, app := range applications {
		b := GetCapBypass(cap, userId, app.ApplicationId, location, devicePlatform, clientType)
		if b != (models.Bypass{}) {
			bypasses = append(bypasses, b)
		}
	}
	return bypasses
}

func GetBypassesForUser(cap models.ConditionalAccessPolicy, users []models.User, appId string, location models.Location, devicePlatform enums.DevicePlatform, clientType enums.ClientType) []models.Bypass {
	var bypasses []models.Bypass
	//TOD: can a cap not contain included users
	for _, user := range users {
		b := GetCapBypass(cap, user.ObjectId, appId, location, devicePlatform, clientType)
		if b != (models.Bypass{}) {
			bypasses = append(bypasses, b)
		}
	}
	return bypasses
}

func GetBypassesUserToApp(cap models.ConditionalAccessPolicy, userId string, appId string, locations []models.Location) []models.Bypass {
	var bypasses []models.Bypass
	for _, location := range locations {
		for _, devicePlatform := range enums.GetAllDevicePlatforms() {
			bypasses = append(bypasses, GetBypassesForClientType(cap, userId, appId, location, devicePlatform)...)
		}
	}
	return bypasses
}

func GetCommonBypasses(allBypasses [][]models.Bypass) []models.Bypass {

	var response []models.Bypass
	if len(allBypasses) == 0 {
		return response
	}
	if len(allBypasses) == 1 {
		return allBypasses[0]
	}
	for _, bypass := range allBypasses[0] {
		isCommon := true
		for _, bypassArray := range allBypasses[1:] {
			if !slices.Contains(bypassArray, bypass) {
				isCommon = false
				break
			}
		}
		if isCommon {
			response = append(response, bypass)
		}
	}
	return response
}

func FindGapsPerUserAndApp(caps []models.ConditionalAccessPolicy, userId string, appId string) ([]models.Bypass, error) {
	var (
		appliedCaps    []models.ConditionalAccessPolicy
		bypasses       [][]models.Bypass
		commonBypasses []models.Bypass
	)
	for _, cap := range caps {
		if CapAppliesToApplication(cap, appId) && CapAppliesToUser(cap, userId) && len(cap.Controls) > 0 {
			appliedCaps = append(appliedCaps, cap)
		}
	}

	// No conditional access policies apply, thus everything is a bypass
	if len(appliedCaps) == 0 {
		var b models.Bypass
		b.ApplicationId = appId
		b.ClientType = enums.AnyClientType
		b.DevicePlatform = enums.AnyDevicePlatform
		b.LocationId = enums.AnyLocation
		b.UserId = userId
		commonBypasses = append(commonBypasses, b)
		return commonBypasses, nil
	}

	locations, err := parsers.ParseLocations()
	if err != nil {
		settings.ErrorLogger.Println(err.Error())
		return commonBypasses, err
	}

	for _, cap := range appliedCaps {
		b := GetBypassesUserToApp(cap, userId, appId, locations)
		bypasses = append(bypasses, b)
	}

	count := 0
	for _, bl := range bypasses {
		count += len(bl)
	}
	settings.InfoLogger.Println("Found a total of " + fmt.Sprint(count) + " bypasses for user " + userId + " to app " + appId)

	commonBypasses = GetCommonBypasses(bypasses)
	commonBypasses = GroupBypassesByClientId(commonBypasses)
	commonBypasses = GroupBypassesByDevicePlatform(commonBypasses)
	commonBypasses = GroupBypassesByLocation(commonBypasses)
	settings.InfoLogger.Println("Out of those, there are " + fmt.Sprint(len(commonBypasses)) + " common bypasses")

	return commonBypasses, nil
}

func FindGapsForUser(caps []models.ConditionalAccessPolicy, userId string) (map[string][]models.Bypass, error) {
	var (
		appliedCaps []models.ConditionalAccessPolicy
		appBypases  map[string][]models.Bypass
	)
	appBypases = make(map[string][]models.Bypass)

	for _, cap := range caps {
		if CapAppliesToUser(cap, userId) && len(cap.Controls) > 0 {
			appliedCaps = append(appliedCaps, cap)
		}
	}

	applications, err := parsers.ParseApplications()
	if err != nil {
		settings.ErrorLogger.Println(err.Error())
		return appBypases, err
	}

	for _, app := range applications {
		bypasses, err := FindGapsPerUserAndApp(caps, userId, app.ApplicationId)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return appBypases, err
		}
		appBypases[app.ApplicationId] = append(appBypases[app.ApplicationId], bypasses...)
	}

	return appBypases, nil
}

func FindGapsForApp(caps []models.ConditionalAccessPolicy, appId string) (map[string][]models.Bypass, error) {
	var (
		appliedCaps  []models.ConditionalAccessPolicy
		userBypasses map[string][]models.Bypass
	)
	userBypasses = make(map[string][]models.Bypass)

	for _, cap := range caps {
		if CapAppliesToApplication(cap, appId) && len(cap.Controls) > 0 {
			appliedCaps = append(appliedCaps, cap)
		}
	}

	users, err := parsers.ParseUsers()
	if err != nil {
		settings.ErrorLogger.Println(err.Error())
		return userBypasses, err
	}

	for _, user := range users {
		bypasses, err := FindGapsPerUserAndApp(caps, user.ObjectId, appId)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return userBypasses, err
		}
		userBypasses[user.ObjectId] = append(userBypasses[user.ObjectId], bypasses...)
	}

	return userBypasses, nil
}

// func FindGapsPerCap
