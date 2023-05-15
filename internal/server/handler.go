package server

import (
	"log"
	"net/http"
)

func handling() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			log.Println(err.Error())
		}
	})
}
