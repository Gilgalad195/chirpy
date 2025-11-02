package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type chirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type cleanedChirp struct {
	CleanedBody string    `json:"cleaned_body"`
	UserID      uuid.UUID `json:"user_id"`
}
type errorResp struct {
	Error string `json:"error"`
}

type userParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(errorResp{Error: msg})
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(dat)
}
