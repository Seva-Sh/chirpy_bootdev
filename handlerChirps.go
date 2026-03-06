package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Seva-Sh/chirpy_bootdev/internal/auth"
	"github.com/Seva-Sh/chirpy_bootdev/internal/database"
)

// helper func that responds with cleaned JSON
func cleanString(s string) string {
	// map of profane words
	profaneMap := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(s, " ")
	for i, word := range words {
		if _, ok := profaneMap[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	// obtain token string
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve token")
		return
	}

	// obtain userID
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating user")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters")
		return
	}

	body := strings.TrimSpace(params.Body)

	if len(body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty")
		return
	} else if len(body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	} else {
		cleanedStr := cleanString(body)
		newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleanedStr,
			UserID: userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error creating a chirp")
			return
		}
		respondWithJSON(w, http.StatusCreated, chirpResponse{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		})
	}
}
