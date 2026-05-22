package main

import (
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.Handle("/", fileServer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // Leaving this mux empty triggers a 404 for all routes
	}

	log.Println("Starting server on http://localhost:8080...")

	log.Fatal(server.ListenAndServe())

}
