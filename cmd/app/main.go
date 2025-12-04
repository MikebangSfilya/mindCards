package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MikebangSfilya/mindCards/internal/cards"
	"github.com/MikebangSfilya/mindCards/internal/config"
	database "github.com/MikebangSfilya/mindCards/internal/db"
	"github.com/MikebangSfilya/mindCards/internal/users"

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

	cfg := config.MustLoad()

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
	router.Use(middleware.Recoverer)

	cardsRepo := cards.NewCardPool(db)
	userRepo := users.NewUserPool(db)

	cardsService := cards.NewService(cardsRepo, logger)
	cardsHandler := cards.New(cardsService)

	//registrated handlers
	cardsHandler.RegistredRoutes(router)
	router.Post("/user", users.SaveUser(userRepo))

	srv := &http.Server{
		Addr:         cfg.Adress,
		Handler:      router,
		ReadTimeout:  cfg.HTTTPServer.Timeout,
		WriteTimeout: cfg.HTTTPServer.Timeout,
		IdleTimeout:  cfg.HTTTPServer.IdleTimeout,
	}
	go func() {

		log.Println(" Server starting")
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
