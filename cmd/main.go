package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	healthHandler "github.com/Impervguin/ds-lab1/internal/handlers/health"
	personHandler "github.com/Impervguin/ds-lab1/internal/handlers/person"
	"github.com/Impervguin/ds-lab1/internal/migrations"
	personRepo "github.com/Impervguin/ds-lab1/internal/repository/pgx/person"
	healthService "github.com/Impervguin/ds-lab1/internal/service/health"
)

func main() {
	cfg, err := NewAppConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
		return
	}

	if err := runMigrations(cfg.DatabaseDSN); err != nil {
		log.Fatalf("run migrations: %v", err)
	}
	log.Println("migrations applied successfully")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("create pgx pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	repo := personRepo.NewPgxPersonRepository(pool)
	handler := personHandler.NewHandler(repo)
	hs := healthService.NewService(pool)
	health := healthHandler.NewHandler(hs)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	handler.Register(r)
	health.Register(r)

	hs.SetReady()
	log.Printf("starting server on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
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
