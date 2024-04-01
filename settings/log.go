package settings

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogging() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ltime|log.Lshortfile)
	if Config[VERBOSE] != VERBOSE_ON {
		InfoLogger.SetOutput(io.Discard)
	}
}
