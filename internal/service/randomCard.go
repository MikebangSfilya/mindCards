package service

import (
	"cards/internal/model"
	"cards/internal/storage"
	"context"
	"log/slog"
)

type RepoCard interface {
	GetCards(ctx context.Context, limit, offset int16) ([]storage.CardRow, error)
}

type CardsActiong struct {
	RepoCards RepoCard
	logger    *slog.Logger
}

func (s *CardsActiong) Educated(ctx context.Context, n int, limit, offset int16) ([]model.MindCard, error) {

	//take []cards from DB
	rows, err := s.RepoCards.GetCards(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	//take []cards from rows
	result := rowsToCards(rows)

	return result, nil
}
