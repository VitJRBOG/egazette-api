package main

import (
	"log"
	"runtime/debug"

	config "bitbucket.org/VitJRBOG/rss_maker/internal/config"
	server "bitbucket.org/VitJRBOG/rss_maker/internal/server"
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
