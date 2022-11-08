package config

import (
	"log"
	"os"
)

type DBConnCfg struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDBConnCfg() DBConnCfg {
	host := os.Getenv("DBMS_HOST")
	port := os.Getenv("DBMS_PORT")
	user := os.Getenv("DBMS_USER")
	password := os.Getenv("DBMS_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("SSL_MODE")

	someEmpty := false

	if host == "" {
		someEmpty = true
		log.Println("DBMS_HOST env variable is empty")
	}

	if port == "" {
		someEmpty = true
		log.Println("DBMS_PORT env variable is empty")
	}

	if user == "" {
		someEmpty = true
		log.Println("DBMS_USER env variable is empty")
	}

	if password == "" {
		someEmpty = true
		log.Println("DBMS_PASSWORD env variable is empty")
	}

	if dbName == "" {
		someEmpty = true
		log.Println("DB_NAME env variable is empty")
	}

	if sslMode == "" {
		someEmpty = true
		log.Println("SSL_MODE env variable is empty")
	}

	if someEmpty {
		log.Fatalln("some env variable is empty")
	}

	return DBConnCfg{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}
}

type ServerCfg struct {
	Port string
}

func NewServerConfig() ServerCfg {
	port := os.Getenv("SERVER_PORT")

	someEmpty := false

	if port == "" {
		someEmpty = true
		log.Println("SERVER_PORT env variable is empty")
	}

	if someEmpty {
		log.Fatalln("some env variable is empty")
	}

	return ServerCfg{
		Port: port,
	}
}
