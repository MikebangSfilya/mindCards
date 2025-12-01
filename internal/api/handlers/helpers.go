package handlers

import (
	dtoout "cards/internal/api/dto/dto_out"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type pagination struct {
	limit  int16
	offset int16
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

// help func to convert str to int
func (h *Handlers) stringToInt(in string) (int, error) {
	if in == "" {
		return 0, nil
	}

	return strconv.Atoi(in)
}

// helper to set default pagination variables
func (h *Handlers) limitOffset(limitStr, offsetStr string) (pagination, error) {

	limit, err := h.stringToInt(limitStr)
	if err != nil {
		return pagination{}, err
	}
	offset, err := h.stringToInt(offsetStr)
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
