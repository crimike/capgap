package capgap

import (
	"capgap/enums"
	"capgap/models"
	"capgap/parsers"
	"capgap/settings"
	"fmt"
	"slices"
	"strings"
)

// TODO stringbuilder
func PrettyPrintBypass(b models.Bypass) string {

	var (
		result string
		loc    string
	)
	if b.LocationId != enums.AnyLocation {
		location := parsers.GetLocationById(b.LocationId)
		if location.PolicyId != "" {
			loc = location.DisplayName + " ("
			for _, ip := range location.IpRange {
				loc += ip + ", "
			}
			loc += ")"
		} else {
			loc = b.LocationId
		}
	} else {
		loc = enums.AnyLocation
	}
	result = "ClientType: " + string(b.ClientType) + "; DevicePlatform: " + string(b.DevicePlatform) + "; Location: " + loc + "\n"
	return result
}

// TODO stringbuilder
func PrettyPrintAppList(appList []string) string {
	var (
		result string
	)
	result = "Applications: [ "
	for _, appId := range appList {
		app := parsers.GetAppById(appId)
		if app == (models.Application{}) {
			result += appId + ", "
		} else {
			result += app.DisplayName + "(" + app.ApplicationId + "), "
		}
	}
	result += "]\n"
	return result
}

func PrettyPrintUserList(userLIst []string) string {
	var (
		builder strings.Builder
	)
	builder.WriteString("Users: [ ")
	for _, userId := range userLIst {
		// user := parsers.GetUserById(userId)
		// if user == (models.User{}) {
		builder.WriteString(userId + ", ")
		// } else {
		// 	result += user.DisplayName + "(" + user.ObjectId + "), "
		// }
	}
	builder.WriteString("]\n")
	return builder.String()
}

func GroupBypassesByClientId(bypassList []models.Bypass) []models.Bypass {
	var response []models.Bypass

	for _, bypass := range bypassList {
		var aux models.Bypass
		aux.UserId = bypass.UserId
		aux.ApplicationId = bypass.ApplicationId
		aux.DevicePlatform = bypass.DevicePlatform
		aux.LocationId = bypass.LocationId
		allCombinations := true

		aux.ClientType = enums.AnyClientType
		if slices.Contains(response, aux) {
			continue
		}

		for _, clientType := range enums.GetAllClientTypes() {
			aux.ClientType = clientType
			if !slices.Contains(bypassList, aux) {
				allCombinations = false
				break
			}
		}
		if allCombinations {
			aux.ClientType = enums.AnyClientType
			response = append(response, aux)
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

		aux.DevicePlatform = enums.AnyDevicePlatform
		if slices.Contains(response, aux) {
			continue
		}

		for _, devicePlatform := range enums.GetAllDevicePlatforms() {
			aux.DevicePlatform = devicePlatform
			if !slices.Contains(bypassList, aux) {
				allCombinations = false
				break
			}
		}
		if allCombinations {
			aux.DevicePlatform = enums.AnyDevicePlatform
			response = append(response, aux)
		} else {
			response = append(response, bypass)
		}
	}

	return response
}

func GroupBypassesByLocation(bypassList []models.Bypass) []models.Bypass {
	var response []models.Bypass

	for _, bypass := range bypassList {
		var aux models.Bypass
		aux.UserId = bypass.UserId
		aux.ApplicationId = bypass.ApplicationId
		aux.ClientType = bypass.ClientType
		aux.DevicePlatform = bypass.DevicePlatform
		allCombinations := true

		aux.LocationId = enums.AnyLocation
		if slices.Contains(response, aux) {
			continue
		}

		for _, location := range parsers.Cache.Locations {
			aux.LocationId = location.PolicyId
			if !slices.Contains(bypassList, aux) {
				allCombinations = false
				break
			}
		}
		if allCombinations {
			aux.LocationId = enums.AnyLocation
			response = append(response, aux)
		} else if settings.Config[settings.ALL_LOCATIONS] == settings.TRUE { //only append bypasses with other locations is -allLocations is specified
			response = append(response, bypass)
		}
	}

	return response
}

// Assumption is that all bypasses have same user/application
func ReportBypassesUserApp(bypasses *[]models.Bypass) {
	if len(*bypasses) == 0 {
		return
	}
	user := parsers.GetUserById((*bypasses)[0].UserId)
	if user == (models.User{}) {
		user = models.User{
			DisplayName: "UNKNOWN USER",
			ObjectId:    (*bypasses)[0].UserId,
		}
	}

	application := parsers.GetAppById((*bypasses)[0].ApplicationId)
	if application == (models.Application{}) {
		application = models.Application{
			DisplayName:   "UNKNOWN APPLICATION",
			ApplicationId: (*bypasses)[0].ApplicationId,
		}
	}

	settings.Reporter.WriteString("Bypasses for user " + user.DisplayName + " going to " + application.DisplayName + "\n")
	for i := range *bypasses {
		settings.Reporter.WriteString(PrettyPrintBypass((*bypasses)[i]))
	}
}

// Assumption is that all bypasses have the same user
func ReportForUser(bypasses *[]models.Bypass) {

	groupedBypasses := make(map[models.Bypass][]string)

	if len(*bypasses) == 0 {
		return
	}
	user := parsers.GetUserById((*bypasses)[0].UserId)
	if user == (models.User{}) {
		user = models.User{
			DisplayName: "UNKNOWN USER",
			ObjectId:    (*bypasses)[0].UserId,
		}
	}

	for i := range *bypasses {
		b := models.Bypass{
			ClientType:     (*bypasses)[i].ClientType,
			DevicePlatform: (*bypasses)[i].DevicePlatform,
			LocationId:     (*bypasses)[i].LocationId,
			UserId:         (*bypasses)[i].UserId,
		}
		groupedBypasses[b] = append(groupedBypasses[b], (*bypasses)[i].ApplicationId)
	}
	settings.Reporter.WriteString("Reporting for " + user.DisplayName + "(" + user.ObjectId + ")\n")
	noControl := models.Bypass{
		ClientType:     enums.AnyClientType,
		DevicePlatform: enums.AnyDevicePlatform,
		LocationId:     enums.AnyLocation,
		UserId:         user.ObjectId,
	}

	allApplicationCount := len(parsers.Cache.Applications)

	// no control for all the applications
	if len(groupedBypasses[noControl]) == allApplicationCount {
		settings.Reporter.WriteString("No control for this user against any app\n")
	} else if len(groupedBypasses[noControl]) != 0 {
		settings.Reporter.WriteString("No control for the following " + fmt.Sprint(len(groupedBypasses[noControl])) + " apps: \n")
		settings.Reporter.WriteString(PrettyPrintAppList(groupedBypasses[noControl]))
	}

	for bypass, appList := range groupedBypasses {
		if bypass == noControl {
			continue
		}
		settings.Reporter.WriteString(PrettyPrintBypass(bypass))
		settings.Reporter.WriteString(PrettyPrintAppList(appList))
	}
}

// Assumption is that all bypasses have the same app
func ReportForApp(bypasses *[]models.Bypass) {
	groupedBypasses := make(map[models.Bypass][]string)
	if len(*bypasses) == 0 {
		return
	}
	//TODO: computationally expensive, will add flag to enable resolving
	app := parsers.GetAppById((*bypasses)[0].ApplicationId)
	if app == (models.Application{}) {
		app = models.Application{
			DisplayName: "UNKNOWN APPLICATION",
			// DisplayName:   "",
			ApplicationId: (*bypasses)[0].ApplicationId,
		}
	}

	for i := range *bypasses {
		b := models.Bypass{
			ClientType:     (*bypasses)[i].ClientType,
			DevicePlatform: (*bypasses)[i].DevicePlatform,
			LocationId:     (*bypasses)[i].LocationId,
			ApplicationId:  (*bypasses)[i].ApplicationId,
		}
		groupedBypasses[b] = append(groupedBypasses[b], (*bypasses)[i].UserId)
	}
	settings.Reporter.WriteString("Reporting for " + app.DisplayName + "(" + app.ApplicationId + ")\n")
	noControl := models.Bypass{
		ClientType:     enums.AnyClientType,
		DevicePlatform: enums.AnyDevicePlatform,
		LocationId:     enums.AnyLocation,
		ApplicationId:  app.ApplicationId,
	}

	allUserCount := len(parsers.Cache.Users)

	// no control for all the applications
	if len(groupedBypasses[noControl]) == allUserCount {
		settings.Reporter.WriteString("No control for this app from any user\n")
	} else {
		userCount := len(groupedBypasses[noControl])
		settings.Reporter.WriteString("No control for the following " + fmt.Sprint(userCount) + " users: \n")
		if userCount > 10 {
			ul := append(groupedBypasses[noControl][:10], "...<rest ommited for brevity")
			settings.Reporter.WriteString(PrettyPrintUserList(ul))
		} else {
			settings.Reporter.WriteString(PrettyPrintUserList(groupedBypasses[noControl]))
		}
	}

	for bypass, userList := range groupedBypasses {
		if bypass == noControl {
			continue
		}
		settings.Reporter.WriteString(PrettyPrintBypass(bypass))
		if len(userList) > 10 {
			ul := append(userList[:10], "...<rest ommited for brevity")
			settings.Reporter.WriteString(PrettyPrintUserList(ul))
		} else {
			settings.Reporter.WriteString(PrettyPrintUserList(userList))
		}
	}
}

func ReportAll(bypasses *[]models.Bypass) {

	perApp := make(map[string][]models.Bypass)

	for i := range *bypasses {
		perApp[(*bypasses)[i].ApplicationId] = append(perApp[(*bypasses)[i].ApplicationId], (*bypasses)[i])
	}

	for _, bypassList := range perApp {
		ReportForApp(&bypassList)
		settings.Reporter.WriteString("======================================\n")
	}
}
