package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
)

type personHandler struct {
	repository *personRepository
}

func (handler *personHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	persons, err := handler.repository.list(ctx)
	if err != nil {
		log.Println("Load persons failed:", err)
		http.Error(w, "Cannot load persons", http.StatusInternalServerError)
		return
	}
	writeJSON(w, persons)
}

func (handler *personHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	person, err := handler.repository.get(ctx, id)
	if err == sql.ErrNoRows {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println("Read person failed:", err)
		http.Error(w, "Cannot load person", http.StatusInternalServerError)
		return
	}

	writeJSON(w, person)
}