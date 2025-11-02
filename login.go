package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Gilgalad195/chirpy/internal/auth"
)

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	userCreds := userParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userCreds)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		writeJSONError(w, http.StatusBadRequest, "error decoding parameters")
		return
	}

	user, err := c.queries.GetUser(r.Context(), userCreds.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		writeJSONError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(userCreds.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		writeJSONError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if !match {
		writeJSONError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	jsonUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	dat, err := json.Marshal(jsonUser)
	if err != nil {
		log.Printf("Error marshaling json: %s", err)
		writeJSONError(w, http.StatusInternalServerError, "error marshaling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
