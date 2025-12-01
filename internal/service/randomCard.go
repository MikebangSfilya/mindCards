package service

import (
	"cards/internal/model"
	"cards/internal/storage"
	"context"
	"log/slog"
)

type RepoCard interface {
	//not real realisaion
	GetCards(ctx context.Context, eduParams eduParams, limit, offset int16) ([]storage.CardRow, error)
}

type CardsActiong struct {
	RepoCards RepoCard
	logger    *slog.Logger
}

type eduParams struct {
}

func (s *CardsActiong) Educated(ctx context.Context, n int, limit, offset int16) ([]model.MindCard, error) {

	eduParams := eduParams{}

	//take []cards from DB
	rows, err := s.RepoCards.GetCards(ctx, eduParams, limit, offset)
	if err != nil {
		return nil, err
	}

	//take []cards from rows
	result := rowsToCard(rows)

	return result, nil
}
