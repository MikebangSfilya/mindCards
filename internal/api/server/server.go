package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handlers interface {
	AddCard(w http.ResponseWriter, r *http.Request)
	DeleteCard(w http.ResponseWriter, r *http.Request)
	UpdateCard(w http.ResponseWriter, r *http.Request)
}
type Server struct {
	handlers handlers
	server   *http.Server
}

func NewServer(handl handlers) *Server {
	return &Server{
		handlers: handl,
	}
}

func (s *Server) Start() error {

	port := ":8080"

	r := chi.NewRouter()
	r.Route("/card", func(r chi.Router) {
		r.Post("/", s.handlers.AddCard)
		r.Delete("/", s.handlers.DeleteCard)
		r.Put("/", s.handlers.UpdateCard)
	})

	s.server = &http.Server{
		Addr:    port,
		Handler: r,
	}
	log.Printf("started")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
