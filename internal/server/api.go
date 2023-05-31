package server

import (
	"egazette-api/internal/rss"
	"egazette-api/internal/sources/jpl"
	"egazette-api/internal/sources/vestirama"
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

func handling() {
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
			description := "JPL is a research and development \n" +
				"lab federally funded by NASA and managed by Caltech.\n\n"
			sourceURL := "https://www.jpl.nasa.gov/news"
			// FIXME: remove hardcoded info

			sendData(w, http.StatusOK, []map[string]string{{
				"description": description,
				"url":         sourceURL,
			}})
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/jpl/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rssFeed, err := jplArticles()

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
			description := "Оренбургская государственная телерадиовещательная компания."
			sourceURL := "https://vestirama.ru/novosti/"
			// FIXME: remove hardcoded info

			sendData(w, http.StatusOK, []map[string]string{{
				"description": description,
				"url":         sourceURL,
			}})

		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
			return
		}
	})

	http.HandleFunc("/source/vestirama/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rssFeed, err := vestiramaArticles()
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
		w.WriteHeader(errInfo.HTTPStatus)
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

func jplArticles() (interface{}, error) {
	data, err := jpl.ComposeRSSFeed()
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't fetch data from JPL",
		}
	}

	return data, nil
}

func vestiramaArticles() (rss.RSS, error) {
	rssFeed, err := vestirama.ComposeRSSFeed()
	if err != nil {
		log.Println(err.Error())
		return rss.RSS{}, Error{
			HTTPStatus: http.StatusInternalServerError,
			Detail:     "couldn't fetch data from Vestirama",
		}
	}

	return rssFeed, nil
}
