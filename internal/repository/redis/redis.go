package redis

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func MustLoad(host, port, password string, db int) *Redis {
	const op = "repository.redis.MustLoad"

	addr := net.JoinHostPort(host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to redis",
			slog.String("op", op),
			slog.String("addr", addr),
			slog.Any("error", err),
		)
		_ = client.Close()
		panic(fmt.Sprintf("%s: %v", op, err))
	}
	slog.Info("redis connected successfully", slog.String("addr", addr))

	return &Redis{Client: client}
}
