package main

import (
	"log"
	"net/http"
)

func (c *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if c.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	c.fileserverHits.Store(0)

	if err := c.queries.ResetUsers(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to reset users: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
