package server

import (
	"egazette-api/internal/config"
	"egazette-api/internal/db"
	"egazette-api/internal/loggers"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Up starts the server.
func Up(serverCfg config.ServerCfg, dbConn db.Connection) {
	infoLogger := loggers.NewInfoLogger()

	handling(dbConn)
	infoLogger.Println("request handling is ready")

	address := fmt.Sprintf(":%s", serverCfg.Port)
	err := http.ListenAndServe(address, logging(http.DefaultServeMux, infoLogger))

	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("server launch error: %s", err)
	}
}

func logging(next http.Handler, infoLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begins := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begins)

		infoLogger.Printf("[%s] %s %s", r.Method, r.RequestURI, timeElapsed)
	})
}
