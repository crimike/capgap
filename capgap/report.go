package capgap

import (
	"capgap/enums"
	"capgap/models"
	"capgap/parsers"
	"capgap/settings"
	"slices"
)

func GroupBypassesByClientId(bypassList []models.Bypass) []models.Bypass {
	var response []models.Bypass

	for _, bypass := range bypassList {
		var aux models.Bypass
		aux.UserId = bypass.UserId
		aux.ApplicationId = bypass.ApplicationId
		aux.DevicePlatform = bypass.DevicePlatform
		aux.LocationId = bypass.LocationId
		allCombinations := true
		for _, clientType := range enums.GetAllClientTypes() {
			aux.ClientType = clientType
			if !slices.Contains(bypassList, aux) {
				allCombinations = false
				break
			}
		}
		if allCombinations {
			aux.ClientType = enums.AnyClientType
			if !slices.Contains(response, aux) {
				response = append(response, aux)
			}
		} else {
			response = append(response, bypass)
		}
	}

	return response
}

func GroupBypassesByDevicePlatform(bypassList []models.Bypass) []models.Bypass {
	var response []models.Bypass

	for _, bypass := range bypassList {
		var aux models.Bypass
		aux.UserId = bypass.UserId
		aux.ApplicationId = bypass.ApplicationId
		aux.ClientType = bypass.ClientType
		aux.LocationId = bypass.LocationId
		allCombinations := true
		for _, devicePlatform := range enums.GetAllDevicePlatforms() {
			aux.DevicePlatform = devicePlatform
			if !slices.Contains(bypassList, aux) {
				allCombinations = false
				break
			}
		}
		if allCombinations {
			aux.DevicePlatform = enums.AnyDevicePlatform
			if !slices.Contains(response, aux) {
				response = append(response, aux)
			}
		} else {
			response = append(response, bypass)
		}
	}

	return response
}

func GroupBypassesByLocation(bypassList []models.Bypass) []models.Bypass {
	var response []models.Bypass

	allLocations, err := parsers.ParseLocations()
	if err != nil {
		settings.ErrorLogger.Println(err)
		return response
	}

	for _, bypass := range bypassList {
		var aux models.Bypass
		aux.UserId = bypass.UserId
		aux.ApplicationId = bypass.ApplicationId
		aux.ClientType = bypass.ClientType
		aux.DevicePlatform = bypass.DevicePlatform
		allCombinations := true

		for _, location := range allLocations {
			aux.LocationId = location.PolicyId
			if !slices.Contains(bypassList, aux) {
				allCombinations = false
				break
			}
		}
		if allCombinations {
			aux.LocationId = enums.AnyLocation
			if !slices.Contains(response, aux) {
				response = append(response, aux)
			}
		} else {
			response = append(response, bypass)
		}
	}

	return response
}
