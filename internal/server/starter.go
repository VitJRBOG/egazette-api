package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	config "github.com/VitJRBOG/RSSMaker/internal/config"
	db "github.com/VitJRBOG/RSSMaker/internal/db"
	mux "github.com/gorilla/mux"
)

func Start(dbConn config.DBConn, serverCfg config.ServerCfg) {
	dbase, err := connectToDB(dbConn)
	if err != nil {
		log.Fatalf(err.Error(), debug.Stack())
	}

	defer func(dbase *sql.DB) {
		err := dbase.Close()
		if err != nil {
			log.Printf("%s\n%s\n", err, debug.Stack())
		}
	}(dbase)

	rtr := mux.NewRouter()
	runHandling(rtr, dbase)

	http.Handle("/", rtr)
	err = http.ListenAndServe(":"+serverCfg.Port, nil)
	if err != nil {
		log.Fatalf(err.Error(), debug.Stack())
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
		return nil, fmt.Errorf("couldn't connect to database: configs are empty")
	}

	return dbase, nil
}
