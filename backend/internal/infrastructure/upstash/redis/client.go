package redis

import (
	"context"

	"github.com/goda6565/ai-consultant/backend/internal/infrastructure/environment"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, e *environment.Environment) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: e.RedisURL,
	})
}
