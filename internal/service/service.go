package service

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"cards/internal/storage"
	"context"
	"fmt"
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

func (s *Service) DeleteCard(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("card not exist")
	}

	return s.Repo.DeleteCard(ctx, id)
}

func (s *Service) UpdateCardDescription(ctx context.Context, id string, cardsUp dtoin.Update) error {
	if id == "" {
		return fmt.Errorf("nil id")
	}
	if cardsUp.NewDeccription == "" {
		return fmt.Errorf("nil desc")
	}

	if err := s.Repo.UptadeCardDescription(ctx, id, cardsUp.NewDeccription); err != nil {
		return err
	}

	return nil

}

func (s *Service) GetCards(ctx context.Context, pagination dtoin.LimitOffset) (map[string]model.MindCard, error) {
	return s.Repo.GetCards(ctx, pagination.Limit, pagination.Offset)
}

func (s *Service) GetCardsByTag(ctx context.Context, tag string, pagination dtoin.LimitOffset) (map[string]model.MindCard, error) {

	rows, err := s.Repo.GetCardsByTag(ctx, tag, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}

	cards := make(map[string]model.MindCard)

	for _, row := range rows {
		card := model.MindCard{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Tag:         row.Tag,
			CreatedAt:   row.CreatedAt,
			LevelStudy:  row.LevelStudy,
			Learned:     row.Learned,
		}
		cards[fmt.Sprintf("%d", card.ID)] = card
	}
	return cards, nil
}

func (s *Service) GetCardById(ctx context.Context, id string) (model.MindCard, error) {
	row, err := s.Repo.GetCardById(ctx, id)
	if err != nil {
		return model.MindCard{}, err
	}

	return model.MindCard{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Tag:         row.Tag,
		CreatedAt:   row.CreatedAt,
		LevelStudy:  row.LevelStudy,
		Learned:     row.Learned,
	}, nil

}
