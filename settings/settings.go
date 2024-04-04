package settings

var (
	Config = make(map[string]string)
)

const (
	ACCESSTOKEN       string = "AccessToken"
	USERID            string = "UserId"
	APPID             string = "ApplicationId"
	CLIENTENDPOINT    string = "ClientEndpoint"
	AADGRAPH          string = "AADGraph"
	MSGRAPH           string = "MSGraph"
	CAPFILE           string = "CapFile"
	APPFILE           string = "AppFile"
	CAPFILE_DIRECTION string = "CapDirection"
	APPFILE_DIRECTION string = "AppDirection"
	LOAD              string = "load"
	SAVE              string = "save"
	TENANT            string = "tenant"
	VERBOSE           string = "verbose"
	VERBOSE_ON        string = "ON"
)
