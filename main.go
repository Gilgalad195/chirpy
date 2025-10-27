package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Gilgalad195/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	apiCfg.platform = os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Unable to connect to database: %s", err)
	}
	dbQueries := database.New(db)
	apiCfg.queries = dbQueries

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.countHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	h := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(h)))

	s := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("encountered an error: %v", err)
	}
}
