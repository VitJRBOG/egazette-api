package egazetteapi

import (
	"egazette-api/internal/config"
	"egazette-api/internal/db"
	"egazette-api/internal/harvester"
	"egazette-api/internal/loggers"
	"egazette-api/internal/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	sources, err := db.SelectSources(dbConnection)
	if err != nil {
		log.Fatalf("launching is not possible: %s", err)
	}

	if len(sources) == 0 {
		log.Fatalf("launching is not possible: no sources found")
	}

	serverRepresentative, harvesterRepresentative := getCompletionHeralds()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go server.Up(&wg, serverRepresentative, serverCfg, dbConnection, sources)

	wg.Add(1)
	go harvester.Harvesting(&wg, harvesterRepresentative, dbConnection, sources)

	wg.Add(1)
	go osSignalsReception(&wg, serverRepresentative, harvesterRepresentative)

	wg.Wait()
	loggers.NewInfoLogger().Println("program exited successfully")
}

func getCompletionHeralds() (chan os.Signal, chan os.Signal) {
	serverRepresentative := make(chan os.Signal, 1)
	harvesterRepresentative := make(chan os.Signal, 1)

	return serverRepresentative, harvesterRepresentative
}

func osSignalsReception(wg *sync.WaitGroup, serverRepresentative, harvesterRepresentative chan os.Signal) {
	signal.Notify(serverRepresentative, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(harvesterRepresentative, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	wg.Done()
}
