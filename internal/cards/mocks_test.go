package cards

import (
	"context"

	"github.com/MikebangSfilya/mindCards/internal/storage"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (m *serviceMock) AddCards(ctx context.Context, userId int, cardParams []Card) ([]*MDAddedDTO, error) {
	args := m.Called(ctx, userId, cardParams)
	if got := args.Get(0); got != nil {
		return got.([]*MDAddedDTO), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *serviceMock) DeleteCard(ctx context.Context, cardId, userId int) error {
	args := m.Called(ctx, cardId, userId)
	return args.Error(0)
}

func (m *serviceMock) UpdateCardDescription(ctx context.Context, cardId, userId int, cardsUp Update) (*MindCard, error) {
	args := m.Called(ctx, cardId, userId, cardsUp)
	if got := args.Get(0); got != nil {
		return got.(*MindCard), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *serviceMock) GetCards(ctx context.Context, userId int, limit, offset int16) ([]*MindCard, error) {
	args := m.Called(ctx, userId, limit, offset)
	if got := args.Get(0); got != nil {
		return got.([]*MindCard), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *serviceMock) GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]*MindCard, error) {
	args := m.Called(ctx, tag, userId, limit, offset)
	if got := args.Get(0); got != nil {
		return got.([]*MindCard), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *serviceMock) GetCardById(ctx context.Context, cardId, userId int) (*MindCard, error) {
	args := m.Called(ctx, cardId, userId)
	if got := args.Get(0); got != nil {
		return got.(*MindCard), args.Error(1)
	}
	return nil, args.Error(1)
}

type repoMock struct {
	mock.Mock
}

func (m *repoMock) AddCard(ctx context.Context, db DBQuerier, userId int, card *MindCard) error {
	args := m.Called(ctx, db, userId, card)
	return args.Error(0)
}

func (m *repoMock) UpdateCardDescription(ctx context.Context, db DBQuerier, cardId, userId int, newDesc string) (storage.CardRow, error) {
	args := m.Called(ctx, db, cardId, userId, newDesc)
	if got := args.Get(0); got != nil {
		return got.(storage.CardRow), args.Error(1)
	}
	return storage.CardRow{}, args.Error(1)
}

func (m *repoMock) DeleteCard(ctx context.Context, db DBQuerier, cardId, userId int) error {
	args := m.Called(ctx, db, cardId, userId)
	return args.Error(0)
}

func (m *repoMock) GetCards(ctx context.Context, db DBQuerier, userId int, limit, offset int16) ([]storage.CardRow, error) {
	args := m.Called(ctx, db, userId, limit, offset)
	if got := args.Get(0); got != nil {
		return got.([]storage.CardRow), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *repoMock) GetCardsByTag(ctx context.Context, db DBQuerier, tag string, userId int, limit, offset int16) ([]storage.CardRow, error) {
	args := m.Called(ctx, db, tag, userId, limit, offset)
	if got := args.Get(0); got != nil {
		return got.([]storage.CardRow), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *repoMock) GetCardById(ctx context.Context, db DBQuerier, cardId, userId int) (storage.CardRow, error) {
	args := m.Called(ctx, db, cardId, userId)
	if got := args.Get(0); got != nil {
		return got.(storage.CardRow), args.Error(1)
	}
	return storage.CardRow{}, args.Error(1)
}

type txManagerMock struct {
	mock.Mock
}

func (m *txManagerMock) WithinTransaction(ctx context.Context, fn func(ctx context.Context, q DBQuerier) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
