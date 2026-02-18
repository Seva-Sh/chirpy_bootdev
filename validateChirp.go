package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// helper func that responds with error JSON
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respBody := returnVals{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

// helper func that responds with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

// func containsProfane(s string) bool {
// 	lowerStr := strings.ToLower(s)
// 	splitStr := strings.Split(lowerStr, " ")
// 	for _, v := range splitStr {
// 		if slices.Contains(profaneWords, v) {
// 			return true
// 		}
// 	}

// 	return false
// }

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

// the handler that decodes the JSON body,
// uses a helper function to respond with the appropriate JSON and status code
func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
	} else {
		cleanedStr := cleanString(params.Body)
		type cleanedResponse struct {
			CleanedBody string `json:"cleaned_body"`
		}
		respondWithJSON(w, http.StatusOK, cleanedResponse{CleanedBody: cleanedStr})
	}
}
