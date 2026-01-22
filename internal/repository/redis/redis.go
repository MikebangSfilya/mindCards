package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

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
	slog.Info("Redis connected successfully", slog.String("addr", addr))

	return &Redis{Client: client}
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%v: failed to marshal value for key %s: %w", key, err)
	}
	return r.Client.Set(ctx, key, data, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string, dest any) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return json.Unmarshal(data, dest)
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
