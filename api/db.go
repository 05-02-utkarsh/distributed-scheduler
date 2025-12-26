package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() {
	connStr := "postgres://scheduler:scheduler@localhost:5432/scheduler_db"

	var err error
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL")
}

// This connects Go → Docker PostgreSQL.
