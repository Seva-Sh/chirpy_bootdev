package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	port := "8080"
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/healthz", handlerReadiness)

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// Ensure the program loggs an error if encountered with server listening failure
	log.Fatal(srv.ListenAndServe())
}
