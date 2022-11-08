package config

import (
	"log"
	"os"
)

type ServerCfg struct {
	Host string
	Port string
}

func NewServerConfig() ServerCfg {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	someEmpty := false

	if host == "" {
		someEmpty = true
		log.Println("SERVER_HOST env variable is empty")
	}

	if port == "" {
		someEmpty = true
		log.Println("SERVER_PORT env variable is empty")
	}

	if someEmpty {
		log.Fatalln("some env variable is empty")
	}

	return ServerCfg{
		Host: host,
		Port: port,
	}
}
