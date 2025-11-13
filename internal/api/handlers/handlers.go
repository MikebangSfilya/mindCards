package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	dtoout "cards/internal/api/dto/dto_out"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Service interface {
	AddCard(ctx context.Context, title, description, tag string) error
	DeleteCard(ctx context.Context, title string) error
}

type Handle interface {
	AddCard(w http.ResponseWriter, r *http.Request)
	DeleteCard(w http.ResponseWriter, r *http.Request)
	UpdateCard(w http.ResponseWriter, r *http.Request)
}

type Handlers struct {
	HTTPhandle Service
}

func New(handle Service) *Handlers {
	return &Handlers{
		HTTPhandle: handle,
	}
}

func (h *Handlers) AddCard(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	var DTOin dtoin.Card
	if err := json.NewDecoder(r.Body).Decode(&DTOin); err != nil {
		slog.Error("failed to decode json",
			"err", err,
			"package", "handlers")
		errDTO := dtoout.NewErr(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.HTTPhandle.AddCard(ctx, DTOin.Title, DTOin.Description, DTOin.Tag); err != nil {
		slog.Error("failed to add card",
			"err", err,
			"package", "handlers")
		errDTO := dtoout.NewErr(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *Handlers) DeleteCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	var dtoDel dtoin.DTODel
	if err := json.NewDecoder(r.Body).Decode(&dtoDel); err != nil {
		slog.Error("failed to decode json",
			"err", err,
			"package", "handlers")
		errDTO := dtoout.NewErr(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := h.HTTPhandle.DeleteCard(ctx, dtoDel.Title); err != nil {
		slog.Error("failed to delete card",
			"err", err,
			"package", "handlers")
		errDTO := dtoout.NewErr(err)
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	resp := dtoout.NewDelDTO(dtoDel.Title)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Warn("wtf?")
		return
	}

}
func (h *Handlers) UpdateCard(w http.ResponseWriter, r *http.Request) {

}
