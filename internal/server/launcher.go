package server

import (
	"fmt"
	"log"
	"net/http"
)

// Up starts the server.
func Up() {
	address := fmt.Sprintf(":%d", 8080)
	err := http.ListenAndServe(address, nil)

	if err != nil {
		log.Fatalf("server launch error: %s", err)
	}
}
