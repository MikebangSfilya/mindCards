package main

import (
	"cards/database"
	"cards/internal/api/handlers"
	"cards/internal/api/server"
	"cards/internal/configurate"
	"cards/internal/repo"
	"cards/internal/service"
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	// the application is listening for the SIGTERM signal to exit
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Printf(".env not found: %v", err)
	}
	cfg := configurate.New()
	db := database.CreateDataBase(cfg)
	if db == nil {
		log.Fatal("Database connection failed")
		return
	}
	defer db.Close()

	repo := repo.New(db)
	service := service.New(repo)
	handle := handlers.New(service)

	srv := server.NewServer(handle)
	go func() {

		log.Println(" Server starting on :8080")
		if err := srv.Start(); err != nil {
			slog.Warn(
				"Server start failed or shutdown",
				"server error", err)
		}
	}()

	<-ctx.Done()

	shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")

	srv.Shutdown(shutdown)
	log.Print("Shutdown end")

}
