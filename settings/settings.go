package settings

// TODO: migrate to a map[string]interface dictionary to support different types
var (
	Config = make(map[string]string)
)

const (
	ACCESSTOKEN        string = "AccessToken"
	USERID             string = "UserId"
	APPID              string = "ApplicationId"
	CLIENTENDPOINT     string = "ClientEndpoint"
	AADGRAPH           string = "AADGraph"
	MSGRAPH            string = "MSGraph"
	RESOURCE_DIRECTION string = "CapDirection"
	LOAD               string = "load"
	SAVE               string = "save"
	TENANT             string = "tenant"
	VERBOSE            string = "verbose"
	LOGFILE            string = "LogFile"
	REPORTFILE         string = "ReportFile"
	FORCE_REPORT       string = "ForceReporting"
	TRUE               string = "TRUE"
	ALL_LOCATIONS      string = "AllLocations"

	CapsFile      string = "caps.json"
	UserFile      string = "users.json"
	AppsFile      string = "apps.json"
	LocationsFile string = "locations.json"
)
