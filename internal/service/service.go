package service

import (
	"cards/internal/model"
	"context"
)

type Repo interface {
	AddCard(ctx context.Context, card model.MindCard) error
}

type Service struct {
	Repo Repo
}

func New(repo Repo) *Service {
	return &Service{
		Repo: repo,
	}
}

func (r *Service) AddCard(ctx context.Context, title, description, tag string) error {
	card := model.NewCard(title, description, tag)
	return r.Repo.AddCard(ctx, *card)
}
