package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := emailParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userEmail)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		return
	}

	user, err := c.queries.CreateUser(r.Context(), userEmail.Email)
	if err != nil {
		log.Printf("An eror occured: %s", err)
		return
	}

	NewUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	dat, err := json.Marshal(NewUser)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}
