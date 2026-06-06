package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Zähler sicher um 1 erhöhen
		cfg.fileserverHits.Add(1)

		// 2. Den ursprünglichen Handler aufrufen, um die Anfrage zu verarbeiten
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// Setze den Content-Type Header auf Plain Text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Setze den HTTP-Statuscode auf 200 OK
	w.WriteHeader(http.StatusOK)

	// Lies den aktuellen Wert sicher aus und schreibe ihn in die Response
	hits := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %d", hits)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	// Setze den HTTP-Statuscode auf 200 OK
	w.WriteHeader(http.StatusOK)

}

func main() {
	cfg := &apiConfig{}
	fileServer := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.Handle("/", fileServer)

	mux.Handle("/app", cfg.middlewareMetricsInc(fileServer))
	mux.HandleFunc("/reset", cfg.handlerReset)
	mux.HandleFunc("/metrics", cfg.handlerMetrics)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // Leaving this mux empty triggers a 404 for all routes
	}

	log.Println("Starting server on http://localhost:8080...")

	log.Fatal(server.ListenAndServe())

}
