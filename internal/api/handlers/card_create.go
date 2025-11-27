package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"context"
	"net/http"
)

// AddCard handler for add card in DB
func (h *Handlers) AddCard(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), baseTimeOut)
	defer cancel()

	var DTOin []dtoin.Card
	if err := decoder(r, &DTOin); err != nil {
		h.handleError(w, err, ErrDecodeJSON, http.StatusBadRequest)
		return
	}
	result := make([]dtoout.MindCardDTO, 0, len(DTOin))
	for i := range DTOin {
		cardDTO, err := h.Service.AddCard(ctx, DTOin[i])
		if err != nil {
			h.handleError(w, err, ErrAddCard, http.StatusBadRequest)
			return
		}
		result = append(result, *cardDTO)
	}

	if err := encoder(w, result); err != nil {
		h.handleError(w, err, ErrEncodeJSON, http.StatusBadRequest)
	}

}
