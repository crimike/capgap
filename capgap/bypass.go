package capgap

import (
	"capgap/enums"
	"capgap/models"
	"capgap/parsers"
	"fmt"
	"log"
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

func FindGapsPerUserAndApp(caps []models.ConditionalAccessPolicy, userId string, appId string) {
	var (
		appliedCaps []models.ConditionalAccessPolicy
		bypasses    [][]models.Bypass
	)
	for _, cap := range caps {
		if CapAppliesToApplication(cap, appId) && CapAppliesToUser(cap, userId) && len(cap.Controls) > 0 {
			appliedCaps = append(appliedCaps, cap)
		}
	}

	locations, err := parsers.ParseLocations()
	if err != nil {
		log.Println(err)
		return
	}

	for _, cap := range appliedCaps {
		b := GetBypassesUserToApp(cap, userId, appId, locations)
		bypasses = append(bypasses, b)
	}

	commonBypasses := GetCommonBypasses(bypasses)

	//TODO: generate report or pretty print bypasses
	commonBypasses = GroupBypassesByClientId(commonBypasses)
	commonBypasses = GroupBypassesByDevicePlatform(commonBypasses)
	commonBypasses = GroupBypassesByLocation(commonBypasses)
	fmt.Println(commonBypasses)
}

// func FindGapsPerUser(caps []models.ConditionalAccessPolicy, userId){
// 	// for application in applications:
// 	// 	FindGapsPerUserAndApp(caps, userId, appId)
// }

// func FindGapsPerApp

// func FindGapsPerCap
