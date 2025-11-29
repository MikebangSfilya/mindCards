package service

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"cards/internal/storage"
	"context"
	"log/slog"
)

// интерфейс для связи с репозиторий

type Repo interface {
	AddCard(ctx context.Context, card *model.MindCard) error
	UptadeCardDescription(ctx context.Context, id, newDesc string) error
	DeleteCard(ctx context.Context, id string) error
	GetCards(ctx context.Context, limit, offset int16) (map[string]model.MindCard, error)
	GetCardById(ctx context.Context, id string) (storage.CardRow, error)
	GetCardsByTag(ctx context.Context, tag string, limit, offset int16) ([]storage.CardRow, error)
}

type Service struct {
	Crud   *CardCRUDService
	logger *slog.Logger
	// RandCard RepoCard
}

func New(repo Repo, logger *slog.Logger) *Service {
	serviceLogger := logger.With("component", "service")
	return &Service{
		Crud:   NewCardCRUDService(repo, serviceLogger),
		logger: serviceLogger,
		// RandCard: randCard,
	}
}

// Add card to DB
func (s *Service) AddCard(ctx context.Context, cardsParams dtoin.Card) (*dtoout.MDAddedDTO, error) {
	return s.Crud.AddCard(ctx, cardsParams)
}

// Delete card from DB
func (s *Service) DeleteCard(ctx context.Context, id string) error {
	return s.Crud.DeleteCard(ctx, id)
}

// Update new description in DB
func (s *Service) UpdateCardDescription(ctx context.Context, id string, cardsUp dtoin.Update) error {
	return s.Crud.UpdateCardDescription(ctx, id, cardsUp)
}

// Возможно не понадобится
func (s *Service) UpdateLvl() {

}

// Get list of cards
func (s *Service) GetCards(ctx context.Context, limit, offset int16) (map[string]model.MindCard, error) {
	return s.Crud.GetCards(ctx, int16(limit), int16(offset))
}

// Get cards filtered by Tag
func (s *Service) GetCardsByTag(ctx context.Context, tag string, limit, offset int16) (map[string]model.MindCard, error) {
	return s.Crud.GetCardsByTag(ctx, tag, limit, offset)
}

// Get one card by unic ID
func (s *Service) GetCardById(ctx context.Context, id string) (model.MindCard, error) {
	return s.Crud.GetCardById(ctx, id)
}
