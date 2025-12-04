package cards

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/MikebangSfilya/mindCards/internal/auth"
	"github.com/go-chi/chi/v5"
)

const (
	baseTimeOut    = time.Duration(15 * time.Second)
	addCardTimeOut = time.Duration(60 * time.Second)
)

// msg errors
const (
	ErrDecodeJSON = "failed to decode JSON"
	ErrEncodeJSON = "failed to encode response"
	ErrAddCard    = "failed to add card"
	ErrDeleteCard = "failed to delete card"
	ErrUpdateCard = "failed to update card"
	ErrValidate   = "wrong params"
)

// ServiceRepo is an interface for connecting to the service layer
type ServiceRepo interface {
	AddCards(ctx context.Context, userId int, cardParams []Card) (*[]MDAddedDTO, error)
	DeleteCard(ctx context.Context, cardId, userId int) error
	UpdateCardDescription(ctx context.Context, cardId, userId int, cardsUp Update) error
	GetCards(ctx context.Context, userId int, limit, offset int16) ([]MindCard, error)
	GetCardsByTag(ctx context.Context, tag string, userId int, limit, offset int16) ([]MindCard, error)
	GetCardById(ctx context.Context, id int) (MindCard, error)
}

type pagination struct {
	limit  int16
	offset int16
}

// Handler stores the service layer dependency
type Handler struct {
	Service ServiceRepo
}

func New(service ServiceRepo) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) AddCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), addCardTimeOut)
		defer cancel()

		usId, err := auth.GetUserID(ctx)
		if err != nil {
			h.handleError(w, err, "authentication required", http.StatusUnauthorized)
			return
		}

		var CardIn []Card
		if err := decoder(r, &CardIn); err != nil {
			h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
			return
		}
		for _, v := range CardIn {
			if err := v.Validate(); err != nil {
				h.handleError(w, err, ErrValidate, http.StatusBadRequest)
				return
			}
		}

		result, err := h.Service.AddCards(ctx, usId, CardIn)
		if err != nil {
			h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
			return
		}

		if err := encoder(w, result); err != nil {
			h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
		}

	}
}

func (h *Handler) DeleteCard() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
		defer cancel()

		delIdstr := chi.URLParam(r, "id")
		delId, err := strconv.Atoi(delIdstr)
		if err != nil {
			h.handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}

		usId, err := auth.GetUserID(ctx)
		if err != nil {
			h.handleError(w, err, "authentication required", http.StatusUnauthorized)
			return
		}

		if err := h.Service.DeleteCard(ctx, delId, usId); err != nil {
			h.handleError(w, err, ErrDeleteCard, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

func (h *Handler) GetCards() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
		defer cancel()

		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		p, err := h.limitOffset(limitStr, offsetStr)
		if err != nil {
			h.handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}

		usId, err := auth.GetUserID(ctx)
		if err != nil {
			h.handleError(w, err, "authentication required", http.StatusUnauthorized)
			return
		}

		cards, err := h.Service.GetCards(ctx, usId, p.limit, p.offset)
		if err != nil {
			h.handleError(w, err, errFailToAdd.Error(), http.StatusInternalServerError)
			return
		}
		slog.Info("Get cards succeful")

		w.Header().Set("Content-Type", "application/json")
		if err := encoder(w, cards); err != nil {
			h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
			return
		}
	}
}

func (h *Handler) GetByTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
		defer cancel()

		usId, err := auth.GetUserID(ctx)
		if err != nil {
			h.handleError(w, err, "authentication required", http.StatusUnauthorized)
			return
		}

		tag := chi.URLParam(r, "tag")

		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		p, err := h.limitOffset(limitStr, offsetStr)
		if err != nil {
			h.handleError(w, err, ErrValidate, http.StatusBadRequest)
		}

		cards, err := h.Service.GetCardsByTag(ctx, tag, usId, p.limit, p.offset)

		if err != nil {
			h.handleError(w, err, errFailToAdd.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := encoder(w, cards); err != nil {
			h.handleError(w, err, ErrEncodeJSON, http.StatusInternalServerError)
			return
		}

	}
}

func (h *Handler) GetById() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
		defer cancel()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}

		card, err := h.Service.GetCardById(ctx, id)
		if err != nil {
			h.handleError(w, err, "failed to get cards", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := encoder(w, card); err != nil {
			h.handleError(w, err, ErrEncodeJSON, http.StatusInternalServerError)
			return
		}

	}
}

func (h *Handler) UpdateCard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
		defer cancel()

		upIdStr := chi.URLParam(r, "id")

		usId, err := auth.GetUserID(ctx)
		if err != nil {
			h.handleError(w, err, "authentication required", http.StatusUnauthorized)
			return
		}

		upId, err := strconv.Atoi(upIdStr)
		if err != nil {
			h.handleError(w, err, err.Error(), http.StatusBadRequest)
			return
		}

		if err := decoder(r, &usId); err != nil {
			h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		}

		dtoUp := Update{}
		if err := decoder(r, &dtoUp); err != nil {
			h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
			return
		}

		if err := dtoUp.Validate(); err != nil {
			h.handleError(w, err, "validate error", http.StatusBadRequest)
			return
		}

		if err := h.Service.UpdateCardDescription(ctx, upId, upId, dtoUp); err != nil {
			h.handleError(w, err, ErrUpdateCard, http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) RegistredRoutes(r chi.Router) {
	r.Route("/card", func(r chi.Router) {
		r.Use(auth.AuthenticateUser)

		r.Post("/", h.AddCards())         //add card
		r.Delete("/{id}", h.DeleteCard()) // Delete card
		r.Put("/{id}", h.UpdateCard())    // Update card
		r.Get("/tag/{tag}", h.GetByTag()) // Get by tag
		r.Get("/", h.GetCards())          // Get all card, limit and offset get by QUERY
		r.Get("/{id}", h.GetById())       // get by unic ID
	})
}

// help func for decode json
func decoder(r *http.Request, dto any) error {
	return json.NewDecoder(r.Body).Decode(dto)
}

// help func for encode json
func encoder(w http.ResponseWriter, resp any) error {
	return json.NewEncoder(w).Encode(resp)
}

// help func responce error DTO
func (h *Handler) handleError(w http.ResponseWriter, err error, msg string, code int) {
	slog.Error(msg, "err", err, "package", "handlers")
	errDTO := NewErr(err)
	http.Error(w, errDTO.ToString(), code)
}

func strToI(s string) (int16, error) {
	if s == "" {
		return 0, nil
	}

	out, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return int16(out), nil
}

// helper to set default pagination variables
func (h *Handler) limitOffset(limitStr, offsetStr string) (pagination, error) {

	limit, err := strToI(limitStr)
	if err != nil {
		return pagination{}, err
	}
	offset, err := strToI(offsetStr)
	if err != nil {

		return pagination{}, err
	}

	p := pagination{
		limit:  int16(limit),
		offset: int16(offset),
	}

	if p.limit == 0 {
		p.limit = 50
	} else if limit > 1000 {
		p.limit = 1000
	}

	if p.offset < 0 {
		p.offset = 0
	}

	return p, nil
}
