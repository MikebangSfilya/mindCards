package cards

import (
	"context"
	"log/slog"
	"time"

	"github.com/MikebangSfilya/mindCards/internal/storage"
)

// интерфейс для связи с репозиторий

type Repo interface {
	AddCard(ctx context.Context, card *MindCard) error
	UpdateCardDescription(ctx context.Context, cardId, userId int, newDesc string) error
	DeleteCard(ctx context.Context, cardId, userId int) error
	GetCards(ctx context.Context, userId int, limit, offset int16) ([]storage.CardRow, error)
	GetCardById(ctx context.Context, cardId int) (storage.CardRow, error)
	GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]storage.CardRow, error)
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

// Add cards to DB
// TODO: collect errors
func (s *Service) AddCards(ctx context.Context, cardParams []Card) (*[]MDAddedDTO, error) {

	jobs := make(chan *MindCard, 50)
	results := make([]MDAddedDTO, 0, len(cardParams))

	go func() {
		defer close(jobs)
		for _, v := range cardParams {
			card, err := NewCard(v.Title, v.Description, v.Tag)
			if err != nil {
				s.logger.Error("failed to create card", "error", err)
				continue
			}

			cardCopy := *card

			jobs <- card

			results = append(results, MDAddedDTO{
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
func (s *Service) DeleteCard(ctx context.Context, cardId, userId int) error {
	return s.Repo.DeleteCard(ctx, cardId, userId)
}

// Update new description in DB
func (s *Service) UpdateCardDescription(ctx context.Context, cardId, UserID int, cardsUp Update) error {

	if err := s.Repo.UpdateCardDescription(ctx, cardId, UserID, cardsUp.NewDescription); err != nil {
		return err
	}

	return nil

}

// Возможно не понадобится
func (s *Service) UpdateLvl() {

}

// Get list of cards
func (s *Service) GetCards(ctx context.Context, userId int, limit, offset int16) ([]MindCard, error) {
	rows, err := s.Repo.GetCards(ctx, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	cards := rowsToCards(rows)

	return cards, nil

}

// Get cards filtered by Tag
func (s *Service) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]MindCard, error) {

	rows, err := s.Repo.GetCardsByTag(ctx, tag, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	cards := rowsToCards(rows)
	return cards, nil
}

// Get one card by unic ID
func (s *Service) GetCardById(ctx context.Context, cardId int) (MindCard, error) {
	row, err := s.Repo.GetCardById(ctx, cardId)
	if err != nil {
		return MindCard{}, err
	}

	return rowToCard(row), nil
}

func rowToCard(row storage.CardRow) MindCard {
	return MindCard{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Tag:         row.Tag,
		CreatedAt:   row.CreatedAt,
		LevelStudy:  row.LevelStudy,
		Learned:     row.Learned,
	}
}

func rowsToCards(rows []storage.CardRow) []MindCard {

	if rows == nil {
		return []MindCard{}
	}

	result := make([]MindCard, len(rows))

	for i, row := range rows {
		result[i] = rowToCard(row)
	}

	return result

}
