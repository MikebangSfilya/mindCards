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
	GetCards(w http.ResponseWriter, r *http.Request)
	GetByTag(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
}
type Server struct {
	handlers handlers
	server   *http.Server
	router   *chi.Mux
}

func NewServer(handl handlers) *Server {
	return &Server{
		handlers: handl,
		router:   chi.NewRouter(),
	}
}

func (s *Server) Start() error {

	port := ":8080"

	s.router.Route("/card", func(r chi.Router) {
		r.Post("/", s.handlers.AddCard)          //add card
		r.Delete("/{id}", s.handlers.DeleteCard) // Delete card
		r.Put("/{id}", s.handlers.UpdateCard)    // Update card
		r.Get("/tag/{tag}", s.handlers.GetByTag) // Get by tag
		r.Get("/", s.handlers.GetCards)          // Get all card(Default limit = 50, offset = 0) #TODO queryParams, not json body
		r.Get("/{id}", s.handlers.GetById)       // get by unic ID
	})

	s.server = &http.Server{
		Addr:    port,
		Handler: s.router,
	}
	log.Printf("started")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
