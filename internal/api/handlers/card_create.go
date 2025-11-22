package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	"context"
	"net/http"
)

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

}
