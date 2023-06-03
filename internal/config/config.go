package config

import (
	"errors"
	"os"
)

// DBConnectionCfg stores the parameters for connecting to the Postgres DB.
type DBConnectionCfg struct {
	HostAddress string
	HostPort    string
	User        string
	Password    string
	DBName      string
	SSLMode     string
}

// NewDBConnectionCfg receives the env variables and returns the struct with them.
func NewDBConnectionCfg() (DBConnectionCfg, error) {
	cfg := DBConnectionCfg{}
	cfg.HostAddress = os.Getenv("POSTGRES_HOST_ADDRESS")
	cfg.HostPort = os.Getenv("POSTGRES_HOST_PORT")
	cfg.User = os.Getenv("POSTGRES_USER")
	cfg.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.DBName = os.Getenv("POSTGRES_DB")
	cfg.SSLMode = os.Getenv("SSL_MODE")

	someIsEmpty := false

	emptyVar := []string{}

	if cfg.HostAddress == "" {
		emptyVar = append(emptyVar, "POSTGRES_HOST_ADDRESS")
		someIsEmpty = true
	}

	if cfg.HostPort == "" {
		emptyVar = append(emptyVar, "POSTGRES_HOST_PORT")
		someIsEmpty = true
	}

	if cfg.User == "" {
		emptyVar = append(emptyVar, "POSTGRES_USER")
		someIsEmpty = true
	}

	if cfg.Password == "" {
		emptyVar = append(emptyVar, "POSTGRES_PASSWORD")
		someIsEmpty = true
	}

	if cfg.DBName == "" {
		emptyVar = append(emptyVar, "POSTGRES_DB")
		someIsEmpty = true
	}

	if cfg.SSLMode == "" {
		emptyVar = append(emptyVar, "SSL_MODE")
		someIsEmpty = true
	}

	if someIsEmpty {
		errMsg := "some desktop environments is empty: "
		for i := range emptyVar {
			if i > 0 && i < len(emptyVar) {
				errMsg += ", "
			}
			errMsg += emptyVar[i]
		}
		return DBConnectionCfg{}, errors.New(errMsg)
	}

	return cfg, nil
}

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
