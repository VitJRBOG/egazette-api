package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	config "RSSFeeder/internal/config"
	db "RSSFeeder/internal/db"

	mux "github.com/gorilla/mux"
)

func Start(dbConn config.DBConn, serverCfg config.ServerCfg) {
	dbase, err := connectToDB(dbConn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer func(dbase *sql.DB) {
		err := dbase.Close()
		if err != nil {
			log.Printf("\n%s\n%s", err, debug.Stack())
		}
	}(dbase)

	rtr := mux.NewRouter()
	runHandling(rtr, dbase)

	http.Handle("/", rtr)
	err = http.ListenAndServe(":"+serverCfg.Port, nil)
	if err != nil {
		log.Fatalf("\n%s\n%s", err.Error(), debug.Stack())
	}
}

func connectToDB(dbConn config.DBConn) (*sql.DB, error) {
	var dbase *sql.DB
	if (dbConn != config.DBConn{}) {
		var err error
		dbase, err = db.Connect(dbConn)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("\ncouldn't connect to database: configs are empty\n%s",
			debug.Stack())
	}

	return dbase, nil
}
