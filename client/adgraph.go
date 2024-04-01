package client

import (
	"capgap/models/adgraph"
	"capgap/settings"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type AdGraphCache struct {
	GroupMembers map[string][]string
	RoleMembers  map[string][]string
	Locations    []adgraph.LocationPolicy
}

var cache AdGraphCache

func (c *AzureClient) InitializeAzureADGraphClient() {
	c.AccessToken = settings.Config[settings.ACCESSTOKEN]
	c.ApiVersion = "1.61-internal"
	c.MainUrl = "https://graph.windows.net/"
	c.Tenant = settings.Config[settings.TENANT]
	c.HttpClient = &http.Client{}

	cache.GroupMembers = make(map[string][]string)
	cache.RoleMembers = make(map[string][]string)
}

func (c *AzureClient) GetNamedLocationsAdGraph() ([]adgraph.LocationPolicy, error) {
	var response []adgraph.LocationPolicy

	if len(cache.Locations) > 0 {
		return cache.Locations, nil
	}

	apiUrl := c.MainUrl + c.Tenant + "/policies?$top=999&$filter=policyType%20eq%206&api-version=" + c.ApiVersion

	for apiUrl != "" {
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Println(err)
			return response, err
		}

		// Add the Bearer token to the request header
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			log.Println(err)
			return response, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return response, err
		}

		if resp.StatusCode != 200 {
			return response, errors.New("Return code is :" + fmt.Sprint(resp.StatusCode))
		}

		var lpr adgraph.LocationPolicyResponse

		err = json.Unmarshal(body, &lpr)
		if err != nil {
			log.Println(err)
			return response, err
		}

		response = append(response, lpr.LocationPolicies...)

		//rebuilding apiUrl
		if strings.Contains(lpr.NextLink, "directoryObjects") {
			apiUrl = c.MainUrl + c.Tenant + "/" + lpr.NextLink + "&$top=999&api-version=" + c.ApiVersion
		} else {
			apiUrl = lpr.NextLink
		}
	}

	cache.Locations = append(cache.Locations, response...)

	return response, nil

}

func (c *AzureClient) GetConditionalAccessPoliciesAdGraph() ([]adgraph.ConditionalAccessPolicy, error) {

	var response []adgraph.ConditionalAccessPolicy

	apiUrl := c.MainUrl + c.Tenant + "/policies?$top=999&$filter=policyType%20eq%2018&api-version=" + c.ApiVersion

	for apiUrl != "" {
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Println(err)
			return response, err
		}

		// Add the Bearer token to the request header
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			log.Println(err)
			return response, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return response, err
		}

		if resp.StatusCode != 200 {
			return response, errors.New("Return code is :" + fmt.Sprint(resp.StatusCode))
		}

		var caps adgraph.ConditionalAccessPolicyResponse

		err = json.Unmarshal(body, &caps)
		if err != nil {
			log.Println(err)
			return response, err
		}

		response = append(response, caps.ConditionalAccessPolicies...)

		//rebuilding apiUrl
		if strings.Contains(caps.NextLink, "directoryObjects") {
			apiUrl = c.MainUrl + c.Tenant + "/" + caps.NextLink + "&$top=999&api-version=" + c.ApiVersion
		} else {
			apiUrl = caps.NextLink
		}
	}

	return response, nil
}

func (c *AzureClient) GetApplicationsAdGraph() ([]adgraph.Application, error) {
	var (
		response []adgraph.Application
	)

	apiUrl := c.MainUrl + c.Tenant + "/applications?$top=999&api-version=" + c.ApiVersion

	for apiUrl != "" {
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Println(err)
			return response, err
		}

		// Add the Bearer token to the request header
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			log.Println(err)
			return response, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return response, err
		}

		if resp.StatusCode != 200 {
			return response, errors.New("Return code is :" + fmt.Sprint(resp.StatusCode))
		}

		var apps adgraph.ApplicationResponse
		err = json.Unmarshal(body, &apps)
		if err != nil {
			log.Println(err)
			return response, err
		}

		response = append(response, apps.Applications...)

		//rebuilding apiUrl
		if strings.Contains(apps.NextLink, "directoryObjects") {
			apiUrl = c.MainUrl + c.Tenant + "/" + apps.NextLink + "&$top=999&api-version=" + c.ApiVersion
		} else {
			apiUrl = apps.NextLink
		}

	}

	return response, nil
}

// TODO: deal with recursive groups
func (c *AzureClient) GetGroupAndMembersAdGraph(groupId string) (adgraph.Group, []string, error) {

	var (
		group   adgraph.Group
		members []string
	)

	apiUrl := c.MainUrl + c.Tenant + "/groups/" + groupId + "?api-version=" + c.ApiVersion

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println(err)
		return group, members, err
	}

	// Add the Bearer token to the request header
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return group, members, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return group, members, err
	}

	if resp.StatusCode != 200 {
		return group, members, errors.New("Return code is :" + fmt.Sprint(resp.StatusCode))
	}

	err = json.Unmarshal(body, &group)
	if err != nil {
		log.Println(err)
		return group, members, err
	}

	//get members

	// check cache first
	val, ok := cache.GroupMembers[groupId]
	if ok {
		return group, val, nil
	}

	apiUrl = c.MainUrl + c.Tenant + "/groups/" + groupId + "/members?$top=999&api-version=" + c.ApiVersion

	for apiUrl != "" {

		req, err = http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Println(err)
			return group, members, err
		}

		// Add the Bearer token to the request header
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)

		resp, err = c.HttpClient.Do(req)
		if err != nil {
			log.Println(err)
			return group, members, err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return group, members, err
		}

		if resp.StatusCode != 200 {
			return group, members, errors.New("Return code is :" + fmt.Sprint(resp.StatusCode))
		}

		var userResponse adgraph.UserResponse
		err = json.Unmarshal(body, &userResponse)
		if err != nil {
			log.Println(err)
			return group, members, err
		}

		for _, member := range userResponse.Users {
			members = append(members, member.ObjectId)
		}

		//rebuilding apiUrl
		if strings.Contains(userResponse.NextLink, "directoryObjects") {
			apiUrl = c.MainUrl + c.Tenant + "/" + userResponse.NextLink + "&$top=999&api-version=" + c.ApiVersion
		} else {
			apiUrl = userResponse.NextLink
		}
	}

	cache.GroupMembers[groupId] = append(cache.GroupMembers[groupId], members...)

	return group, members, nil
}

// TODO: deal with recursive groups within the role
func (c *AzureClient) GetRoleMembersAdGraph(roleId string) ([]string, error) {

	var (
		roleMembers []string
	)

	// check cache first
	val, ok := cache.RoleMembers[roleId]
	if ok {
		return val, nil
	}

	apiUrl := c.MainUrl + c.Tenant + "/roleAssignments/?$filter=roleDefinitionId%20eq%20%27" + roleId + "%27&$top=999&api-version=" + c.ApiVersion

	for apiUrl != "" {
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			log.Println(err)
			return roleMembers, err
		}

		// Add the Bearer token to the request header
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)

		resp, err := c.HttpClient.Do(req)
		if err != nil {
			log.Println(err)
			return roleMembers, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return roleMembers, err
		}

		if resp.StatusCode != 200 {
			return roleMembers, errors.New("Return code is :" + fmt.Sprint(resp.StatusCode))
		}

		var roleAssignmentResponse adgraph.RoleAssignmentResponse
		err = json.Unmarshal(body, &roleAssignmentResponse)
		if err != nil {
			log.Println(err)
			return roleMembers, err
		}

		for _, member := range roleAssignmentResponse.RoleAssignments {
			roleMembers = append(roleMembers, member.PrincipalId)
		}

		//rebuild apiUrl
		if strings.Contains(roleAssignmentResponse.NextLink, "roleAssignments") {
			apiUrl = c.MainUrl + c.Tenant + "/" + roleAssignmentResponse.NextLink + "&$top=999&api-version=" + c.ApiVersion
		} else {
			apiUrl = ""
		}
	}

	cache.RoleMembers[roleId] = append(cache.RoleMembers[roleId], roleMembers...)

	return roleMembers, nil
}
