package main

import (
	"log"
	"runtime/debug"

	config "github.com/VitJRBOG/RSSMaker/internal/config"
	server "github.com/VitJRBOG/RSSMaker/internal/server"
)

func main() {
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
