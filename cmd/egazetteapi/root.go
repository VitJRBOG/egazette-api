package egazetteapi

import (
	"egazette-api/internal/config"
	"egazette-api/internal/server"
	"log"
)

// Execute starts the main functions of the program.
func Execute() {
	serverCfg, err := config.NewServerCfg()
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	server.Up(serverCfg)
}
