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
		save           bool
		load           bool
		tenantId       string
		verboseLogging bool
		logToFile      string
	)
	flag.StringVar(&accessToken, "accessToken", "", "JWT access token for the specified scope")
	flag.StringVar(&userId, "userId", "", "User ObjectId for which to check gaps")
	flag.StringVar(&appId, "appId", "", "Application ID for which to check gaps")
	flag.BoolVar(&aadGraph, "aad", true, "Whether to use AAD Graph or MS Graph - current default is AAD Graph")
	flag.BoolVar(&msGraph, "msgraph", false, "Whether to use AAD Graph or MS Graph - current default is AAD Graph")
	flag.BoolVar(&save, "save", false, "If enabled, saves the conditional access policies, users, apps and locations to file(JSON format) - useful during testing")
	flag.BoolVar(&load, "load", false, "If present, conditional access policies, users, apps and locations will be loaded from the file given(JSON format)")
	flag.StringVar(&tenantId, "tenant", "", "Specify tenant ID ")
	flag.StringVar(&logToFile, "log", "", "Specify log filename to log to instead of STDOUT")
	flag.BoolVar(&verboseLogging, "v", false, "Verbose logging")
	flag.Usage = PrintUsage
	flag.Parse()
	//check params
	if save && load {
		return errors.New("Cannot save and load resources at the same time")
	}
	if accessToken == "" {
		return errors.New("Access token needs to be provided")
	}
	if load {
		if _, err := os.Stat(settings.CapsFile); err != nil {
			return fmt.Errorf("File provided("+settings.CapsFile+") does not exist: [%w]", err)
		}
		if _, err := os.Stat(settings.AppsFile); err != nil {
			return fmt.Errorf("File provided("+settings.AppsFile+") does not exist: [%w]", err)
		}
		if _, err := os.Stat(settings.UserFile); err != nil {
			return fmt.Errorf("File provided("+settings.UserFile+") does not exist: [%w]", err)
		}
		if _, err := os.Stat(settings.LocationsFile); err != nil {
			return fmt.Errorf("File provided("+settings.LocationsFile+") does not exist: [%w]", err)
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
	settings.Config[settings.LOGFILE] = logToFile
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
	if save {
		settings.Config[settings.RESOURCE_DIRECTION] = settings.SAVE
	}
	if load {
		settings.Config[settings.RESOURCE_DIRECTION] = settings.LOAD
	}

	return nil
}

func RunCapGap() {

	caps, err := parsers.ParseConditionalAccessPolicyList()
	if err != nil {
		fmt.Println("Could not retrieve conditional access policies: " + err.Error())
		return
	}
	if settings.Config[settings.USERID] != "" && settings.Config[settings.APPID] != "" {
		capgap.FindGapsPerUserAndApp(caps, settings.Config[settings.USERID], settings.Config[settings.APPID])
	} else if settings.Config[settings.APPID] == "" {
		userGaps, err := capgap.FindGapsForUser(caps, settings.Config[settings.USERID])
		if err != nil {
			fmt.Println("Could not find common bypasses for user: " + err.Error())
			return
		}
		sortedGaps := capgap.SortBypassesByAppId(userGaps)
		fmt.Println(sortedGaps)
	} else if settings.Config[settings.USERID] == "" {
		appGaps, err := capgap.FindGapsForApp(caps, settings.Config[settings.APPID])
		if err != nil {
			fmt.Println("Could not find common bypasses for app: " + err.Error())
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

	defer settings.EndLogging()

	RunCapGap()

}
