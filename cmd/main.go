package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Impervguin/ds-lab1/internal/migrations"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://program:test@localhost:5432/persons?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	if err := migrations.Run(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	log.Println("migrations applied successfully")
}
