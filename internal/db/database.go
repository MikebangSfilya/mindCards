package database

import (
	"cards/internal/config"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDataBase(cfg config.Config) *pgxpool.Pool {

	// connStr := "postgres://test:test@localhost:5434/testdb?sslmode=disable"
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUSer,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	slog.Info(
		"Database connection",
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"database", cfg.DBName,
		"user", cfg.DBUSer,
	)

	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {

		slog.Error(
			"Database connection fail",
			"error", err,
			"db", "pgxpool")
		return nil
	}

	if err := dbPool.Ping(context.Background()); err != nil {

		slog.Error(
			"Database ping fail",
			"error", err,
			"db", "pgxpool")
		return nil
	}

	slog.Info("Database connected successfully!")
	return dbPool
}
