package parsers

import (
	"capgap/client"
	"capgap/models"
	"capgap/settings"
	"encoding/json"
	"os"
	"strconv"
)

func ParseLocations() ([]models.Location, error) {
	var (
		c         client.AzureClient
		locations []models.Location
		err       error
	)
	c.InitializeClient()

	if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
		locations, err = ParseLocationsADGraph(&c)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return locations, err
		}
	} // else msgraph
	return locations, nil
}

func ParseApplications() ([]models.Application, error) {
	var (
		c            client.AzureClient
		applications []models.Application
		err          error
	)
	c.InitializeClient()

	if settings.Config[settings.APPFILE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading applications from " + settings.Config[settings.APPFILE])
		data, err := os.ReadFile(settings.Config[settings.APPFILE])
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
		if settings.Config[settings.APPFILE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving applications to file: " + settings.Config[settings.APPFILE])
			content, _ := json.MarshalIndent(applications, "", " ")
			err = os.WriteFile(settings.Config[settings.APPFILE], content, 0644)
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

	if settings.Config[settings.USERFILE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading users from " + settings.Config[settings.USERFILE])
		data, err := os.ReadFile(settings.Config[settings.USERFILE])
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
		if settings.Config[settings.USERFILE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving users to file: " + settings.Config[settings.USERFILE])
			content, _ := json.MarshalIndent(users, "", " ")
			err = os.WriteFile(settings.Config[settings.USERFILE], content, 0644)
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
	)

	c.InitializeClient()

	if settings.Config[settings.CAPFILE_DIRECTION] == settings.LOAD {
		settings.InfoLogger.Println("Reading CAPS from " + settings.Config[settings.CAPFILE])
		data, err := os.ReadFile(settings.Config[settings.CAPFILE])
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
			caps, err := ParseConditionalAccessPolicyListADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return caps, err
			}
		} // else MSGRAPH
		if settings.Config[settings.CAPFILE_DIRECTION] == settings.SAVE {
			settings.InfoLogger.Println("Saving conditional access policies to file: " + settings.Config[settings.CAPFILE])
			content, _ := json.MarshalIndent(caps, "", " ")
			err := os.WriteFile(settings.Config[settings.CAPFILE], content, 0644)
			if err != nil {
				settings.ErrorLogger.Println("Could not save conditional access policies to file: " + err.Error())
			}
		}
	}

	return caps, nil
}
