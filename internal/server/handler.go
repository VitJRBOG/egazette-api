package server

import (
	"database/sql"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	mux "github.com/gorilla/mux"
)

func runHandling(rtr *mux.Router, dbase *sql.DB) {
	rtr.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello world!"))
	}).Methods("GET", "POST")

	rtr.HandleFunc("/getRSSFeed/{id:[0-9]}", func(rw http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Printf("%s\n%s\n", err.Error(), debug.Stack())
		} else {
			data, err := getRSSFeed(dbase, id)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				log.Printf("%s\n%s\n", err.Error(), debug.Stack())
			}

			rw.Header().Set("Content-Type", "application/xml")
			rw.Write(data)
		}
	}).Methods("GET", "POST")
}
