package main

import (
	"log"
	"net/http"
)

// handler that resets apiConfig count to 0
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform == "dev" {
		err := cfg.db.Reset(r.Context())
		if err != nil {
			log.Printf("Error reseting users: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cfg.fileserverHits.Store(0)
		respondWithStatus(w, 200, "200 OK")
	} else {
		respondWithError(w, http.StatusForbidden, "403 Forbidden")
	}

}
