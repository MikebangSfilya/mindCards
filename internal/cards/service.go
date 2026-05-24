package cards

import (
	"context"
	"errors"
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

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, q DBQuerier) error) error
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
	txManager TxManager
	Pool      *pgxpool.Pool
	Redis     *redis2.Redis
	logger    *slog.Logger
}

func NewService(repo Repo, txManager TxManager, pool *pgxpool.Pool, logger *slog.Logger, redis *redis2.Redis) *Service {
	serviceLogger := logger.With("component", "service")
	return &Service{
		Repo:      repo,
		txManager: txManager,
		Pool:      pool,
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
			results = append(results, &MDAddedDTO{
				Title:       card.Title,
				Description: card.Description,
				Tag:         card.Tag,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("batch add cards failed: %w", err)
	}

	s.invalidateUserCollectionCaches(ctx, userId)
	s.logger.Info("batch add completed", "count", len(results))

	return results, nil
}

// DeleteCard card from DB
func (s *Service) DeleteCard(ctx context.Context, cardId, userId int) error {
	key := fmt.Sprintf("cards:u:%d:c:%d", userId, cardId)

	err := s.txManager.WithinTransaction(ctx, func(ctx context.Context, q DBQuerier) error {
		if err := s.Repo.DeleteCard(ctx, q, cardId, userId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("delete card failed: %w", err)
	}

	s.invalidateUserCaches(ctx, userId, key)
	s.logger.Info("batch delete completed", "cardId", cardId)
	return nil
}

// UpdateCardDescription new description in DB
func (s *Service) UpdateCardDescription(ctx context.Context, cardId, UserID int, cardsUp Update) (*MindCard, error) {
	key := fmt.Sprintf("cards:u:%d:c:%d", UserID, cardId)

	row, err := s.Repo.UpdateCardDescription(ctx, s.Pool, cardId, UserID, cardsUp.NewDescription)
	if err != nil {
		return nil, fmt.Errorf("update card description: %w", err)
	}

	s.logger.Info("update completed", "cardId", cardId)
	s.invalidateUserCaches(ctx, UserID, key)
	return rowToCard(row), nil
}

// Getcards list of cards
func (s *Service) GetCards(ctx context.Context, userId int, limit, offset int16) ([]*MindCard, error) {
	key := fmt.Sprintf("cards:u:%d:l:%d:o:%d", userId, limit, offset)

	var cachedCards []*MindCard
	if s.getFromCache(ctx, key, &cachedCards) {
		return cachedCards, nil
	}

	rows, err := s.Repo.GetCards(ctx, s.Pool, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	cards := rowsToCards(rows)

	go s.setCache(key, cards, 5*time.Minute)

	return cards, nil
}

// Get cards filtered by Tag
func (s *Service) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]*MindCard, error) {
	key := fmt.Sprintf("cards:u:%d:t:%s:l:%d:o:%d", userId, tag, limit, offset)

	var cachedCards []*MindCard
	if s.getFromCache(ctx, key, &cachedCards) {
		return cachedCards, nil
	}

	rows, err := s.Repo.GetCardsByTag(ctx, s.Pool, tag, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	cards := rowsToCards(rows)

	go s.setCache(key, cards, 5*time.Minute)

	return cards, nil
}

// Get one card by ID
func (s *Service) GetCardById(ctx context.Context, cardID, userID int) (*MindCard, error) {
	key := fmt.Sprintf("cards:u:%d:c:%d", userID, cardID)

	var card MindCard
	if s.getFromCache(ctx, key, &card) {
		return &card, nil
	}

	row, err := s.Repo.GetCardById(ctx, s.Pool, cardID, userID)
	if err != nil {
		return nil, err
	}

	cardDB := rowToCard(row)

	go s.setCache(key, cardDB, 24*time.Hour)
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

func (s *Service) getFromCache(ctx context.Context, key string, dest any) bool {
	err := s.Redis.Get(ctx, key, dest)
	if err == nil {
		return true
	}
	if errors.Is(err, redis2.ErrCacheMiss) {
		return false
	}
	s.logger.Warn("cache failure, clearing key", "key", key, "error", err)

	_ = s.Redis.Delete(ctx, key)
	return false
}

func (s *Service) setCache(key string, dest any, ttl time.Duration) {
	setCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.Redis.Set(setCtx, key, dest, ttl); err != nil {
		s.logger.Warn("redis set failed", "err", err)
	}
}

func (s *Service) invalidateUserCaches(ctx context.Context, userID int, keys ...string) {
	for _, key := range keys {
		if err := s.Redis.Delete(ctx, key); err != nil {
			s.logger.Warn("redis delete failed", "key", key, "err", err)
		}
	}

	s.invalidateUserCollectionCaches(ctx, userID)
}

func (s *Service) invalidateUserCollectionCaches(ctx context.Context, userID int) {
	prefixes := []string{
		fmt.Sprintf("cards:u:%d:l:", userID),
		fmt.Sprintf("cards:u:%d:t:", userID),
	}

	for _, prefix := range prefixes {
		if err := s.Redis.DeleteByPrefix(ctx, prefix); err != nil {
			s.logger.Warn("redis prefix delete failed", "prefix", prefix, "err", err)
		}
	}
}
