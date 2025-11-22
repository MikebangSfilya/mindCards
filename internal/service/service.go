package service

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"context"
	"fmt"
	"log/slog"
)

// интерфейс для связи с репозиторий

type Repo interface {
	AddCard(ctx context.Context, card *model.MindCard) error
	UptadeCardDescription(ctx context.Context, updt []string)
	DeleteCard(ctx context.Context, title string) error
	GetCards(ctx context.Context, limit, offset int16) (map[string]model.MindCard, error)
	GetCardById(ctx context.Context, id int) model.MindCard
	GetCardsByTag(ctx context.Context, tag string, limit, offset int16) (map[string]model.MindCard, error)
}

type Service struct {
	Repo Repo
}

func New(repo Repo) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) AddCard(ctx context.Context, cardsParams dtoin.Card) (*dtoout.MindCardDTO, error) {
	card, err := model.NewCard(cardsParams.Title, cardsParams.Description, cardsParams.Tag)
	if err != nil {
		return nil, err
	}
	if err := s.Repo.AddCard(ctx, card); err != nil {
		slog.Error("failed to add card",
			"error", err,
			"package", "service")
		return nil, err
	}

	return &dtoout.MindCardDTO{
		ID:          card.ID,
		Title:       card.Title,
		Description: card.Description,
		Tag:         card.Tag,
		CreatedAt:   card.CreatedAt,
		LevelStudy:  card.LevelStudy,
		Learned:     card.Learned,
	}, nil
}

func (s *Service) DeleteCard(ctx context.Context, title string) error {
	if title == "" {
		return fmt.Errorf("card not exist")
	}

	return s.Repo.DeleteCard(ctx, title)
}

func (s *Service) UpdateCardDescription(ctx context.Context, cardsUp dtoin.Update) error {
	if cardsUp.Title == "" {
		return fmt.Errorf("nil title")
	}
	if cardsUp.NewDeccription == "" {
		return fmt.Errorf("nil desc")
	}

	cardToRepo := []string{cardsUp.Title, cardsUp.NewDeccription}

	s.Repo.UptadeCardDescription(ctx, cardToRepo)

	return nil

}

func (s *Service) GetCards(ctx context.Context, pagination dtoin.LimitOffset) (map[string]model.MindCard, error) {
	return s.Repo.GetCards(ctx, pagination.Limit, pagination.Offset)
}
func (s *Service) GetCardsByTag(ctx context.Context, tag string, pagination dtoin.LimitOffset) (map[string]model.MindCard, error) {
	return s.Repo.GetCardsByTag(ctx, tag, pagination.Limit, pagination.Offset)
}

func (s *Service) GetCardById(ctx context.Context, id int) model.MindCard {
	return s.Repo.GetCardById(ctx, id)
}
