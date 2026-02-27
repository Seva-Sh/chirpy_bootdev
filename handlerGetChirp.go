package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("chirpID")

	requestUUID, err := uuid.Parse(requestID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), requestUUID)
	// check for not found chirp
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "No chirp found")
		return
	}
	// check for any other error
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
