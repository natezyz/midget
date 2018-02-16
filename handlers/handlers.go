package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/natezyz/midget/storage"
)

type EncodeRequest struct {
	Url string `json:"url"`
}

type MidgetResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Encode(storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			respond(&MidgetResponse{
				Status:  404,
				Message: "method not found",
			}, w)
			return
		}

		log.Printf("Method is %s\n", r.Method)

		decoder := json.NewDecoder(r.Body)
		var req EncodeRequest
		if err := decoder.Decode(&req); err != nil {
			log.Printf("Encode(): Parse error: %v\n", err)
			return
		}
		defer r.Body.Close()

		respond(&MidgetResponse{
			Status:  200,
			Message: storage.Store(req.Url),
		}, w)
	})
}

func Decode(storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respond(&MidgetResponse{
				Status:  404,
				Message: "method not found",
			}, w)
			return
		}
		code := r.URL.Path[len("/decode/"):]
		log.Printf("Decoding code: %s\n", string(code))

		url, err := storage.Retrieve(code)
		if err != nil {
			respond(&MidgetResponse{
				Status:  204,
				Message: "invalid code " + code,
			}, w)
			return
		}
		log.Printf("Decoded %s to %s\n", code, url)
		respond(&MidgetResponse{
			Status:  200,
			Message: url,
		}, w)
	})
}

func Redirect(storage storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respond(&MidgetResponse{
				Status:  404,
				Message: "method not found\n",
			}, w)
			return
		}
		code := r.URL.Path[len("/"):]
		log.Printf("Redirect(): Found code %s\n", code)

		url, err := storage.Retrieve(code)
		if err != nil {
			respond(&MidgetResponse{
				Status:  404,
				Message: "Error: URL not found for code: " + code,
			}, w)
			return
		}
		log.Printf("Redirecting to url %s\n", url)

		http.Redirect(w, r, url, 301)
	})
}

func respond(response *MidgetResponse, w http.ResponseWriter) error {
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	w.Write(js)
	return nil
}
