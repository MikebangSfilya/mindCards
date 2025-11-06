package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handlers interface {
	Add(ctx context.Context)
	Delete(ctx context.Context)
	Update(ctx context.Context)
}

type Server struct {
	handle handlers
}

func NewServer(handl handlers) *Server {
	return &Server{
		handle: handl,
	}
}

func (s *Server) Start() {
	r := chi.NewRouter()
	r.Get("/card", func(w http.ResponseWriter, r *http.Request) {}) // заглушка

}
