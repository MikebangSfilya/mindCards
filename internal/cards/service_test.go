package cards

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	myRedis "github.com/MikebangSfilya/mindCards/internal/repository/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

func (s *Service) noCache(ctx context.Context, cardID, userID int) (*MindCard, error) {
	row, err := s.Repo.GetCardById(ctx, nil, cardID, userID)
	if err != nil {
		return nil, err
	}
	return rowToCard(row), nil
}

func BenchmarkService_FullCycle(b *testing.B) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Step 1: Infrastructure setup
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("mindcards"),
		postgres.WithUsername("user"),
		postgres.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		b.Fatalf("failed to start postgres: %v", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			b.Logf("failed to terminate postgres container: %v", err)
		}
	}()

	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		b.Fatalf("failed to start redis: %v", err)
	}
	defer func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			b.Logf("failed to terminate redis container: %v", err)
		}
	}()

	pgConnStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	redisHost, _ := redisContainer.Host(ctx)
	redisPort, _ := redisContainer.MappedPort(ctx, "6379")

	// Step 2: Dependency initialization
	pool, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		b.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	setupFullSchema(b, pool)

	redisRepo, err := myRedis.New(redisHost, redisPort.Port(), "", 0)
	if err != nil {
		b.Fatalf("failed to connect to redis: %v", err)
	}

	repo := NewCardPool(pool)
	txM := NewTxManager(pool)
	service := NewService(repo, txM, logger, redisRepo)

	userID := 1
	_, err = pool.Exec(ctx, "INSERT INTO users (user_id, email, password_hash) VALUES ($1, 'bench@test.com', 'hash')", userID)
	if err != nil {
		b.Fatalf("failed to create user: %v", err)
	}

	b.ResetTimer()

	// Step 3: Run benchmarks

	// Test A: Add cards (Transaction check)
	b.Run("AddCards_Single_Transaction", func(b *testing.B) {
		params := []Card{{Title: "Bench Title", Description: "Desc", Tag: "benchmark"}}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.AddCards(ctx, userID, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Test B: Get list (Cache Hit check)
	b.Run("GetCards_List_CacheHit", func(b *testing.B) {
		seedCardsForList(b, service, userID, 10)
		_, _ = service.GetCards(ctx, userID, 10, 0)
		time.Sleep(100 * time.Millisecond)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.GetCards(ctx, userID, 10, 0)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Test C: Update (Cache Invalidation check)
	b.Run("Update_Description_And_Invalidate", func(b *testing.B) {
		cardID := 9999
		_, err := pool.Exec(ctx, `
          INSERT INTO memory_cards (card_id, user_id, title, card_description, tag, created_at, level_study, learned) 
          VALUES ($1, $2, 'UpdTitle', 'OldDesc', 'tag', NOW(), 0, false)
          ON CONFLICT (card_id) DO NOTHING`, cardID, userID)
		if err != nil {
			b.Fatal(err)
		}

		upd := Update{NewDescription: "New Bench Desc"}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.UpdateCardDescription(ctx, cardID, userID, upd)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// Test D: Battle (Postgres vs Redis)
	battleCardID := 5555
	_, err = pool.Exec(ctx, `
       INSERT INTO memory_cards (card_id, user_id, title, card_description, tag, created_at, level_study, learned) 
       VALUES ($1, $2, 'BattleTitle', 'BattleDesc', 'vs', NOW(), 0, false)
       ON CONFLICT (card_id) DO NOTHING`, battleCardID, userID)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Battle_Postgres_Direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := service.noCache(ctx, battleCardID, userID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	_, _ = service.GetCardById(ctx, battleCardID, userID)
	time.Sleep(100 * time.Millisecond)

	b.Run("Battle_Redis_Hit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := service.GetCardById(ctx, battleCardID, userID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Helpers
func setupFullSchema(b *testing.B, pool *pgxpool.Pool) {
	sql := `
    CREATE TABLE IF NOT EXISTS users (
       user_id SERIAL PRIMARY KEY, 
       email VARCHAR(100) NOT NULL UNIQUE, 
       password_hash VARCHAR(250) NOT NULL
    );
    CREATE TABLE IF NOT EXISTS memory_cards (
       card_id SERIAL PRIMARY KEY,
       user_id INT NOT NULL,
       title VARCHAR(255) NOT NULL,
       card_description TEXT,
       tag VARCHAR(100),
       created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
       level_study SMALLINT DEFAULT 0 CHECK (level_study >= 0 AND level_study <= 5),
       learned BOOLEAN DEFAULT FALSE,
       FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
    );
    CREATE INDEX IF NOT EXISTS idx_title ON memory_cards (title);
    CREATE INDEX IF NOT EXISTS idx_memory_cards_user_id ON memory_cards (user_id);`
	_, err := pool.Exec(context.Background(), sql)
	if err != nil {
		b.Fatalf("failed to setup schema: %v", err)
	}
}

func seedCardsForList(b *testing.B, s *Service, userID, count int) {
	params := make([]Card, count)
	for i := 0; i < count; i++ {
		params[i] = Card{Title: fmt.Sprintf("Title %d", i), Description: "Description", Tag: "seed"}
	}
	_, err := s.AddCards(context.Background(), userID, params)
	if err != nil {
		b.Fatalf("seeding failed: %v", err)
	}
}
