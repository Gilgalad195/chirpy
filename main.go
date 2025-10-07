package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	s := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("encountered an error: %v", err)
	}

}
