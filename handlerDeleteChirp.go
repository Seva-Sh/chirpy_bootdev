package main

import (
	"net/http"

	"github.com/Seva-Sh/chirpy_bootdev/internal/auth"
	"github.com/Seva-Sh/chirpy_bootdev/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// get the token
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve token")
		return
	}

	// get user
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "The given token is invalid")
		return
	}

	// get chirp via getting ID
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID")
		return
	}

	if chirp.UserID == userID {
		// here we can delete the chirp
		err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
			ID:     chirp.ID,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Failed to delete chirp")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	} else {
		respondWithError(w, http.StatusForbidden, "The user is not the author of the Chirp")
		return
	}
}
