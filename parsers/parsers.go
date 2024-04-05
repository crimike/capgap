package parsers

import (
	"capgap/client"
	"capgap/models"
	"capgap/settings"
	"encoding/json"
	"os"
	"strconv"
)

//TODO move cache of objects here instead of the client(it will also apply for MSGraph)

func ParseLocations() ([]models.Location, error) {
	var (
		c         client.AzureClient
		locations []models.Location
		err       error
	)
	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading locations from " + settings.LocationsFile)
		data, err := os.ReadFile(settings.LocationsFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return locations, err
		}
		err = json.Unmarshal(data, &locations)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return locations, err
		}
		settings.InfoLogger.Println("A total of " + strconv.Itoa(len(locations)) + " locations were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			locations, err = ParseLocationsADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return locations, err
			}
		} // else msgraph
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving locations to file: " + settings.LocationsFile)
			content, _ := json.MarshalIndent(locations, "", " ")
			err = os.WriteFile(settings.LocationsFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save locations to file: " + err.Error())
			}
		}
	}
	return locations, nil
}

func ParseApplications() ([]models.Application, error) {
	var (
		c            client.AzureClient
		applications []models.Application
		err          error
	)
	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading applications from " + settings.AppsFile)
		data, err := os.ReadFile(settings.AppsFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return applications, err
		}
		err = json.Unmarshal(data, &applications)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return applications, err
		}
		settings.InfoLogger.Println("A total of " + strconv.Itoa(len(applications)) + " applications were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			applications, err = ParseApplicationsADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return applications, err
			}
		} // else msgraph
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving applications to file: " + settings.AppsFile)
			content, _ := json.MarshalIndent(applications, "", " ")
			err = os.WriteFile(settings.AppsFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save applications to file: " + err.Error())
			}
		}
	}

	return applications, nil
}

func ParseUsers() ([]models.User, error) {
	var (
		c     client.AzureClient
		users []models.User
		err   error
	)
	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading users from " + settings.UserFile)
		data, err := os.ReadFile(settings.UserFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return users, err
		}
		err = json.Unmarshal(data, &users)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return users, err
		}
		settings.InfoLogger.Println("A total of " + strconv.Itoa(len(users)) + " users were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			users, err = ParseUsersADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return users, err
			}
		} // else msgraph
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving users to file: " + settings.UserFile)
			content, _ := json.MarshalIndent(users, "", " ")
			err = os.WriteFile(settings.UserFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save users to file: " + err.Error())
			}
		}
	}

	return users, nil
}

func ParseConditionalAccessPolicyList() ([]models.ConditionalAccessPolicy, error) {

	var (
		c    client.AzureClient
		caps []models.ConditionalAccessPolicy
		err  error
	)

	c.InitializeClient()

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading CAPS from " + settings.CapsFile)
		data, err := os.ReadFile(settings.CapsFile)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return caps, err
		}
		err = json.Unmarshal(data, &caps)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return caps, err
		}
		settings.InfoLogger.Println("A total of " + strconv.Itoa(len(caps)) + " conditional access policies were retrieved from file")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			caps, err = ParseConditionalAccessPolicyListADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return caps, err
			}
		} // else MSGRAPH
		if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving conditional access policies to file: " + settings.CapsFile)
			content, _ := json.MarshalIndent(caps, "", " ")
			err := os.WriteFile(settings.CapsFile, content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save conditional access policies to file: " + err.Error())
			}
		}
	}

	return caps, nil
}
