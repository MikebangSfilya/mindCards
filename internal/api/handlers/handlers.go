package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"cards/internal/model"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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
	DeleteCard(ctx context.Context, title string) error
	UpdateCardDescription(ctx context.Context, cardsUp dtoin.Update) error
	GetCards(ctx context.Context, pagination dtoin.LimitOffset) (map[string]model.MindCard, error)
	GetCardsByTag(ctx context.Context, tag string, pagination dtoin.LimitOffset) (map[string]model.MindCard, error)
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

// AddCard handler for add card in DB
func (h *Handlers) AddCard(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	var DTOin dtoin.Card
	if err := decoder(r, &DTOin); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}

	cardDTO, err := h.Service.AddCard(ctx, DTOin)
	if err != nil {
		h.handleError(w, err, ErrAddCard, http.StatusBadRequest)
		return
	}
	if err := encoder(w, cardDTO); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
	}

	// w.Write([]byte("dada"))
	// w.WriteHeader(http.StatusOK)
}

// Delete card from DB
func (h *Handlers) DeleteCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	var dtoDel dtoin.DTODel
	if err := decoder(r, &dtoDel); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteCard(ctx, dtoDel.Title); err != nil {
		h.handleError(w, err, ErrDeleteCard, http.StatusBadRequest)
		return
	}
	resp := dtoout.NewDelDTO(dtoDel.Title)
	if err := encoder(w, resp); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
		return
	}

}

// Update card in DB
func (h *Handlers) UpdateCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	dtoUp := dtoin.Update{}
	if err := decoder(r, &dtoUp); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateCardDescription(ctx, dtoUp); err != nil {
		h.handleError(w, err, ErrUpdateCard, http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetCards(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	dtoIn := dtoin.LimitOffset{}
	if err := decoder(r, &dtoIn); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}
	dtoIn.PaginationDefault()

	cards, err := h.Service.GetCards(ctx, dtoIn)
	if err != nil {
		h.handleError(w, err, "potom", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := encoder(w, cards); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
		return
	}
}

func (h *Handlers) GetByTag(w http.ResponseWriter, r *http.Request) {

	tag := chi.URLParam(r, "tag")

	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	dtoIn := dtoin.LimitOffset{}
	if err := decoder(r, &dtoIn); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}

	dtoIn.PaginationDefault()

	cards, err := h.Service.GetCardsByTag(ctx, tag, dtoIn)
	if err != nil {
		h.handleError(w, err, "potom", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := encoder(w, cards); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusInternalServerError)
		return
	}

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
func (h *Handlers) handleError(w http.ResponseWriter, err error, msg string, code int) {
	slog.Error(msg, "err", err, "package", "handlers")
	errDTO := dtoout.NewErr(err)
	http.Error(w, errDTO.ToString(), code)
}
