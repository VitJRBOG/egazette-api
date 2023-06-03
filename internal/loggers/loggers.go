package loggers

import (
	"log"
	"os"
)

// InitializeDefaultLogger sets the parameters for a standard logger.
func InitializeDefaultLogger() {
	log.SetOutput(os.Stderr)
	log.SetPrefix("[WARNING] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

// NewInfoLogger sets the parameters for the info messages logger.
func NewInfoLogger() *log.Logger {
	return log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
}
