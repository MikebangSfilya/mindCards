package cards

import (
	"context"
	"fmt"
	"strings"

	"github.com/MikebangSfilya/mindCards/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRepository struct {
	db *pgxpool.Pool
}

func NewPool(db *pgxpool.Pool) *pgxRepository {
	repo := &pgxRepository{
		db: db,
	}

	return repo
}

func (r *pgxRepository) AddCard(ctx context.Context, card *MindCard) error {
	query := `
	  INSERT INTO memory_cards 
    (user_id, title, card_description, tag, created_at, level_study, learned)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING card_id
    `

	card.Tag = strings.ToLower(card.Tag)

	err := r.db.QueryRow(ctx, query, card.UserID, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned).Scan(&card.ID)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}

func (r *pgxRepository) DeleteCard(ctx context.Context, card_id, user_id string) error {
	query := `
    DELETE FROM memory_cards 
    WHERE card_id = $1 AND user_id = $2
    `

	result, err := r.db.Exec(ctx, query, card_id, user_id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("card with title '%s' not found", card_id)
	}
	return nil
}

func (r *pgxRepository) UpdateCardDescription(ctx context.Context, card_id, user_id, newDesc string) error {
	query := `
	UPDATE memory_cards
	SET card_description = $1
	WHERE card_id = $2 AND user_id = $3
	`

	_, err := r.db.Exec(ctx, query, newDesc, card_id, user_id)
	if err != nil {
		return err
	}
	return nil

}

func (r *pgxRepository) GetCards(ctx context.Context, limit, offset int16, user_id string) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	LIMIT $1 OFFSET $2
	WHERE user_id = $3
	`

	rows, err := r.db.Query(ctx, query, limit, offset, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *pgxRepository) GetCardsByTag(ctx context.Context, tag string, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE tag = $1
	LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, tag, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *pgxRepository) GetCardById(ctx context.Context, id string) (storage.CardRow, error) {
	query := `
	SELECT *
	FROM memory_cards
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	return scanRow(row)
}

func scanRow(row pgx.Row) (storage.CardRow, error) {
	card := storage.CardRow{}
	err := row.Scan(
		&card.ID,
		&card.Title,
		&card.Description,
		&card.Tag,
		&card.CreatedAt,
		&card.LevelStudy,
		&card.Learned,
	)
	if err != nil {
		return storage.CardRow{}, err
	}

	return card, nil
}

func scanRows(rows pgx.Rows) ([]storage.CardRow, error) {
	var cardsRow []storage.CardRow
	for rows.Next() {
		var Row storage.CardRow
		err := rows.Scan(
			&Row.ID,
			&Row.Title,
			&Row.Description,
			&Row.Tag,
			&Row.CreatedAt,
			&Row.LevelStudy,
			&Row.Learned,
		)
		if err != nil {
			return nil, err
		}
		cardsRow = append(cardsRow, Row)
	}

	return cardsRow, rows.Err()
}
