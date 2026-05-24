package cards

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/MikebangSfilya/mindCards/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestServiceDeleteCardWithoutRedis(t *testing.T) {
	repo := &repoMock{}
	txManager := &txManagerMock{}
	txManager.On("WithinTransaction", mock.Anything, mock.AnythingOfType("func(context.Context, cards.DBQuerier) error")).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(context.Context, DBQuerier) error)
			_ = fn(args.Get(0).(context.Context), nil)
		}).
		Return(nil).
		Once()
	repo.On("DeleteCard", mock.Anything, mock.Anything, 11, 2).Return(nil).Once()

	service := NewService(repo, txManager, nil, testLogger(), nil)

	require.NoError(t, service.DeleteCard(context.Background(), 11, 2))
	repo.AssertExpectations(t)
	txManager.AssertExpectations(t)
}

func TestServiceUpdateCardDescriptionWithoutRedis(t *testing.T) {
	now := time.Now().UTC()
	repo := &repoMock{}
	repo.On("UpdateCardDescription", mock.Anything, mock.Anything, 11, 2, "Updated description").Return(storage.CardRow{
		CardID:      11,
		UserID:      2,
		Title:       "Go",
		Description: "Updated description",
		Tag:         "lang",
		CreatedAt:   now,
	}, nil).Once()

	service := NewService(repo, &txManagerMock{}, nil, testLogger(), nil)

	card, err := service.UpdateCardDescription(context.Background(), 11, 2, Update{NewDescription: "Updated description"})
	require.NoError(t, err)
	assert.Equal(t, int64(11), card.CardID)
	assert.Equal(t, "Updated description", card.Description)
	repo.AssertExpectations(t)
}

func TestServiceGetCardByIDWithoutRedisFallsBackToRepo(t *testing.T) {
	now := time.Now().UTC()
	repo := &repoMock{}
	repo.On("GetCardById", mock.Anything, mock.Anything, 15, 3).Return(storage.CardRow{
		CardID:      15,
		UserID:      3,
		Title:       "Title",
		Description: "Description",
		Tag:         "tag",
		CreatedAt:   now,
	}, nil).Once()

	service := NewService(repo, &txManagerMock{}, nil, testLogger(), nil)

	card, err := service.GetCardById(context.Background(), 15, 3)
	require.NoError(t, err)
	assert.Equal(t, int64(15), card.CardID)
	assert.Equal(t, "Description", card.Description)
	repo.AssertExpectations(t)
}
