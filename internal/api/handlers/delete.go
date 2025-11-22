package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Delete card from DB
func (h *Handlers) DeleteCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	delId := chi.URLParam(r, "id")

	if err := h.Service.DeleteCard(ctx, delId); err != nil {
		h.handleError(w, err, ErrDeleteCard, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
