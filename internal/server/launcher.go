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
	httpLogger := loggers.NewHTTPLogger()

	handling(dbConn)
	httpLogger.Println("request handling is ready")

	address := fmt.Sprintf(":%s", serverCfg.Port)
	err := http.ListenAndServe(address, logging(http.DefaultServeMux, httpLogger))

	if err != nil {
		log.Fatalf("server launch error: %s", err)
	}
}

func logging(next http.Handler, httpLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begins := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begins)

		httpLogger.Printf("[%s] %s %s", r.Method, r.RequestURI, timeElapsed)
	})
}
