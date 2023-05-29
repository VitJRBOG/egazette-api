package config

import (
	"errors"
	"os"
)

// ServerCfg stores the parameters for launching the server.
type ServerCfg struct {
	Port string
}

// NewServerCfg receives the env variables and returns the struct with them.
func NewServerCfg() (ServerCfg, error) {
	cfg := ServerCfg{}

	cfg.Port = os.Getenv("SERVER_PORT")

	if cfg.Port == "" {
		return ServerCfg{}, errors.New("SERVER_PORT env variable is empty")
	}

	return cfg, nil
}
