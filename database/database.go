package database

import (
	"cards/internal/configurate"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDataBase(cfg configurate.Config) *pgxpool.Pool {

	consStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUSer, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	dbPool, err := pgxpool.New(context.Background(), consStr)
	if err != nil {
		return nil
	}
	return dbPool
}
