package redis

import (
	"context"
	"fmt"

	"github.com/K1ver/e-commerce-api/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.Db,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  cfg.Redis.PoolTimeout,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
