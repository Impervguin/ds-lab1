package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	personHandler "github.com/Impervguin/ds-lab1/internal/handlers/person"
	"github.com/Impervguin/ds-lab1/internal/migrations"
	personRepo "github.com/Impervguin/ds-lab1/internal/repository/pgx/person"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://program:test@localhost:5432/persons?sslmode=disable"
	}

	if err := runMigrations(dsn); err != nil {
		log.Fatalf("run migrations: %v", err)
	}
	log.Println("migrations applied successfully")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("create pgx pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := personRepo.NewPgxPersonRepository(pool)
	handler := personHandler.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handler.Register(r)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	return migrations.Run(db)
}
