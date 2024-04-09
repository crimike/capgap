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

func GetBypassesUserToApp(cap models.ConditionalAccessPolicy, userId string, appId string) []models.Bypass {
	var bypasses []models.Bypass
	for _, location := range parsers.Cache.Locations {
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

func FindGapsPerUserAndApp(userId string, appId string) []models.Bypass {
	var (
		appliedCaps    []models.ConditionalAccessPolicy
		bypasses       [][]models.Bypass
		commonBypasses []models.Bypass
	)
	for _, cap := range parsers.Cache.ConditionalAccessPolicies {
		if CapAppliesToApplication(cap, appId) && CapAppliesToUser(cap, userId) && len(cap.Controls) > 0 {
			appliedCaps = append(appliedCaps, cap)
		}
	}

	// No conditional access policies apply, thus everything is a bypass
	if len(appliedCaps) == 0 {
		return []models.Bypass{
			{
				ApplicationId:  appId,
				ClientType:     enums.AnyClientType,
				DevicePlatform: enums.AnyDevicePlatform,
				LocationId:     enums.AnyLocation,
				UserId:         userId,
			},
		}
	}

	for _, cap := range appliedCaps {
		b := GetBypassesUserToApp(cap, userId, appId)
		bypasses = append(bypasses, b)
	}

	commonBypasses = GetCommonBypasses(bypasses)
	bypasses = nil
	commonBypasses = GroupBypassesByClientId(commonBypasses)
	commonBypasses = GroupBypassesByDevicePlatform(commonBypasses)
	commonBypasses = GroupBypassesByLocation(commonBypasses)

	return commonBypasses
}

func FindGapsForUser(userId string) []models.Bypass {
	var (
		appBypases []models.Bypass
	)

	settings.DebugLogger.Println("Parsing gaps for user " + userId + " against all apps")

	for i := range parsers.Cache.Applications {
		x := FindGapsPerUserAndApp(userId, parsers.Cache.Applications[i].ApplicationId)
		appBypases = append(appBypases, x...)
	}

	return appBypases
}

func FindGapsForApp(appId string) []models.Bypass {
	var (
		userBypasses []models.Bypass
		useAllUsers  bool
	)

	settings.DebugLogger.Println("Parsing gaps for application " + appId + " against all users")

	// Optimization: attempt to parse applied caps on the application, and if not all users are included, we could limit the list to the included users
	useAllUsers = false
	for _, cap := range parsers.Cache.ConditionalAccessPolicies {
		if CapAppliesToApplication(cap, appId) && len(cap.Controls) > 0 {
			if cap.IncludedUsers[0] == enums.AllUsers {
				useAllUsers = true
			}
		}
	}

	if useAllUsers {
		for i := range parsers.Cache.Users {
			ub := FindGapsPerUserAndApp(parsers.Cache.Users[i].ObjectId, appId)
			userBypasses = append(userBypasses, ub...)
		}
	} else {
		var userIds []string
		for _, cap := range parsers.Cache.ConditionalAccessPolicies {
			if CapAppliesToApplication(cap, appId) && len(cap.Controls) > 0 {
				userIds = append(userIds, cap.IncludedUsers...)
			}
		}
		for i := range userIds {
			ub := FindGapsPerUserAndApp(userIds[i], appId)
			userBypasses = append(userBypasses, ub...)
		}
	}

	return userBypasses
}

// due to optimization, reporting happens inline
func FindAllGaps() {
	appCount := len(parsers.Cache.Applications)
	workerCount := 50

	results := make(chan []models.Bypass, workerCount)
	jobs := make(chan string, appCount)

	for w := 0; w < workerCount; w++ {
		go func(jobs <-chan string, results chan<- []models.Bypass) {
			for appId := range jobs {
				b := FindGapsForApp(appId)
				results <- b
			}

		}(jobs, results)
	}

	for i := range parsers.Cache.Applications {
		jobs <- parsers.Cache.Applications[i].ApplicationId
	}

	close(jobs)
	settings.InfoLogger.Println("All jobs started")

	for i := 0; i < appCount; i++ {
		bps := <-results
		if i%10 == 0 {
			settings.InfoLogger.Println("Parsed " + fmt.Sprint(i) + "/" + fmt.Sprint(appCount) + " applications. Current ammount of bypasses: " + fmt.Sprint(len(bps)))
		}
		ReportAll(&bps)
	}
}

// func FindGapsPerCap
