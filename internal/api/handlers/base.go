package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"context"
	"time"
)

const (
	baseTimeOut = time.Duration(15 * time.Second)
)

const (
	ErrDecodeJSON = "failed to decode JSON"
	ErrEncodeJSON = "failed to encode response"
	ErrAddCard    = "failed to add card"
	ErrDeleteCard = "failed to delete card"
	ErrUpdateCard = "failed to update card"
)

// Service is an interface for connecting to the service layer
type Service interface {
	AddCard(ctx context.Context, cardsParams dtoin.Card) (*dtoout.MindCardDTO, error)
	DeleteCard(ctx context.Context, id string) error
	UpdateCardDescription(ctx context.Context, id string, cardsUp dtoin.Update) error
	GetCards(ctx context.Context, pagination dtoin.LimitOffset) (map[string]model.MindCard, error)
	GetCardsByTag(ctx context.Context, tag string, pagination dtoin.LimitOffset) (map[string]model.MindCard, error)
	GetCardById(ctx context.Context, id string) (model.MindCard, error)
}

// Handlers stores the service layer dependency
type Handlers struct {
	Service Service
}

func New(service Service) *Handlers {
	return &Handlers{
		Service: service,
	}
}
