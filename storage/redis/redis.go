package redis

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/singhJasvinder101/go_blog/internal/config"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config, db int) (*RedisClient, error) {
	// println("redis url", cfg.RedisConn+ strconv.Itoa(db))

	opts, err := redis.ParseURL(cfg.RedisConn + strconv.Itoa(db))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	return &RedisClient{
		Client: client,
	}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, val interface{}, exp time.Duration) error {
	return r.Client.Set(ctx, key, val, exp).Err()
}
