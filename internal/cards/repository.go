package cards

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/MikebangSfilya/mindCards/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository struct {
	db *pgxpool.Pool
}

type CardTransaction struct {
	tx pgx.Tx
}

func NewCardPool(db *pgxpool.Pool) Repo {
	repo := CardRepository{
		db: db,
	}

	return &repo
}

func (r *CardRepository) BeginTransaction(ctx context.Context) (Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	return &CardTransaction{
		tx: tx,
	}, nil
}

func (ct *CardTransaction) AddCard(ctx context.Context, userId int, card *MindCard) error {
	query := `
	INSERT INTO memory_cards 
    (user_id, title, card_description, tag, created_at, level_study, learned)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING card_id
    `

	card.Tag = strings.ToLower(card.Tag)

	err := ct.tx.QueryRow(ctx, query, userId, card.Title, card.Description, card.Tag, card.CreatedAt, card.LevelStudy, card.Learned).Scan(&card.CardID)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	return nil
}

func (ct *CardTransaction) DeleteCard(ctx context.Context, cardId, userId int) error {
	query := `
    DELETE FROM memory_cards 
    WHERE card_id = $1 AND user_id = $2
    `

	result, err := ct.tx.Exec(ctx, query, cardId, userId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("card with title '%v' not found", cardId)
	}
	return nil
}

func (ct *CardTransaction) UpdateCardDescription(ctx context.Context, cardId, userId int, newDesc string) (storage.CardRow, error) {
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
	err := ct.tx.QueryRow(ctx, query, newDesc, cardId, userId).Scan(
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

func (ct *CardTransaction) GetCards(ctx context.Context, userId int, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE user_id = $1
	LIMIT $2 OFFSET $3
	
	`

	rows, err := ct.tx.Query(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (ct *CardTransaction) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE tag = $1 AND user_id = $2
	LIMIT $3 OFFSET $4
	`

	rows, err := ct.tx.Query(ctx, query, tag, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (ct *CardTransaction) GetCardById(ctx context.Context, cardId, userId int) (storage.CardRow, error) {
	query := `
	SELECT card_id, user_id, title, card_description, tag, created_at, level_study, learned
	FROM memory_cards
	WHERE card_id = $1 AND user_id = $2
	`

	row := ct.tx.QueryRow(ctx, query, cardId, userId)

	return scanRow(row)
}

func (ct *CardTransaction) Commit(ctx context.Context) error {
	return ct.tx.Commit(ctx)
}

func (ct *CardTransaction) Rollback(ctx context.Context) error {
	return ct.tx.Rollback(ctx)
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
