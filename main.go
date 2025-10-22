package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.countHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)

	h := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h)))

	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("encountered an error: %v", err)
	}
}
