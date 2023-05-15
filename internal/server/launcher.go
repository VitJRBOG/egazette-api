package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Up starts the server.
func Up() {
	handling()

	address := fmt.Sprintf(":%d", 8080) // FIXME: remove hardcoded port
	err := http.ListenAndServe(address, logging(http.DefaultServeMux))

	if err != nil {
		log.Fatalf("server launch error: %s", err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begins := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begins)

		log.Printf("[%s] %s %s", r.Method, r.RequestURI, timeElapsed)
	})
}
