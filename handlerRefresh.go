package main

import (
	"net/http"
	"time"

	"github.com/Seva-Sh/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "The given token is invalid")
		return
	}

	expiresIn := time.Duration(3600) * time.Second
	// create a JWT token
	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwt,
	})
}
