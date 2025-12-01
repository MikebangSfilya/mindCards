package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) GetCards(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	p, err := h.limitOffset(limitStr, offsetStr)
	if err != nil {
		h.handleError(w, err, err.Error(), http.StatusBadRequest)
		return
	}

	cards, err := h.Service.GetCards(ctx, p.limit, p.offset)
	if err != nil {
		h.handleError(w, err, "failed to get cards", http.StatusInternalServerError)
		return
	}
	slog.Info("Get cards succeful")

	w.Header().Set("Content-Type", "application/json")
	if err := encoder(w, cards); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
		return
	}
}

func (h *Handlers) GetByTag(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	tag := chi.URLParam(r, "tag")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	p, err := h.limitOffset(limitStr, offsetStr)
	if err != nil {
		h.handleError(w, err, err.Error(), http.StatusBadRequest)
	}

	cards, err := h.Service.GetCardsByTag(ctx, tag, p.limit, p.offset)

	if err != nil {
		h.handleError(w, err, "failed to get cards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := encoder(w, cards); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetById(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	id := chi.URLParam(r, "id")

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

func (h *Handlers) GetEducation(w http.ResponseWriter, r *http.Request) {

}
