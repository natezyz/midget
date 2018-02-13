package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/natezyz/midget/storage"
)

type Request struct {
	Action string
	Id     string
	Params map[string]string
}

func encode(url string, w http.ResponseWriter, storage storage.Storage) {
	w.Write([]byte(storage.Store(url)))
}

func decode(encoded string, w http.ResponseWriter, storage storage.Storage) {
	url, err := storage.Retrieve(encoded)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error: URL not found: " + err.Error() + "\n"))
		return
	}
	log.Printf("Process(): Decoded %s to %s\n", encoded, url)
	w.Write([]byte(url))
}

func Process(storage storage.Storage) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req Request
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("Process(): Parse error: %v\n", err)
			return
		}
		defer r.Body.Close()

		switch req.Action {
		case "encode":
			encode(req.Params["url"], w, storage)
		case "decode":
			decode(req.Params["code"], w, storage)
		default:
			log.Printf("Process(): Action %s not supported\n", req.Action)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Action " + req.Action + " not supported\n"))
		}
	}

	return http.HandlerFunc(handler)
}

func Redirect(storage storage.Storage) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/"):]
		log.Printf("Redirect(): Found code %s\n", code)

		url, err := storage.Retrieve(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Error: URL not found: " + err.Error() + "\n"))
			return
		}
		log.Printf("Redirecting to url %s\n", url)

		http.Redirect(w, r, url, 301)
	}

	return http.HandlerFunc(handler)
}
