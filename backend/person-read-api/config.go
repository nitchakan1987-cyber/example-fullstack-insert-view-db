package main

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"time"
)

func openDatabase() (*sql.DB, error) {
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
	params.Set("TrustServerCertificate", "true")
	dsn.RawQuery = params.Encode()

	db, err := sql.Open("sqlserver", dsn.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}

func pingDatabase(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}