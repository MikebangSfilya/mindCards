package cards

import (
	"context"
	"fmt"
	"strings"

	"github.com/MikebangSfilya/mindCards/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type cardRepository struct {
	db *pgxpool.Pool
}

func NewCardPool(db *pgxpool.Pool) *cardRepository {
	repo := &cardRepository{
		db: db,
	}

	return repo
}

func (r *cardRepository) AddCard(ctx context.Context, userId int, card *MindCard) error {
	query := `
	INSERT INTO memory_cards 
    (user_id, title, card_description, tag, created_at, level_study, learned)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING card_id
    `

	card.Tag = strings.ToLower(card.Tag)

	err := r.db.QueryRow(ctx, query, userId, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned).Scan(&card.CardID)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}

func (r *cardRepository) DeleteCard(ctx context.Context, cardId, userId int) error {
	query := `
    DELETE FROM memory_cards 
    WHERE card_id = $1 AND user_id = $2
    `

	result, err := r.db.Exec(ctx, query, cardId, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("card with title '%v' not found", cardId)
	}
	return nil
}

func (r *cardRepository) UpdateCardDescription(ctx context.Context, cardId, userId int, newDesc string) error {
	query := `
	UPDATE memory_cards
	SET card_description = $1
	WHERE card_id = $2 AND user_id = $3
	`

	_, err := r.db.Exec(ctx, query, newDesc, cardId, userId)
	if err != nil {
		return err
	}
	return nil

}

func (r *cardRepository) GetCards(ctx context.Context, userId int, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE user_id = $1
	LIMIT $2 OFFSET $3
	
	`

	rows, err := r.db.Query(ctx, query, userId, limit, offset)
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

func (r *cardRepository) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE tag = $1 AND user_id = $2
	LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, tag, userId, limit, offset)
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

func (r *cardRepository) GetCardById(ctx context.Context, cardId, userId int) (storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE card_id = $1 AND user_id = $2
	`

	row := r.db.QueryRow(ctx, query, cardId, userId)

	return scanRow(row)
}

func scanRow(row pgx.Row) (storage.CardRow, error) {
	card := storage.CardRow{}
	err := row.Scan(
		&card.CardID,
		&card.UserID,
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
			&Row.CardID,
			&Row.UserID,
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
