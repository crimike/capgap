package settings

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	logFile     *os.File
)

func InitLogging() {
	var f *os.File

	if Config[LOGFILE] != "" {
		f, _ = os.OpenFile(Config[LOGFILE], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		logFile = f
	} else {
		f = os.Stdout
	}

	InfoLogger = log.New(f, "INFO: ", log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "ERROR: ", log.Ltime|log.Lshortfile)
	if Config[LOGFILE] == "" && Config[VERBOSE] != VERBOSE_ON {
		InfoLogger.SetOutput(io.Discard)
	}
}

func EndLogging() {
	if Config[LOGFILE] != "" {
		logFile.Close()
	}
}
