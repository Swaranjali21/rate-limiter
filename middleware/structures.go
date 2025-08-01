package middleware

import (
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiter and Result structs
type RateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

type RateLimitResult struct {
	Allowed   bool
	Remaining int
	ResetTime time.Time
}