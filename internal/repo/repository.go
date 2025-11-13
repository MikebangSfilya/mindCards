package repo

import (
	"cards/internal/model"
	"context"
	"fmt"

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
	query := `
	INSERT INTO memory_cards 
	(title, description, tag, created_at, level_study, learned)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}

func (r *pgxRepository) DeleteCard(ctx context.Context, title string) error {
	query := `
	DELETE 
	FROM memory_cards 
	WHERE title = $1
	`

	result, err := r.db.Exec(ctx, query, title)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("card with title '%s' not found", title)
	}
	return nil
}
