package service

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"context"
	"fmt"
	"log/slog"
	"sync"
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

func (s *CardCRUDService) AddSliceCard(ctx context.Context, cardParams []dtoin.Card) {

}

// Add card to DB
func (s *CardCRUDService) AddCard(ctx context.Context, cardsParams dtoin.Card) (*dtoout.MDAddedDTO, error) {
	var card *model.MindCard
	var cardErr error
	var wg sync.WaitGroup
	var mu sync.RWMutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		card, cardErr = model.NewCard(cardsParams.Title, cardsParams.Description, cardsParams.Tag)
		if cardErr != nil {
			s.logger.Error("failed to add card", "error", cardErr)
			return
		}
	}()

	wg.Wait()

	if err := s.Repo.AddCard(ctx, card); err != nil {
		s.logger.Error("failed to add card", "error", err)
		return nil, err
	}

	s.logger.Info("adding card", "title", cardsParams.Title)

	return &dtoout.MDAddedDTO{
		ID:          card.ID,
		Title:       card.Title,
		Description: card.Description,
		Tag:         card.Tag,
	}, nil
}

// Delete card from DB
func (s *CardCRUDService) DeleteCard(ctx context.Context, id string) error {
	if id == "" {
		s.logger.Warn("failed to delete card", "Warn", ErrNotExist)
		return ErrNotExist
	}

	return s.Repo.DeleteCard(ctx, id)
}

// Update new description in DB
func (s *CardCRUDService) UpdateCardDescription(ctx context.Context, id string, cardsUp dtoin.Update) error {
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

// Возможно не понадобится
func (s *CardCRUDService) UpdateLvl() {

}

// Get list of cards
func (s *CardCRUDService) GetCards(ctx context.Context, limit, offset int16) (map[string]model.MindCard, error) {
	return s.Repo.GetCards(ctx, limit, offset)
}

// Get cards filtered by Tag
func (s *CardCRUDService) GetCardsByTag(ctx context.Context, tag string, limit, offset int16) (map[string]model.MindCard, error) {

	rows, err := s.Repo.GetCardsByTag(ctx, tag, limit, offset)
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

// Get one card by unic ID
func (s *CardCRUDService) GetCardById(ctx context.Context, id string) (model.MindCard, error) {
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
