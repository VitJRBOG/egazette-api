package cmd

import (
	config "RSSFeeder/internal/config"
	server "RSSFeeder/internal/server"
)

func Execute() {
	dbConnCfg := config.NewDBConnCfg()
	serverCfg := config.NewServerConfig()

	server.Start(dbConnCfg, serverCfg)
}
