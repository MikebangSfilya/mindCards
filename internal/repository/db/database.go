package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MikebangSfilya/mindCards/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDataBase(cfg config.Config) *pgxpool.Pool {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
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
		"user", cfg.DBUser,
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
