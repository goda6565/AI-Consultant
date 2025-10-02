package redis

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/redis/go-redis/v9"
)

func ProvideRedisClient(ctx context.Context, e *environment.Environment) (*redis.Client, func()) {
	opt, err := redis.ParseURL(e.RedisURL)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	return client, func() {
		if err := client.Close(); err != nil {
			panic(err)
		}
	}
}
