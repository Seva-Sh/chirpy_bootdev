package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Seva-Sh/chirpy_bootdev/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	passwordCheck, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !passwordCheck {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresIn := time.Hour
	if params.ExpiresInSeconds != nil {
		if *params.ExpiresInSeconds > 3600 {
			expiresIn = time.Duration(3600) * time.Second
		} else {
			expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
		}
	}
	jwt, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     jwt,
	})
}
