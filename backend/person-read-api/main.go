package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func main() {
	db, err := openDatabase()
	if err != nil {
		log.Fatal("Database configuration failed: ", err)
	}
	defer db.Close()

	err = pingDatabase(db)
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	log.Println("SQL Server connected")

	server := newServer(db)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           server,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Go API listening on :8080")
	log.Fatal(server.ListenAndServe())
}
