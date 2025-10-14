package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (c *apiConfig) countHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("Hits: %v", c.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (c *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	c.fileserverHits.Store(0)
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", apiCfg.countHandler)
	mux.HandleFunc("POST /reset", apiCfg.resetHandler)

	h := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h)))

	s := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("encountered an error: %v", err)
	}

}
