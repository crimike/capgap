package main

// azuread graph or microsoft graph
// login

import (
	"capgap/capgap"
	"capgap/parsers"
	"capgap/settings"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

func PrintUsage() {
	fmt.Println("CapGap is meant to discover Azure Conditional Access Policy bypasses for certain combinations.")
	fmt.Println()
	flag.PrintDefaults()
}

func ParseCommandLine() error {
	var (
		accessToken    string
		userId         string
		appId          string
		aadGraph       bool
		msGraph        bool
		saveToFile     string
		loadFromFile   string
		tenantId       string
		verboseLogging bool
	)
	flag.StringVar(&accessToken, "accessToken", "", "JWT access token for the specified scope")
	flag.StringVar(&userId, "userId", "", "User ObjectId for which to check gaps")
	flag.StringVar(&appId, "appId", "", "Application ID for which to check gaps")
	flag.BoolVar(&aadGraph, "aad", true, "Whether to use AAD Graph or MS Graph - current default is AAD Graph")
	flag.BoolVar(&msGraph, "msgraph", false, "Whether to use AAD Graph or MS Graph - current default is AAD Graph")
	flag.StringVar(&saveToFile, "save", "", "If enabled, saves the conditional access policies to file(JSON format) - useful during testing")
	flag.StringVar(&loadFromFile, "load", "", "If present, conditional access policies will be loaded from the file given(JSON format)")
	flag.StringVar(&tenantId, "tenant", "", "Specify tenant ID ")
	flag.BoolVar(&verboseLogging, "v", false, "Verbose logging")
	flag.Usage = PrintUsage
	flag.Parse()
	//check params
	if saveToFile != "" && loadFromFile != "" {
		return errors.New("Cannot save and load conditional access policies at the same time")
	}
	if accessToken == "" {
		return errors.New("Access token needs to be provided")
	}
	if loadFromFile != "" {
		if _, err := os.Stat(loadFromFile); err != nil {
			return fmt.Errorf("File provided("+loadFromFile+") does not exist: [%w]", err)
		}
	}
	if aadGraph && msGraph {
		return errors.New("You need to choose between AADGraph and MSGraph")
	}
	if userId == "" {
		return errors.New("User ID needs to be specified")
	}
	if tenantId == "" {
		return errors.New("Tenant needs to be specified")
	}
	settings.Config[settings.ACCESSTOKEN] = accessToken
	settings.Config[settings.USERID] = userId
	settings.Config[settings.APPID] = appId
	settings.Config[settings.TENANT] = tenantId
	if msGraph {
		settings.Config[settings.CLIENTENDPOINT] = settings.MSGRAPH
	} else {
		settings.Config[settings.CLIENTENDPOINT] = settings.AADGRAPH
	}
	if saveToFile != "" {
		settings.Config[settings.CAPFILE] = saveToFile
		settings.Config[settings.CAPFILE_DIRECTION] = settings.SAVECAP
	}
	if loadFromFile != "" {
		settings.Config[settings.CAPFILE] = loadFromFile
		settings.Config[settings.CAPFILE_DIRECTION] = settings.LOADCAP
	}
	if verboseLogging {
		settings.Config[settings.VERBOSE] = settings.VERBOSE_ON
	}
	settings.InitLogging()

	return nil
}

func RunCapGap() {

	caps, err := parsers.ParseConditionalAccessPolicyList()
	if err != nil {
		settings.ErrorLogger.Fatalln("Could not retrieve conditional access policies: " + err.Error())
	}
	if settings.Config[settings.USERID] != "" && settings.Config[settings.APPID] != "" {
		capgap.FindGapsPerUserAndApp(caps, settings.Config[settings.USERID], settings.Config[settings.APPID])
	} else if settings.Config[settings.APPID] == "" {
		capgap.FindGapsPerUser(caps, settings.Config[settings.USERID])
	}
}

func main() {

	err := ParseCommandLine()
	if err != nil {
		flag.Usage()
		log.Panicln(err)
	}

	RunCapGap()

}
