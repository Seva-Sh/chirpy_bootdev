package main

import (
	"net/http"
	"sort"

	"github.com/Seva-Sh/chirpy_bootdev/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	sortStr := r.URL.Query().Get("sort")

	// sort via author_id logic
	var chirps []database.Chirp
	var err error
	var author_uuid uuid.UUID
	if author_id != "" {
		author_uuid, err = uuid.Parse(author_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error parsing author ID")
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), author_uuid)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
			return
		}
	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
			return
		}
	}

	// create a slice of json style chirps
	chirpResponses := []chirpResponse{}
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if sortStr == "desc" {
		sort.Slice(chirpResponses, func(i, j int) bool { return chirpResponses[j].CreatedAt.Before(chirpResponses[i].CreatedAt) })

	} else {
		sort.Slice(chirpResponses, func(i, j int) bool { return chirpResponses[i].CreatedAt.Before(chirpResponses[j].CreatedAt) })
	}

	respondWithJSON(w, http.StatusOK, chirpResponses)
}
