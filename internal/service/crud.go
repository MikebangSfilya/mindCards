package service

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type CardCRUDService struct {
	Repo   Repo
	logger *slog.Logger
}

func NewCardCRUDService(repo Repo, logger *slog.Logger) *CardCRUDService {

	return &CardCRUDService{
		Repo:   repo,
		logger: logger,
	}
}

// Add cards to DB
func (s *CardCRUDService) AddCards(ctx context.Context, cardParams []dtoin.Card) (*[]dtoout.MDAddedDTO, error) {

	jobs := make(chan *model.MindCard, 50)
	results := make([]dtoout.MDAddedDTO, 0, len(cardParams))

	go func() {
		defer close(jobs)
		for _, v := range cardParams {

			card, err := model.NewCard(v.Title, v.Description, v.Tag)
			if err != nil {
				s.logger.Error("failed to create card", "error", err)

				continue
			}

			cardCopy := *card

			jobs <- card

			results = append(results, dtoout.MDAddedDTO{
				Title:       cardCopy.Title,
				Description: cardCopy.Description,
				Tag:         cardCopy.Tag,
			})
		}
	}()

	for job := range jobs {
		go func() {
			dbContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.Repo.AddCard(dbContext, job); err != nil {
				s.logger.Error("failed to add card", "error", err)

			}

			s.logger.Info("adding card", "title", job.Title)
		}()
	}

	return &results, nil
}

// Delete card from DB
func (s *CardCRUDService) DeleteCard(ctx context.Context, id string) error {
	if id == "" {
		s.logger.Warn("failed to delete card", "error", ErrNotExist)
		return ErrNotExist
	}

	return s.Repo.DeleteCard(ctx, id)
}

// Update new description in DB
func (s *CardCRUDService) UpdateCardDescription(ctx context.Context, id string, cardsUp dtoin.Update) error {
	if id == "" {
		return fmt.Errorf("nil id")
	}
	if cardsUp.NewDescription == "" {
		return fmt.Errorf("nil desc")
	}

	if err := s.Repo.UpdateCardDescription(ctx, id, cardsUp.NewDescription); err != nil {
		return err
	}

	return nil

}

// Возможно не понадобится
func (s *CardCRUDService) UpdateLvl() {

}

// Get list of cards
func (s *CardCRUDService) GetCards(ctx context.Context, limit, offset int16) ([]model.MindCard, error) {
	rows, err := s.Repo.GetCards(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	cards := rowsToCard(rows)

	return cards, nil

}

// Get cards filtered by Tag
func (s *CardCRUDService) GetCardsByTag(ctx context.Context, tag string, limit, offset int16) ([]model.MindCard, error) {

	rows, err := s.Repo.GetCardsByTag(ctx, tag, limit, offset)
	if err != nil {
		return nil, err
	}

	cards := rowsToCard(rows)
	return cards, nil
}

// Get one card by unic ID
func (s *CardCRUDService) GetCardById(ctx context.Context, id string) (model.MindCard, error) {
	row, err := s.Repo.GetCardById(ctx, id)
	if err != nil {
		return model.MindCard{}, err
	}

	return rowToCard(row), nil
}
