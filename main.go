package main

import (
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()

	s := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	s.ListenAndServe()
}
