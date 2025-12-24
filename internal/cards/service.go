package cards

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/MikebangSfilya/mindCards/internal/storage"
)

type Transaction interface {
	AddCard(ctx context.Context, userId int, card *MindCard) error
	UpdateCardDescription(ctx context.Context, cardId, userId int, newDesc string) (storage.CardRow, error)
	DeleteCard(ctx context.Context, cardId, userId int) error
	GetCards(ctx context.Context, userId int, limit, offset int16) ([]storage.CardRow, error)
	GetCardById(ctx context.Context, cardId, userId int) (storage.CardRow, error)
	GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]storage.CardRow, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Repo interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// general Service struct
type Service struct {
	Repo   Repo
	logger *slog.Logger
}

func NewService(repo Repo, logger *slog.Logger) *Service {
	serviceLogger := logger.With("component", "service")
	return &Service{
		Repo:   repo,
		logger: serviceLogger,
	}
}

// AddCards cards to DB
func (s *Service) AddCards(ctx context.Context, userId int, cardParams []Card) ([]*MDAddedDTO, error) {

	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	
	defer tx.Rollback(ctx)

	results := make([]*MDAddedDTO, 0, len(cardParams))

	for i, cardParam := range cardParams {
		card := NewCard(cardParam.Title, cardParam.Description, cardParam.Tag)

		if err := tx.AddCard(ctx, userId, card); err != nil {
			s.logger.Error("failed to add card", "index", i, "error", err)
			_ = tx.Rollback(ctx)
			return nil, fmt.Errorf("add card at index %d: %w", i, err)
		}
		results = append(results, &MDAddedDTO{
			Title:       cardParam.Title,
			Description: cardParam.Description,
			Tag:         cardParam.Tag,
		})

	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit", "error", err)
		return nil, fmt.Errorf("commit: %w", err)
	}

	s.logger.Info("batch add completed", "count", len(results))

	return results, nil
}

// DeleteCard card from DB
func (s *Service) DeleteCard(ctx context.Context, cardId, userId int) error {
	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := tx.DeleteCard(ctx, cardId, userId); err != nil {
		s.logger.Error("failed to delete card", "error", err)
		return fmt.Errorf("delete card: %w", err)
	}

	return tx.Commit(ctx)
}

// UpdateCardDescription new description in DB
func (s *Service) UpdateCardDescription(ctx context.Context, cardId, UserID int, cardsUp Update) (*MindCard, error) {
	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	row, err := tx.UpdateCardDescription(ctx, cardId, UserID, cardsUp.NewDescription)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return rowToCard(row), nil
}

// UpdateLvl Возможно не понадобится
func (s *Service) UpdateLvl() {

}

// Get list of cards
func (s *Service) GetCards(ctx context.Context, userId int, limit, offset int16) ([]*MindCard, error) {
	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.GetCards(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return rowsToCards(rows), nil

}

// Get cards filtered by Tag
func (s *Service) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]*MindCard, error) {
	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.GetCardsByTag(ctx, tag, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return rowsToCards(rows), nil
}

// Get one card by ID
func (s *Service) GetCardById(ctx context.Context, cardId, userId int) (*MindCard, error) {
	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction", "error", err)
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	row, err := tx.GetCardById(ctx, cardId, userId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return rowToCard(row), nil
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
