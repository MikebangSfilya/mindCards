package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetCards(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err := h.stringToInt(limitStr)
	if err != nil {
		h.handleError(w, fmt.Errorf("invalid limit parameter: %w", err), "Invalid limit", http.StatusBadRequest)
		return
	}
	offset, err := h.stringToInt(offsetStr)
	if err != nil {
		h.handleError(w, fmt.Errorf("invalid offset parameter: %w", err), "Invalid limit", http.StatusBadRequest)
		return
	}

	if limit == 0 {
		limit = 50
	} else if limit > 1000 {
		limit = 1000
	}

	if offset < 0 {
		offset = 0
	}

	cards, err := h.Service.GetCards(ctx, int16(limit), int16(offset))
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

func (h *Handlers) GetById(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	card, err := h.Service.GetCardById(ctx, id)
	if err != nil {
		h.handleError(w, err, "potom", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := encoder(w, card); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusInternalServerError)
		return
	}

}
