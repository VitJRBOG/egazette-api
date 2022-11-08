package server

import (
	"fmt"
	"log"
	"net/http"

	config "RSSFeeder/internal/config"
)

func Start(serverCfg config.ServerCfg) {
	handling()

	address := fmt.Sprintf("%s:%s", serverCfg.Host, serverCfg.Port)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatalf("server error: %s", err)
	}
}
