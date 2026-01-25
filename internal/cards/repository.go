package cards

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MikebangSfilya/mindCards/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBQuerier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type CardRepository struct {
	pool *pgxpool.Pool
}

func NewCardPool(db *pgxpool.Pool) *CardRepository {
	repo := CardRepository{
		pool: db,
	}

	return &repo
}

func (ct *CardRepository) AddCard(ctx context.Context, db DBQuerier, userId int, card *MindCard) error {
	query := `
	INSERT INTO memory_cards 
    (user_id, title, card_description, tag, created_at, level_study, learned)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING card_id
    `

	card.Tag = strings.ToLower(card.Tag)

	err := db.QueryRow(ctx, query, userId, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned).Scan(&card.CardID)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}

func (ct *CardRepository) DeleteCard(ctx context.Context, db DBQuerier, cardId, userId int) error {
	query := `
    DELETE FROM memory_cards 
    WHERE card_id = $1 AND user_id = $2
    `

	result, err := db.Exec(ctx, query, cardId, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("card with title '%v' not found", cardId)
	}
	return nil
}

func (ct *CardRepository) UpdateCardDescription(ctx context.Context, db DBQuerier, cardId, userId int, newDesc string) (storage.CardRow, error) {
	query := `
        UPDATE memory_cards
        SET card_description = $1
        WHERE card_id = $2 AND user_id = $3
        RETURNING 
            card_id,
            user_id,
            title,
            card_description,
            tag,
            created_at,
            level_study,
            learned
    `

	var card storage.CardRow
	err := db.QueryRow(ctx, query, newDesc, cardId, userId).Scan(
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
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return storage.CardRow{}, fmt.Errorf("card not found or access denied")
		}
		return storage.CardRow{}, fmt.Errorf("failed to update card description: %w", err)
	}

	return card, nil

}

func (ct *CardRepository) GetCards(ctx context.Context, db DBQuerier, userId int, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE user_id = $1
	LIMIT $2 OFFSET $3
	
	`

	rows, err := db.Query(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (ct *CardRepository) GetCardsByTag(ctx context.Context, db DBQuerier, tag string, userId int, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE tag = $1 AND user_id = $2
	LIMIT $3 OFFSET $4
	`

	rows, err := db.Query(ctx, query, tag, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (ct *CardRepository) GetCardById(ctx context.Context, db DBQuerier, cardId, userId int) (storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE card_id = $1 AND user_id = $2
	`

	row := db.QueryRow(ctx, query, cardId, userId)

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
