package main

import (
	"log"
	"net/http"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func main() {
	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", fileServer))

	mux.HandleFunc("/healthz", HealthzHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // Leaving this mux empty triggers a 404 for all routes
	}

	log.Println("Starting server on http://localhost:8080...")

	log.Fatal(server.ListenAndServe())

}
