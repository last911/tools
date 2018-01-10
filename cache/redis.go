package cache

import (
	"github.com/go-redis/redis"
)

// Options alias redis.Options
type Options = redis.Options

// Redis *redis.Client
type Redis struct {
	*redis.Client
}

// NewRedis return *Redis
// Options{Addr: host:port}
func NewRedis(opts *redis.Options) *Redis {
	return &Redis{redis.NewClient(opts)}
}
