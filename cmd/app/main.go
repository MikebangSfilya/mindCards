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
	"github.com/MikebangSfilya/mindCards/internal/repository/db"
	redis2 "github.com/MikebangSfilya/mindCards/internal/repository/redis"
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

	red := initRedis(cfg)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	router := chi.NewRouter()

	applyMiddleware(router)

	cardsRepo := cards.NewCardPool(db)
	userRepo := users.NewUserPool(db)

	txMan := cards.NewTxManager(db)
	cardsService := cards.NewService(cardsRepo, txMan, logger, red)
	cardsHandler := cards.New(cardsService)

	//registrated handlers
	cardsHandler.RegistredRoutes(router)
	router.Post("/user", users.SaveUser(userRepo))

	srv := newServer(cfg, router)

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

	db.Close()
	red.Client.Close()
	log.Print("Shutdown end")

}

func newServer(cfg config.Config, router chi.Router) *http.Server {
	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTTPServer.Timeout,
		WriteTimeout: cfg.HTTTPServer.Timeout,
		IdleTimeout:  cfg.HTTTPServer.IdleTimeout,
	}
}

func applyMiddleware(r chi.Router) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}

func initRedis(cfg config.Config) *redis2.Redis {
	rd := redis2.MustLoad(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	return rd
}
