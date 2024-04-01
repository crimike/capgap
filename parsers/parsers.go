package parsers

import (
	"capgap/client"
	"capgap/models"
	"capgap/settings"
	"encoding/json"
	"log"
	"os"
)

func ParseLocations() ([]models.Location, error) {
	var (
		c         client.AzureClient
		locations []models.Location
	)
	c.InitializeClient()

	if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
		locations, err := ParseLocationsADGraph(&c)
		if err != nil {
			log.Println(err)
			return locations, err
		}
	} // else msgraph
	return locations, nil
}

func ParseConditionalAccessPolicyList() ([]models.ConditionalAccessPolicy, error) {

	var (
		c    client.AzureClient
		caps []models.ConditionalAccessPolicy
	)

	c.InitializeClient()

	if settings.Config[settings.CAPFILE_DIRECTION] == settings.LOADCAP {
		data, err := os.ReadFile(settings.CAPFILE)
		err = json.Unmarshal(data, &caps)
		if err != nil {
			log.Println(err)
			return caps, err
		}
	} else {
		if settings.Config[settings.CLIENTENDPOINT] == settings.AADGRAPH {
			caps, err := ParseConditionalAccessPolicyListADGraph(&c)
			if err != nil {
				log.Println(err)
				return caps, err
			}
		} // else MSGRAPH
		if settings.Config[settings.CAPFILE_DIRECTION] == settings.SAVECAP {
			content, _ := json.MarshalIndent(caps, "", " ")
			_ = os.WriteFile(settings.CAPFILE, content, 0644)
		}
	}

	return caps, nil
}
