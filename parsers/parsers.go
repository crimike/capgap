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

	if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
		applications, err = ParseApplicationsADGraph(&c)
		if err != nil {
			settings.ErrorLogger.Println(err)
			return applications, err
		}
	} // else msgraph
	return applications, nil
}

func ParseConditionalAccessPolicyList() ([]models.ConditionalAccessPolicy, error) {

	var (
		c    client.AzureClient
		caps []models.ConditionalAccessPolicy
	)

	c.InitializeClient()

	if settings.Config[settings.CAPFILE_DIRECTION] == settings.LOADCAP {
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
		settings.InfoLogger.Println("A total of " + strconv.Itoa(len(caps)) + " conditional access policies were retrieved")
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			caps, err := ParseConditionalAccessPolicyListADGraph(&c)
			if err != nil {
				settings.ErrorLogger.Println(err)
				return caps, err
			}
		} // else MSGRAPH
		if settings.Config[settings.CAPFILE_DIRECTION] == settings.SAVECAP {
			settings.InfoLogger.Println("Saving conditional access policies to file: " + settings.Config[settings.CAPFILE])
			content, _ := json.MarshalIndent(caps, "", " ")
			_ = os.WriteFile(settings.Config[settings.CAPFILE], content, 0644)
		}
	}

	return caps, nil
}
