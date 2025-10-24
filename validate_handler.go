package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	cleanedChirp := cleanChirp(chirp)

	if len(cleanedChirp.CleanedBody) > 140 {
		respBody := errorResp{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Something went wrong"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	} else {
		respBody := cleanedChirp
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Something went wrong"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}
}

func cleanChirp(chirp parameters) cleaned {
	filtered := cleaned{}
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
