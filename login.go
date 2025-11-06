package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Gilgalad195/chirpy/internal/auth"
	"github.com/Gilgalad195/chirpy/internal/database"
)

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	loginCreds := loginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginCreds)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		writeJSONError(w, http.StatusBadRequest, "error decoding parameters")
		return
	}
	if loginCreds.ExpiresInSeconds <= 0 || loginCreds.ExpiresInSeconds > 3600 {
		loginCreds.ExpiresInSeconds = 3600
	}

	user, err := c.queries.GetUser(r.Context(), loginCreds.Email)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(loginCreds.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if !match {
		writeJSONError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, c.secret, time.Duration(loginCreds.ExpiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("Error creating user token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "An error occured")
		return
	}

	token, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "An error occurred")
		return
	}

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	}

	refreshToken, err := c.queries.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "An error occured")
		return
	}

	jsonResponse := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken.Token,
	}

	dat, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error marshaling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
