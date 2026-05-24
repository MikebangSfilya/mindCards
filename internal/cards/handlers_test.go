package cards

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDecoder(t *testing.T) {
	dto := Card{
		Title:       "title",
		Description: "Desc",
		Tag:         "Go",
	}
	testCase, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(testCase))
	dtoIn := Card{}
	err := decoder(req, &dtoIn)
	require.NoError(t, err)
	require.NotNil(t, dtoIn)
	assert.Equal(t, dto.Title, dtoIn.Title)
	assert.Equal(t, dto.Description, dtoIn.Description)

	assert.Equal(t, dto.Tag, dtoIn.Tag)
}

func TestBase(t *testing.T) {
}

func TestGetByTag_InvalidPaginationStopsRequest(t *testing.T) {
	svc := &serviceMock{}
	handler := New(svc)

	router := chi.NewRouter()
	handler.RegistredRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/card/tag/go?limit=bad", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	svc.AssertNotCalled(t, "GetCardsByTag", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	assert.Contains(t, rec.Body.String(), "invalid syntax")
}

func TestAddCards_ServiceErrorReturnsServerError(t *testing.T) {
	svc := &serviceMock{}
	svc.On("AddCards", mock.Anything, 1, []Card{{
		Title:       "Go",
		Description: "Long enough description",
		Tag:         "lang",
	}}).Return(nil, errors.New("db failure")).Once()
	handler := New(svc)

	router := chi.NewRouter()
	handler.RegistredRoutes(router)

	body := `[{"title":"Go","description":"Long enough description","tag":"lang"}]`
	req := httptest.NewRequest(http.MethodPost, "/card/", bytes.NewBufferString(body))
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "db failure")
	svc.AssertExpectations(t)
}

func TestUpdateCard_WritesValidJSONResponse(t *testing.T) {
	now := time.Now().UTC()
	svc := &serviceMock{}
	svc.On("UpdateCardDescription", mock.Anything, 7, 1, Update{
		NewDescription: "Updated description",
	}).Return(&MindCard{
		CardID:      7,
		UserID:      1,
		Title:       "Go",
		Description: "Updated description",
		Tag:         "lang",
		CreatedAt:   now,
	}, nil).Once()
	handler := New(svc)

	router := chi.NewRouter()
	handler.RegistredRoutes(router)

	req := httptest.NewRequest(http.MethodPut, "/card/7", bytes.NewBufferString(`{"description":"Updated description"}`))
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got MindCard
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, int64(7), got.CardID)
	assert.Equal(t, "Updated description", got.Description)
	svc.AssertExpectations(t)
}
