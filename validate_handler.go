package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorResp struct {
		Error string `json:"error"`
	}
	type validResp struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	if len(params.Body) > 140 {
		respBody := errorResp{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Something went wrong"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	} else {
		respBody := validResp{
			Valid: true,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(`{"error":"Something went wrong"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}

}
