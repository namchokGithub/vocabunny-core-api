package infrastructure

import (
	"context"
	"fmt"

	"github.com/namchokGithub/vocabunny-core-api/configs"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg configs.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

func PingRedis(ctx context.Context, client *redis.Client) error {
	if client == nil {
		return nil
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}

	return nil
}
