package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/Gilgalad195/chirpy/internal/auth"
	"github.com/Gilgalad195/chirpy/internal/database"
)

func (c *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	inputChirp := chirpParams{}
	err := decoder.Decode(&inputChirp)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("An error occured: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "user could not be authenticated")
		return
	}
	id, err := auth.ValidateJWT(userToken, c.secret)
	if err != nil {
		log.Printf("An error occured: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "user is not authorized")
		return
	}

	inputChirp.UserID = id

	cleanedChirp := cleanChirp(inputChirp)

	if len(cleanedChirp.CleanedBody) > 140 {
		respBody := errorResp{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Something went wrong"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	ch := database.CreateChirpParams{
		Body:   cleanedChirp.CleanedBody,
		UserID: cleanedChirp.UserID,
	}

	chirp, err := c.queries.CreateChirp(r.Context(), ch)
	if err != nil {
		log.Printf("An error occured: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	newChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	dat, err := json.Marshal(newChirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)

}

func cleanChirp(chirp chirpParams) cleanedChirp {
	filtered := cleanedChirp{
		UserID: chirp.UserID,
	}
	inputSlice := strings.Split(chirp.Body, " ")
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range inputSlice {
		if slices.Contains(forbiddenWords, strings.ToLower(word)) {
			inputSlice[i] = "****"
		}
	}
	filtered.CleanedBody = strings.Join(inputSlice, " ")
	return filtered
}
