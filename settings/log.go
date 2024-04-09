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
	logFile     *os.File
	Reporter    *os.File
)

func InitLogging() {
	var f *os.File

	if Config[LOGFILE] != "" {
		f, err := os.OpenFile(Config[LOGFILE], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("Could not open log file for writing")
		}
		logFile = f
	} else {
		f = os.Stdout
	}
	if Config[REPORTFILE] != "" {
		rep, err := os.OpenFile(Config[REPORTFILE], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("Could not open report file for writing")
		}
		Reporter = rep
	} else {
		Reporter = os.Stdout
	}

	DebugLogger = log.New(f, "DEBUG: ", log.Ltime|log.Lshortfile)
	InfoLogger = log.New(f, "INFO: ", log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "ERROR: ", log.Ltime|log.Lshortfile)
	if f != os.Stdout {
		mw := io.MultiWriter(os.Stdout, f)
		InfoLogger.SetOutput(mw)
		ErrorLogger.SetOutput(mw)
	}
	if Config[LOGFILE] == "" && Config[VERBOSE] != VERBOSE_ON {
		DebugLogger.SetOutput(io.Discard)
	}
	dt := time.Now()
	InfoLogger.Println("Starting CAPGAP at " + dt.Format("15:04:05.000000"))
}

func EndLogging() {
	dt := time.Now()
	InfoLogger.Println("Parsing finished at " + dt.Format("15:04:05.000000"))
	if Config[LOGFILE] != "" {
		logFile.Close()
	}
}
