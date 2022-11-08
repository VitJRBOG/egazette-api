package cmd

import (
	"log"
	"runtime/debug"

	config "RSSFeeder/internal/config"
	server "RSSFeeder/internal/server"
)

func Execute() {
	dbConn, err := config.GetDBConnectionData()
	if err != nil {
		log.Fatalf("%s\n%s\n", err.Error(), debug.Stack())
	}

	serverCfg, err := config.GetServerConfig()
	if err != nil {
		log.Fatalf("%s\n%s\n", err.Error(), debug.Stack())
	}

	server.Start(dbConn, serverCfg)
}
