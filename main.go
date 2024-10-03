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
	fmt.Println("CapGap is meant to discover Azure Conditional Access Policy bypasses for given combinations.")
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
		reportToFile   string
		force          bool
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
	flag.StringVar(&reportToFile, "report", "", "Specify report filename to write results to instead of STDOUT")
	flag.BoolVar(&verboseLogging, "v", false, "Verbose logging")
	flag.Usage = PrintUsage
	flag.Parse()
	//check params
	if save && load {
		return errors.New("Cannot save and load resources at the same time")
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
	} else {
		if accessToken == "" {
			return errors.New("Access token needs to be provided")
		}
		if aadGraph && msGraph {
			return errors.New("You need to choose between AADGraph and MSGraph")
		}
		if tenantId == "" {
			return errors.New("Tenant needs to be specified")
		}
	}

	if verboseLogging {
		settings.Config[settings.VERBOSE] = settings.VERBOSE_ON
	}
	settings.Config[settings.LOGFILE] = logToFile
	settings.Config[settings.REPORTFILE] = reportToFile
	settings.InitLogging()
	settings.Config[settings.ACCESSTOKEN] = accessToken
	settings.Config[settings.USERID] = userId
	settings.Config[settings.APPID] = appId
	settings.Config[settings.TENANT] = tenantId
	if msGraph {
		settings.DebugLogger.Println("Using MSGraph Client")
		settings.Config[settings.CLIENTENDPOINT] = settings.MSGRAPH
	} else {
		settings.DebugLogger.Println("Using AAD Graph Client")
		settings.Config[settings.CLIENTENDPOINT] = settings.AADGRAPH
	}
	if save {
		settings.Config[settings.RESOURCE_DIRECTION] = settings.SAVE
	}
	if load {
		settings.Config[settings.RESOURCE_DIRECTION] = settings.LOAD
		//if loading, parse all objects
		parsers.ParseAll()
	}
	if force {
		settings.Config[settings.FORCE_REPORT] = settings.FORCE
	}

	return nil
}

func RunCapGap() {

	settings.InfoLogger.Println("Parsing CAP list")
	err := parsers.ParseConditionalAccessPolicyList()
	if err != nil {
		settings.ErrorLogger.Fatalln("Could not parse conditional access policies: " + err.Error())
	}

	settings.InfoLogger.Println("Parsing locations")
	err = parsers.ParseLocations()
	if err != nil {
		settings.ErrorLogger.Fatalln("Could not parse locations: " + err.Error())
	}

	if settings.Config[settings.RESOURCE_DIRECTION] == settings.SAVE {
		settings.InfoLogger.Println("Saving all information to file. All other actions skipped...")
		err = parsers.ParseApplications()
		if err != nil {
			settings.ErrorLogger.Fatalln("Could not parse applications: " + err.Error())
		}
		err = parsers.ParseUsers()
		if err != nil {
			settings.ErrorLogger.Fatalln("Could not parse users: " + err.Error())
		}
		return
	}

	if settings.Config[settings.USERID] == "" {
		err = parsers.ParseUsers()
		if err != nil {
			settings.ErrorLogger.Fatalln("Could not parse users: " + err.Error())
		}
	}

	if settings.Config[settings.APPID] == "" {
		err = parsers.ParseApplications()
		if err != nil {
			settings.ErrorLogger.Fatalln("Could not parse applications: " + err.Error())
		}
	}

	if settings.Config[settings.USERID] != "" && settings.Config[settings.APPID] != "" {
		userAppGaps := capgap.FindGapsPerUserAndApp(settings.Config[settings.USERID], settings.Config[settings.APPID])
		capgap.ReportBypassesUserApp(&userAppGaps)
	} else if settings.Config[settings.APPID] == "" && settings.Config[settings.USERID] == "" {
		capgap.FindAllGaps()
	} else if settings.Config[settings.USERID] != "" {
		userGaps := capgap.FindGapsForUser(settings.Config[settings.USERID])
		settings.InfoLogger.Println("Finished finding all bypasses(" + fmt.Sprint(len(userGaps)) + "), writing the report")
		if len(userGaps) == 0 {
			settings.Reporter.WriteString("No bypasses for user " + settings.Config[settings.USERID])
		} else {
			capgap.ReportForUser(&userGaps)
		}
	} else if settings.Config[settings.APPID] != "" {
		appGaps := capgap.FindGapsForApp(settings.Config[settings.APPID])
		settings.InfoLogger.Println("Finished finding all bypasses(" + fmt.Sprint(len(appGaps)) + "), writing the report")
		if len(appGaps) == 0 {
			settings.Reporter.WriteString("No bypasses for app " + settings.Config[settings.APPID])
		} else {
			capgap.ReportForApp(&appGaps)
		}
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
