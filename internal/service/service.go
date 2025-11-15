package service

import (
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
}

type Service struct {
	Repo Repo
}

func New(repo Repo) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) AddCard(ctx context.Context, title, description, tag string) error {
	card := model.NewCard(title, description, tag)

	if err := s.Repo.AddCard(ctx, card); err != nil {
		slog.Error("failed to add card",
			"error", err,
			"package", "service")
		return err
	}

	return nil
}

func (s *Service) DeleteCard(ctx context.Context, title string) error {
	if title == "" {
		return fmt.Errorf("card not exist")
	}

	return s.Repo.DeleteCard(ctx, title)
}

func (s *Service) UptadeCardDescription(ctx context.Context, title, description string) error {
	if title == "" {
		return fmt.Errorf("nil title")
	}
	if description == "" {
		return fmt.Errorf("nil desc")
	}
	//Упросить
	cardToRepo := []string{title, description}

	s.Repo.UptadeCardDescription(ctx, cardToRepo)

	return nil

}
