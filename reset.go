package main

import (
	"net/http"
)

func (c *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.fileserverHits.Store(0)
}
