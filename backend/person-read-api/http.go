package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func newServer(db *sql.DB) http.Handler {
	people := &personHandler{repository: &personRepository{db: db}}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{
			"status":  "ok",
			"service": "person-read-api",
		})
	})
	mux.HandleFunc("GET /api/persons", people.list)
	mux.HandleFunc("GET /api/persons/{id}", people.get)

	return allowAngular(mux)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Write JSON failed:", err)
	}
}

func allowAngular(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		if r.Header.Get("Origin") == "http://localhost:4200" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}