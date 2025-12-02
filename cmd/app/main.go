package main

import (
	"cards/internal/cards"
	"cards/internal/config"
	database "cards/internal/db"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	// the application is listening for the SIGTERM signal to exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Printf(".env not found: %v", err)
	}
	cfg := config.New()
	db := database.CreateDataBase(cfg)
	if db == nil {
		log.Fatal("Database connection failed")
		return
	}
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	repo := cards.NewPool(db)
	service := cards.NewService(repo, logger)
	handle := cards.New(service)

	handle.RegistredRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {

		log.Println(" Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil {
			slog.Warn(
				"Server start failed or shutdown",
				"server error", err)
		}
	}()

	<-ctx.Done()

	shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")

	if err := srv.Shutdown(shutdown); err != nil {
		log.Print("shutdown fail")
	}
	log.Print("Shutdown end")

}
