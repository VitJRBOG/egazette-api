package egazetteapi

import (
	"egazette-api/internal/config"
	"egazette-api/internal/db"
	"egazette-api/internal/harvester"
	"egazette-api/internal/loggers"
	"egazette-api/internal/server"
	"fmt"
	"log"
)

// Execute starts the main functions of the program.
func Execute() {
	loggers.InitializeDefaultLogger()

	dbConnectionCfg, err := config.NewDBConnectionCfg()
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	serverCfg, err := config.NewServerCfg()
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConnectionCfg.User, dbConnectionCfg.Password,
		dbConnectionCfg.HostAddress, dbConnectionCfg.HostPort,
		dbConnectionCfg.DBName,
		dbConnectionCfg.SSLMode)

	dbConnection, err := db.NewConnection(dsn)
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	// FIXME: It's crud. Need to describe an OS signals listener.

	go server.Up(serverCfg, dbConnection)

	harvester.Harvesting(dbConnection)
}
