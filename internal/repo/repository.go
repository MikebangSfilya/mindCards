package repo

import (
	"cards/internal/model"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo interface {
	AddCard(ctx context.Context, card model.MindCard) error
	UptadeCardDescription(ctx context.Context)
	DeleteCard(ctx context.Context)
}

type pgxRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *pgxRepository {
	repo := &pgxRepository{
		db: db,
	}

	return repo
}

func (r *pgxRepository) AddCard(ctx context.Context, card model.MindCard) error {
	if r == nil {
		slog.Error("repository is nil")
		return fmt.Errorf("repository is nil")
	}
	if r.db == nil {
		slog.Error("database connection is nil")
		return fmt.Errorf("database connection is nil")
	}

	query := `
	INSERT INTO memory_cards 
	(title, description, tag, created_at, level_study, learned)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	return err
}
