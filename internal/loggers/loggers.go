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

// NewHTTPLogger sets the parameters for the HTTP requests logger.
func NewHTTPLogger() *log.Logger {
	return log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
}
