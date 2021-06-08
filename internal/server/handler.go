package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	mux "github.com/gorilla/mux"
)

func runHandling(rtr *mux.Router, dbase *sql.DB) {
	rtr.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		manual := "Методы API:\n" +
			"/getRSSFeed/[id] - где вместо [id] подставляется целочисленный идентификатор " +
			"добавленного источника, " +
			"возвращает xml с публикациями источника. " +
			"Для чтения публикаций данную ссылку можно вставить в RSS-ридер.\n" +
			"/addRSSSource - добавляет новый источник для RSS-ленты. " +
			"Возвращает JSON в формате: {\"feed_id\":[id]}, где [id] - целочисленный " +
			"идентификатор источника. " +
			"Пока что можно добавить только паблики ВК. Для этого необходимо выполнить " +
			"post-запрос с параметрами в виде JSON в формате:\n" +
			"{\n" +
			"    \"source_name\": \"SomePublic\", // название источника (на английском)\n" +
			"    \"url\": \"https://vk.com/some_public\", // ссылка на паблик или группу\n" +
			"    \"access_token\": \"q1w2e3r4t5y6u7i8o9p0\", // бессрочный токен ВК\n" +
			"    \"vk_id\": -123456, // идентификатор паблика или группы (целое число)\n" +
			"}"
		rw.Write([]byte(manual))
	}).Methods("GET", "POST")

	rtr.HandleFunc("/getRSSFeed/{id:[0-9]+}", func(rw http.ResponseWriter, r *http.Request) {
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

	rtr.HandleFunc("/addRSSSource", func(rw http.ResponseWriter, r *http.Request) {
		var vkRSSSource VKRSSSource
		err := json.NewDecoder(r.Body).Decode(&vkRSSSource)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Printf("%s\n%s\n", err.Error(), debug.Stack())
		} else {
			data, err := addVKRSSSource(dbase, vkRSSSource)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				log.Printf("%s\n%s\n", err.Error(), debug.Stack())
			} else {
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(data)
			}
		}

		defer r.Body.Close()
	}).Methods("POST")
}
