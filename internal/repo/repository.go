package repo

import (
	"cards/internal/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *pgxRepository {
	repo := &pgxRepository{
		db: db,
	}

	return repo
}

func (r *pgxRepository) AddCard(ctx context.Context, card *model.MindCard) error {
	query := `
	INSERT INTO memory_cards 
	(title, description, tag, created_at, level_study, learned)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	err := r.db.QueryRow(ctx, query, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned).Scan(&card.ID)
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

func (r *pgxRepository) UptadeCardDescription(ctx context.Context, updt []string) {

}

func (r *pgxRepository) GetCards(ctx context.Context, limit, offset int16) (map[int]model.MindCard, error) {
	query := `
	SELECT *
	FROM memory_cards
	LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make(map[int]model.MindCard)
	for rows.Next() {
		var card model.MindCard
		err := rows.Scan(
			&card.ID,
			&card.Title,
			&card.Description,
			&card.Tag,
			&card.CreatedAt,
			&card.LevelStudy,
			&card.Learned,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		cards[int(card.ID)] = card
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return cards, err
}

func (r *pgxRepository) GetCardById(ctx context.Context, id int) model.MindCard {
	return model.MindCard{}
}
