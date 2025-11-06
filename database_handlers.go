package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Gilgalad195/chirpy/internal/auth"
	"github.com/Gilgalad195/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	loginCreds := loginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginCreds)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		writeJSONError(w, http.StatusBadRequest, "error decoding parameters")
		return
	}

	if loginCreds.Email == "" || loginCreds.Password == "" {
		writeJSONError(w, http.StatusBadRequest, "Email and Password are required.")
		return
	}

	hashedPass, err := auth.HashPassword(loginCreds.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "A server error occured")
		return
	}

	createUserParams := database.CreateUserParams{
		Email:          loginCreds.Email,
		HashedPassword: hashedPass,
	}

	user, err := c.queries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		log.Printf("An error occured: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	dat, err := json.Marshal(newUser)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error marshaling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

func (c *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := c.queries.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("An error occured: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error retrieving chirps")
		return
	}

	chirps := make([]Chirp, 0, len(data))
	for _, ch := range data {
		newChirp := Chirp{
			ID:        ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body:      ch.Body,
			UserID:    ch.UserID,
		}
		chirps = append(chirps, newChirp)
	}

	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error marshaling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (c *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	chirpId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("An error occured: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid UUID")
		return
	}
	chirp, err := c.queries.GetChirp(r.Context(), chirpId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "chirp not found")
			return
		}
		log.Printf("An error occured: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error retrieving chirp")
		return
	}

	jsonChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	dat, err := json.Marshal(jsonChirp)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error marshaling json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}
