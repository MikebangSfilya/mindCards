package cards

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	redis2 "github.com/MikebangSfilya/mindCards/internal/repository/redis"
	"github.com/MikebangSfilya/mindCards/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManager struct {
	pool *pgxpool.Pool
}

type Repo interface {
	AddCard(ctx context.Context, db DBQuerier, userId int, card *MindCard) error
	UpdateCardDescription(ctx context.Context, db DBQuerier, cardId, userId int, newDesc string) (storage.CardRow, error)
	DeleteCard(ctx context.Context, db DBQuerier, cardId, userId int) error
	GetCards(ctx context.Context, db DBQuerier, userId int, limit, offset int16) ([]storage.CardRow, error)
	GetCardsByTag(ctx context.Context, db DBQuerier, tag string, userId int, limit, offset int16) ([]storage.CardRow, error)
	GetCardById(ctx context.Context, db DBQuerier, cardId, userId int) (storage.CardRow, error)
}

func NewTxManager(pool *pgxpool.Pool) *txManager {
	return &txManager{pool: pool}
}

func (tm *txManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context, q DBQuerier) error) error {
	tx, err := tm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("txManager.WithinTransaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

type Service struct {
	Repo      Repo
	txManager txManager
	Redis     *redis2.Redis
	logger    *slog.Logger
}

func NewService(repo *CardRepository, txManager *txManager, logger *slog.Logger, redis *redis2.Redis) *Service {
	serviceLogger := logger.With("component", "service")
	return &Service{
		Repo:      repo,
		txManager: *txManager,
		Redis:     redis,
		logger:    serviceLogger,
	}
}

// AddCards cards to DB
func (s *Service) AddCards(ctx context.Context, userId int, cardParams []Card) ([]*MDAddedDTO, error) {

	results := make([]*MDAddedDTO, 0, len(cardParams))

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context, q DBQuerier) error {
		for _, cardParam := range cardParams {
			card := NewCard(cardParam.Title, cardParam.Description, cardParam.Tag)

			if err := s.Repo.AddCard(ctx, q, userId, card); err != nil {
				return err
			}
			results = append(results, &MDAddedDTO{})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("batch add cards failed: %w", err)
	}

	s.logger.Info("batch add completed", "count", len(results))

	return results, nil
}

// DeleteCard card from DB
func (s *Service) DeleteCard(ctx context.Context, cardId, userId int) error {
	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context, q DBQuerier) error {
		if err := s.Repo.DeleteCard(ctx, q, cardId, userId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("batch delete cards failed: %w", err)
	}

	s.logger.Info("batch delete completed", "cardId", cardId)
	return nil
}

// UpdateCardDescription new description in DB
func (s *Service) UpdateCardDescription(ctx context.Context, cardId, UserID int, cardsUp Update) (*MindCard, error) {
	row, err := s.Repo.UpdateCardDescription(ctx, nil, cardId, UserID, cardsUp.NewDescription)
	if err != nil {
		return nil, fmt.Errorf("update card description: %w", err)
	}

	s.logger.Info("update completed", "cardId", cardId)
	return rowToCard(row), nil
}

// UpdateLvl Возможно не понадобится
func (s *Service) UpdateLvl() {

}

// Getcards list of cards
func (s *Service) GetCards(ctx context.Context, userId int, limit, offset int16) ([]*MindCard, error) {
	key := fmt.Sprintf("cards:u:%d:l:%d:o:%d", userId, limit, offset)

	var cachedCards []*MindCard
	err := s.Redis.Get(ctx, key, &cachedCards)
	if err == nil {
		return cachedCards, nil
	}

	rows, err := s.Repo.GetCards(ctx, nil, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	cards := rowsToCards(rows)
	go func() {
		setCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.Redis.Set(setCtx, key, cards, 1*time.Minute); err != nil {
			s.logger.Warn("redis set failed", "err", err)
		}
	}()
	return cards, nil
}

// Get cards filtered by Tag
func (s *Service) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]*MindCard, error) {

	rows, err := s.Repo.GetCardsByTag(ctx, nil, tag, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	return rowsToCards(rows), nil
}

// Get one card by ID
func (s *Service) GetCardById(ctx context.Context, cardID, userID int) (*MindCard, error) {
	key := fmt.Sprintf("cards:u:%d:c:%d", userID, cardID)

	var card MindCard
	err := s.Redis.Get(ctx, key, &card)
	if err == nil {
		return &card, nil
	}

	row, err := s.Repo.GetCardById(ctx, nil, cardID, userID)
	if err != nil {
		return nil, err
	}

	cardDB := rowToCard(row)

	go func() {
		setCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := s.Redis.Set(setCtx, key, cardDB, time.Hour*24); err != nil {
			s.logger.Error("failed to set card", "error", err)
		}
	}()
	return cardDB, nil
}

func rowToCard(row storage.CardRow) *MindCard {
	return &MindCard{
		CardID:      row.CardID,
		UserID:      row.UserID,
		Title:       row.Title,
		Description: row.Description,
		Tag:         row.Tag,
		CreatedAt:   row.CreatedAt,
		LevelStudy:  row.LevelStudy,
		Learned:     row.Learned,
	}
}

func rowsToCards(rows []storage.CardRow) []*MindCard {

	if rows == nil {
		return nil
	}

	result := make([]*MindCard, len(rows))

	for i, row := range rows {
		result[i] = rowToCard(row)
	}

	return result

}
