package main

import (
	"log"
	"net/http"
	"os"

	"github.com/natezyz/midget/handlers"
	"github.com/natezyz/midget/storage"
)

func main() {
	s := &storage.Map{}
	s.Init()

	http.Handle("/", handlers.Redirect(s))
	http.Handle("/process", handlers.Process(s))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting midget on port: %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
