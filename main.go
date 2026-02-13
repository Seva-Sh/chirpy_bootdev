package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// stateful config
type apiConfig struct {
	fileserverHits atomic.Int32
}

// middleware pattern that wraps handlers and increments the counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// handler that writes the number of request that have been counted
func (cfg *apiConfig) handlerCount(w http.ResponseWriter, r *http.Request) {
	val32 := cfg.fileserverHits.Load()
	str := fmt.Sprintf("Hits: %d", val32)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(str))
}

// handler that resets apiConfig count to 0
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}

func main() {
	port := "8080"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerCount)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	// ensure the program loggs an error if encountered with server listening failure
	log.Fatal(srv.ListenAndServe())
}
