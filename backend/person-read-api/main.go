package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"strconv"

	_ "github.com/microsoft/go-mssqldb"
)

type Person struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	BirthDate string `json:"birthDate"`
	Age       int    `json:"age"`
	Address   string `json:"address"`
}

func main() {
	dsn := &url.URL{
		Scheme: "sqlserver",
		User: url.UserPassword(
			os.Getenv("GO_DB_USER"),
			os.Getenv("GO_DB_PASSWORD"),
		),
		Host: os.Getenv("GO_DB_HOST"),
	}

	params := url.Values{}
	params.Set("database", os.Getenv("GO_DB_NAME"))
	params.Set("encrypt", "true")
	// ใช้กับ SQL Server ในเครื่องพัฒนา
	params.Set("TrustServerCertificate", "true")
	dsn.RawQuery = params.Encode()

	db, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		log.Fatal("Database configuration failed: ", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = db.PingContext(ctx)
	cancel()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	log.Println("SQL Server connected")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{
			"status":  "ok",
			"service": "person-read-api",
		})
	})

	mux.HandleFunc("GET /api/persons", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx, `
			SELECT Id, FirstName, LastName, BirthDate, Address
			FROM dbo.Persons
			ORDER BY Id ASC
		`)
		if err != nil {
			log.Println("Query persons failed:", err)
			http.Error(w, "Cannot load persons", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		persons := make([]Person, 0)
		thailand := time.FixedZone("UTC+7", 7*60*60)
		currentYear := time.Now().In(thailand).Year()

		for rows.Next() {
			var person Person
			var birthDate time.Time

			err := rows.Scan(
				&person.ID,
				&person.FirstName,
				&person.LastName,
				&birthDate,
				&person.Address,
			)
			if err != nil {
				log.Println("Scan person failed:", err)
				http.Error(w, "Cannot read person", http.StatusInternalServerError)
				return
			}

			person.BirthDate = birthDate.Format("2006-01-02")
			person.Age = currentYear - birthDate.Year()
			persons = append(persons, person)
		}

		if err := rows.Err(); err != nil {
			log.Println("Read rows failed:", err)
			http.Error(w, "Cannot load persons", http.StatusInternalServerError)
			return
		}

		writeJSON(w, persons)
	})

	mux.HandleFunc("GET /api/persons/{id}", func(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var person Person
	var birthDate time.Time

	err = db.QueryRowContext(ctx, `
		SELECT Id, FirstName, LastName, BirthDate, Address
		FROM dbo.Persons
		WHERE Id = @Id
	`, sql.Named("Id", id)).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&birthDate,
		&person.Address,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println("Read person failed:", err)
		http.Error(w, "Cannot load person", http.StatusInternalServerError)
		return
	}

	thailand := time.FixedZone("UTC+7", 7*60*60)
	person.BirthDate = birthDate.Format("2006-01-02")
	person.Age = time.Now().In(thailand).Year() - birthDate.Year()

	writeJSON(w, person)
})

	server := &http.Server{
		Addr:              ":8080",
		Handler:           allowAngular(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Go API listening on :8080")
	log.Fatal(server.ListenAndServe())
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