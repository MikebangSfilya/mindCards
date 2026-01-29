package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
		slog.String("host", cfg.DBHost),
		slog.String("port", cfg.DBPort),
		slog.String("database", cfg.DBName),
		slog.String("user", cfg.DBUser))

	dbPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {

		slog.Error("Database connection fail",
			slog.Group("Error",
				slog.String("error", err.Error()),
				slog.Time("time", time.Now())))

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
