package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Gilgalad195/chirpy/internal/auth"
)

func (c *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting user token: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "no token was found")
		return
	}

	user, err := c.queries.GetUserFromRefreshToken(r.Context(), tok)
	if err != nil {
		log.Printf("Token not found: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "an error occured")
		return
	}

	accTok, err := auth.MakeJWT(user.ID, c.secret, 60*time.Minute)
	if err != nil {
		log.Printf("Error creating access token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "an error occured")
		return
	}

	t := tokenParams{
		Token: accTok,
	}

	dat, err := json.Marshal(t)
	if err != nil {
		log.Printf("Error marshaling json: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error marshaling json")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)

}

func (c *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting user token: %v", err)
		writeJSONError(w, http.StatusUnauthorized, "no token was found")
		return
	}

	err = c.queries.RevokeRefreshToken(r.Context(), tok)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "an error occured")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)
}
