package main

import (
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

type emailParams struct {
	Email string `json:"email"`
}

// type validResp struct {
// 	Valid bool `json:"valid"`
// }
