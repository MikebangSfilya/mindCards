package handlers

import (
	dtoin "cards/internal/api/dto/dto_in"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service interface {
	AddCard(ctx context.Context, title, description, tag string) error
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

	var DTOin dtoin.DTOin
	if err := json.NewDecoder(r.Body).Decode(&DTOin); err != nil {
		errDTO := fmt.Errorf("error %v", err)
		http.Error(w, errDTO.Error(), http.StatusBadRequest)
		return
	}

	if err := h.HTTPhandle.AddCard(ctx, DTOin.Title, DTOin.Description, DTOin.Tag); err != nil {
		errDTO := fmt.Errorf("error %v", err)
		http.Error(w, errDTO.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (h *Handlers) DeleteCard(w http.ResponseWriter, r *http.Request) {

}
func (h *Handlers) UpdateCard(w http.ResponseWriter, r *http.Request) {

}
