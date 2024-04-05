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
		accessToken       string
		userId            string
		appId             string
		aadGraph          bool
		msGraph           bool
		saveCapToFile     string
		loadCapFromFile   string
		saveAppToFile     string
		loadAppFromFile   string
		saveUsersToFile   string
		loadUsersFromFile string
		tenantId          string
		verboseLogging    bool
	)
	flag.StringVar(&accessToken, "accessToken", "", "JWT access token for the specified scope")
	flag.StringVar(&userId, "userId", "", "User ObjectId for which to check gaps")
	flag.StringVar(&appId, "appId", "", "Application ID for which to check gaps")
	flag.BoolVar(&aadGraph, "aad", true, "Whether to use AAD Graph or MS Graph - current default is AAD Graph")
	flag.BoolVar(&msGraph, "msgraph", false, "Whether to use AAD Graph or MS Graph - current default is AAD Graph")
	flag.StringVar(&saveCapToFile, "save", "", "If enabled, saves the conditional access policies to file(JSON format) - useful during testing")
	flag.StringVar(&loadCapFromFile, "loadCaps", "", "If present, conditional access policies will be loaded from the file given(JSON format)")
	flag.StringVar(&saveAppToFile, "saveApps", "", "If enabled, saves the applications to file(JSON format) - useful during testing")
	flag.StringVar(&loadAppFromFile, "loadApps", "", "If present, applications will be loaded from the file given(JSON format)")
	flag.StringVar(&saveUsersToFile, "saveUsers", "", "If enabled, saves the users to file(JSON format) - useful during testing")
	flag.StringVar(&loadUsersFromFile, "loadUsers", "", "If present, users will be loaded from the file given(JSON format)")
	flag.StringVar(&tenantId, "tenant", "", "Specify tenant ID ")
	flag.BoolVar(&verboseLogging, "v", false, "Verbose logging")
	flag.Usage = PrintUsage
	flag.Parse()
	//check params
	if saveCapToFile != "" && loadCapFromFile != "" {
		return errors.New("Cannot save and load conditional access policies at the same time")
	}
	if saveAppToFile != "" && loadAppFromFile != "" {
		return errors.New("Cannot save and load apps at the same time")
	}
	if saveUsersToFile != "" && loadUsersFromFile != "" {
		return errors.New("Cannot save and load users at the same time")
	}
	if accessToken == "" {
		return errors.New("Access token needs to be provided")
	}
	if loadCapFromFile != "" {
		if _, err := os.Stat(loadCapFromFile); err != nil {
			return fmt.Errorf("File provided("+loadCapFromFile+") does not exist: [%w]", err)
		}
	}
	if loadAppFromFile != "" {
		if _, err := os.Stat(loadAppFromFile); err != nil {
			return fmt.Errorf("File provided("+loadAppFromFile+") does not exist: [%w]", err)
		}
	}
	if loadUsersFromFile != "" {
		if _, err := os.Stat(loadUsersFromFile); err != nil {
			return fmt.Errorf("File provided("+loadUsersFromFile+") does not exist: [%w]", err)
		}
	}
	if aadGraph && msGraph {
		return errors.New("You need to choose between AADGraph and MSGraph")
	}
	if tenantId == "" {
		return errors.New("Tenant needs to be specified")
	}
	if verboseLogging {
		settings.Config[settings.VERBOSE] = settings.VERBOSE_ON
	}
	settings.InitLogging()
	settings.Config[settings.ACCESSTOKEN] = accessToken
	settings.Config[settings.USERID] = userId
	settings.Config[settings.APPID] = appId
	settings.Config[settings.TENANT] = tenantId
	if msGraph {
		settings.InfoLogger.Println("Using MSGraph Client")
		settings.Config[settings.CLIENTENDPOINT] = settings.MSGRAPH
	} else {
		settings.InfoLogger.Println("Using AAD Graph Client")
		settings.Config[settings.CLIENTENDPOINT] = settings.AADGRAPH
	}
	if saveCapToFile != "" {
		settings.Config[settings.CAPFILE] = saveCapToFile
		settings.Config[settings.CAPFILE_DIRECTION] = settings.SAVE
	}
	if loadCapFromFile != "" {
		settings.Config[settings.CAPFILE] = loadCapFromFile
		settings.Config[settings.CAPFILE_DIRECTION] = settings.LOAD
	}
	if saveAppToFile != "" {
		settings.Config[settings.APPFILE] = saveAppToFile
		settings.Config[settings.APPFILE_DIRECTION] = settings.SAVE
	}
	if loadAppFromFile != "" {
		settings.Config[settings.APPFILE] = loadAppFromFile
		settings.Config[settings.APPFILE_DIRECTION] = settings.LOAD
	}
	if saveUsersToFile != "" {
		settings.Config[settings.USERFILE] = saveUsersToFile
		settings.Config[settings.USERFILE_DIRECTION] = settings.SAVE
	}
	if loadUsersFromFile != "" {
		settings.Config[settings.USERFILE] = loadUsersFromFile
		settings.Config[settings.USERFILE_DIRECTION] = settings.LOAD
	}

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
		userGaps, err := capgap.FindGapsForUser(caps, settings.Config[settings.USERID])
		if err != nil {
			settings.ErrorLogger.Println("Could not find common bypasses for user")
			return
		}
		sortedGaps := capgap.SortBypassesByAppId(userGaps)
		fmt.Println(sortedGaps)
	} else if settings.Config[settings.USERID] == "" {
		appGaps, err := capgap.FindGapsForApp(caps, settings.Config[settings.APPID])
		if err != nil {
			settings.ErrorLogger.Println("Could not find common bypasses for app")
			return
		}
		sortedGaps := capgap.SortBypassesByUserId(appGaps)
		fmt.Println(sortedGaps)
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
