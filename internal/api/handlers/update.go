package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Update card in DB
func (h *Handlers) UpdateCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	upId := chi.URLParam(r, "id")

	dtoUp := dtoin.Update{}
	if err := decoder(r, &dtoUp); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}

	if err := dtoUp.Validate(); err != nil {
		h.handleError(w, err, "validate error", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateCardDescription(ctx, upId, dtoUp); err != nil {
		h.handleError(w, err, ErrUpdateCard, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
