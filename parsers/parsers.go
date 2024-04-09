package parsers

import (
	"capgap/client"
	"capgap/models"
	"capgap/settings"
	"encoding/json"
	"os"
	"slices"
	"strconv"
)

type DataCache struct {
	Locations                 []models.Location
	Applications              []models.Application
	Users                     []models.User
	ConditionalAccessPolicies []models.ConditionalAccessPolicy
}

var Cache DataCache

func ParseAll() error {
	err := ParseLocations()
	if err != nil {
		return err
	}
	err = ParseApplications()
	if err != nil {
		return err
	}
	err = ParseUsers()
	if err != nil {
		return err
	}
	err = ParseConditionalAccessPolicyList()
	if err != nil {
		return err
	}
	return nil
}

func GetUserById(userId string) models.User {
	//if len of users == 0, these weren't parsed so we can try to query
	idx := slices.IndexFunc(Cache.Users, func(u models.User) bool { return u.ObjectId == userId })
	if idx == -1 {
		return models.User{}
	}
	return Cache.Users[idx]
}

func GetAppById(appId string) models.Application {
	// if len of applications == 0, these weren't parsed so we can try to query
	idx := slices.IndexFunc(Cache.Applications, func(a models.Application) bool { return a.ApplicationId == appId })
	if idx == -1 {
		return models.Application{}
	}
	return Cache.Applications[idx]
}

func GetLocationById(locationId string) models.Location {
	// if len of locations == 0, these weren't parsed so we can try to query - chances are that for locations we would need to parse all of them anyway
	idx := slices.IndexFunc(Cache.Locations, func(l models.Location) bool { return l.PolicyId == locationId })
	if idx == -1 {
		return models.Location{}
	}
	return Cache.Locations[idx]
}

func ParseLocations() error {
	var (
		c         client.AzureClient
		locations []models.Location
		err       error
	)

	if len(Cache.Locations) > 0 {
		return nil
	}

	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.DebugLogger.Println("Reading locations from " + settings.LocationsFile)
		data, err := os.ReadFile(settings.LocationsFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		err = json.Unmarshal(data, &locations)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		settings.DebugLogger.Println("A total of " + strconv.Itoa(len(locations)) + " locations were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			locations, err = parseLocationsADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return err
			}
		} // else msgraph
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.DebugLogger.Println("Saving locations to file: " + settings.LocationsFile)
			content, _ := json.MarshalIndent(locations, "", " ")
			err = os.WriteFile(settings.LocationsFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save locations to file: " + err.Error())
				return err
			}
		}
	}
	Cache.Locations = append(Cache.Locations, locations...)
	return nil
}

func ParseApplications() error {
	var (
		c            client.AzureClient
		applications []models.Application
		err          error
	)

	if len(Cache.Applications) > 0 {
		return nil
	}

	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.DebugLogger.Println("Reading applications from " + settings.AppsFile)
		data, err := os.ReadFile(settings.AppsFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		err = json.Unmarshal(data, &applications)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		settings.DebugLogger.Println("A total of " + strconv.Itoa(len(applications)) + " applications were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			applications, err = parseApplicationsADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return err
			}
		} // else msgraph
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.DebugLogger.Println("Saving applications to file: " + settings.AppsFile)
			content, _ := json.MarshalIndent(applications, "", " ")
			err = os.WriteFile(settings.AppsFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save applications to file: " + err.Error())
				return err
			}
		}
	}

	Cache.Applications = append(Cache.Applications, applications...)
	return nil
}

func ParseUsers() error {
	var (
		c     client.AzureClient
		users []models.User
		err   error
	)

	if len(Cache.Users) > 0 {
		return nil
	}

	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.DebugLogger.Println("Reading users from " + settings.UserFile)
		data, err := os.ReadFile(settings.UserFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		err = json.Unmarshal(data, &users)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		settings.DebugLogger.Println("A total of " + strconv.Itoa(len(users)) + " users were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			users, err = parseUsersADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return err
			}
		} // else msgraph
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.DebugLogger.Println("Saving users to file: " + settings.UserFile)
			content, _ := json.MarshalIndent(users, "", " ")
			err = os.WriteFile(settings.UserFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save users to file: " + err.Error())
				return err
			}
		}
	}

	Cache.Users = append(Cache.Users, users...)
	return nil
}

func ParseConditionalAccessPolicyList() error {

	var (
		c    client.AzureClient
		caps []models.ConditionalAccessPolicy
		err  error
	)

	if len(Cache.ConditionalAccessPolicies) > 0 {
		return nil
	}

	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.DebugLogger.Println("Reading CAPS from " + settings.CapsFile)
		data, err := os.ReadFile(settings.CapsFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		err = json.Unmarshal(data, &caps)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return err
		}
		settings.DebugLogger.Println("A total of " + strconv.Itoa(len(caps)) + " conditional access policies were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			caps, err = parseConditionalAccessPolicyListADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return err
			}
		} // else MSGRAPH
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.DebugLogger.Println("Saving conditional access policies to file: " + settings.CapsFile)
			content, _ := json.MarshalIndent(caps, "", " ")
			err := os.WriteFile(settings.CapsFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save conditional access policies to file: " + err.Error())
				return err
			}
		}
	}

	Cache.ConditionalAccessPolicies = append(Cache.ConditionalAccessPolicies, caps...)

	return nil
}
