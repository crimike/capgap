package settings

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	LogFile     *os.File
	Reporter    *os.File
)

func InitLogging() {
	var writer io.Writer

	if Config[LOGFILE] != "" {
		f, err := os.OpenFile(Config[LOGFILE].(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("Could not open log file for writing")
		}
		LogFile = f
		writer = io.MultiWriter(os.Stdout, LogFile)
	} else {
		writer = os.Stdout
	}
	if Config[REPORTFILE] != "" {
		rep, err := os.OpenFile(Config[REPORTFILE].(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("Could not open report file for writing")
		}
		Reporter = rep
	} else {
		Reporter = os.Stdout
	}

	InfoLogger = log.New(writer, "INFO: ", log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(writer, "ERROR: ", log.Ltime|log.Lshortfile)

	if Config[VERBOSE].(bool) {
		DebugLogger = log.New(writer, "DEBUG: ", log.Ltime|log.Lshortfile)
	} else if Config[LOGFILE] != "" {
		DebugLogger = log.New(LogFile, "DEBUG: ", log.Ltime|log.Lshortfile)
	} else {
		DebugLogger = log.New(io.Discard, "DEBUG: ", log.Ltime|log.Lshortfile)
	}

	dt := time.Now()
	InfoLogger.Println("Starting CAPGAP at " + dt.Format("15:04:05.000000"))
	DebugLogger.Println("debug test")
}

func EndLogging() {
	dt := time.Now()
	InfoLogger.Println("Parsing finished at " + dt.Format("15:04:05.000000"))
	if Config[LOGFILE] != "" {
		// 	fmt.Println("Closing log file")
		LogFile.Sync()
		LogFile.Close()
	}
	if Config[REPORTFILE] != "" {
		Reporter.Sync()
		Reporter.Close()
	}
}
