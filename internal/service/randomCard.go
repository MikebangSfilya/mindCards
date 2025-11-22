package service

import (
	"cards/internal/storage"
	"context"
)

type RepoCard interface {
	GetCardById(ctx context.Context, id string) (storage.CardRow, error)
}

type CardsActiong struct {
	RepoCards RepoCard
}
