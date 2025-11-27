package handlers

import (
	dtoout "cards/internal/api/dto/dto_out"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

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

func (h *Handlers) stringToInt(in string) (int, error) {
	if in == "" {
		return 0, nil
	}

	return strconv.Atoi(in)
}
