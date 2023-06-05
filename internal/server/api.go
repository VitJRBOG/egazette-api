package server

import (
	"egazette-api/internal/db"
	"egazette-api/internal/models"
	"egazette-api/internal/rss"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

// Error stores info about error.
type Error struct {
	HTTPStatus int
	Detail     string
}

// Error returns a text representation of error info.
func (e Error) Error() string {
	return fmt.Sprintf("status %d: %s", e.HTTPStatus, e.Detail)
}

func handling(dbConn db.Connection, sources []models.Source) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			text := "/source - list of sources"

			sendText(w, http.StatusOK, text)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			text := fmt.Sprintf("%s/source/jpl\n%s/source/vestirama",
				r.Host, r.Host)

			sendText(w, http.StatusOK, text)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/jpl", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			text := fmt.Sprintf("%s/source/jpl/info - about source\n%s/source/jpl/articles - list of articles",
				r.Host, r.Host)

			sendText(w, http.StatusOK, text)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/jpl/info", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			source, err := jplInfo(sources)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusOK, []map[string]string{{
				"name":        source.Name,
				"description": source.Description,
				"url":         source.HomeURL,
			}})
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/jpl/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rssFeed, err := jplArticles(dbConn, sources)

			if err != nil {
				sendError(w, err)
				return
			}

			sendRSSFeed(w, rssFeed)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/vestirama", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			text := fmt.Sprintf("%s/source/vestirama/info - about source\n%s/source/vestirama/articles - list of articles",
				r.Host, r.Host)

			sendText(w, http.StatusOK, text)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/vestirama/info", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			source, err := vestiramaInfo(sources)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusOK, []map[string]string{{
				"name":        source.Name,
				"description": source.Description,
				"url":         source.HomeURL,
			}})
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/vestirama/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rssFeed, err := vestiramaArticles(dbConn, sources)
			if err != nil {
				sendError(w, err)
				return
			}

			sendRSSFeed(w, rssFeed)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/natgeo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			text := fmt.Sprintf("%s/source/natgeo/info - about source\n%s/source/natgeo/articles - list of articles",
				r.Host, r.Host)

			sendText(w, http.StatusOK, text)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/natgeo/info", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			source, err := natgeoInfo(sources)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, http.StatusOK, []map[string]string{{
				"name":        source.Name,
				"description": source.Description,
				"url":         source.HomeURL,
			}})
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/natgeo/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rssFeed, err := natgeoArticles(dbConn, sources)
			if err != nil {
				sendError(w, err)
				return
			}

			sendRSSFeed(w, rssFeed)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})
}

func sendText(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(text))
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}
}

func sendRSSFeed(w http.ResponseWriter, values interface{}) {
	data, err := xml.Marshal(values)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}
}

func sendData(w http.ResponseWriter, status int, values interface{}) {
	response := map[string]interface{}{
		"response": values,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}
}

func sendError(w http.ResponseWriter, reqError error) {
	response := map[string]interface{}{
		"error": "internal server error",
	}
	w.WriteHeader(http.StatusInternalServerError)

	if errInfo, ok := reqError.(Error); ok {
		if errInfo.HTTPStatus != http.StatusInternalServerError {
			w.WriteHeader(errInfo.HTTPStatus)
		}
		response["error"] = errInfo.Detail
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func jplInfo(sources []models.Source) (models.Source, error) {
	source := models.FindSourceByAPIName(sources, "jpl")

	if source.Name == "" && source.Description == "" && source.HomeURL == "" && source.APIName == "" {
		return models.Source{}, Error{
			HTTPStatus: http.StatusNoContent,
			Detail:     "couldn't find data about JPL",
		}
	}

	return source, nil
}

func vestiramaInfo(sources []models.Source) (models.Source, error) {
	source := models.FindSourceByAPIName(sources, "vestirama")

	if source.Name == "" && source.Description == "" && source.HomeURL == "" && source.APIName == "" {
		return models.Source{}, Error{
			HTTPStatus: http.StatusNoContent,
			Detail:     "couldn't find data about Vestirama",
		}
	}

	return source, nil
}

func natgeoInfo(sources []models.Source) (models.Source, error) {
	source := models.FindSourceByAPIName(sources, "natgeo")

	if source.Name == "" && source.Description == "" && source.HomeURL == "" && source.APIName == "" {
		return models.Source{}, Error{
			HTTPStatus: http.StatusNoContent,
			Detail:     "couldn't find data about NatGeo",
		}
	}

	return source, nil
}

func jplArticles(dbConn db.Connection, sources []models.Source) (rss.RSS, error) {
	source := models.FindSourceByAPIName(sources, "jpl")

	if source.Name == "" && source.Description == "" && source.HomeURL == "" && source.APIName == "" {
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusNoContent,
			Detail:     "couldn't find data about JPL",
		}
	}

	articles, err := db.SelectArticles(dbConn, source.Name)
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't fetch data with JPL articles",
		}
	}

	rssFeed, err := rss.ComposeRSSFeed(source, articles)
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't compose RSS feed with JPL articles",
		}
	}

	return rssFeed, nil
}

func vestiramaArticles(dbConn db.Connection, sources []models.Source) (rss.RSS, error) {
	source := models.FindSourceByAPIName(sources, "vestirama")

	if source.Name == "" && source.Description == "" && source.HomeURL == "" && source.APIName == "" {
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusNoContent,
			Detail:     "couldn't find data about Vestirama",
		}
	}

	articles, err := db.SelectArticles(dbConn, source.Name)
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't fetch data with Vestirama articles",
		}
	}

	rssFeed, err := rss.ComposeRSSFeed(source, articles)
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't compose RSS feed with Vestirama articles",
		}
	}

	return rssFeed, nil
}

func natgeoArticles(dbConn db.Connection, sources []models.Source) (rss.RSS, error) {
	source := models.FindSourceByAPIName(sources, "natgeo")

	if source.Name == "" && source.Description == "" && source.HomeURL == "" && source.APIName == "" {
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusNoContent,
			Detail:     "couldn't find data about NatGeo",
		}
	}

	articles, err := db.SelectArticles(dbConn, source.Name)
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't fetch data with NatGeo articles",
		}
	}

	rssFeed, err := rss.ComposeRSSFeed(source, articles)
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't compose RSS feed with NatGeo articles",
		}
	}

	return rssFeed, nil
}
