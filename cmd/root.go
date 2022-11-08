package cmd

import (
	config "RSSFeeder/internal/config"
	server "RSSFeeder/internal/server"
)

func Execute() {
	serverCfg := config.NewServerConfig()

	server.Start(serverCfg)
}
