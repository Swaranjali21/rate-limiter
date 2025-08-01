package middleware

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func (rl *RateLimiter) CheckLimit(clientID string) (result RateLimitResult, err error) {
	ctx := context.Background()
	currentTime := time.Now()
	windowStart := currentTime.Add(-rl.window)
	key := fmt.Sprintf("rate_limit:%s", clientID)

	pipe := rl.client.Pipeline()
	cutoff := strconv.FormatInt(windowStart.UnixNano(), 10)

	pipe.ZRemRangeByScore(ctx, key, "0", cutoff)
	countCmd := pipe.ZCard(ctx, key)
	oldestCmd := pipe.ZRangeWithScores(ctx, key, 0, 0)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return
	}

	currentWindowCount := countCmd.Val()
	oldRequests := oldestCmd.Val()

	if currentWindowCount >= int64(rl.limit) {
		var resetTime time.Time
		if len(oldRequests) > 0 {
			firstRequestScore := int64(oldRequests[0].Score)
			oldestTimestamp := time.Unix(0, firstRequestScore)
			resetTime = oldestTimestamp.Add(rl.window)
		} else {
			resetTime = currentTime.Add(rl.window)
		}
		result = RateLimitResult{
			Allowed:   false,
			Remaining: 0,
			ResetTime: resetTime,
		}
		return
	}

	// Add new request
	newPipe := rl.client.Pipeline()
	newPipe.ZAdd(ctx, key, &redis.Z{
		Score:  float64(currentTime.UnixNano()),
		Member: fmt.Sprintf("req_%d", currentTime.UnixNano()),
	})
	newPipe.Expire(ctx, key, rl.window+time.Minute)
	_, err = newPipe.Exec(ctx)
	if err != nil {
		return
	}

	result = RateLimitResult{
		Allowed:   true,
		Remaining: rl.limit - int(currentWindowCount) - 1,
		ResetTime: currentTime.Add(rl.window),
	}
	return
}

// Fiber middleware
func (rl *RateLimiter) FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		// get the clientID
		clientID := getClientID(c)
		// check are they allowed to make request
		result, err := rl.CheckLimit(clientID)
		if err != nil {
			log.Printf("Rate limit error: %v", err)
			err = c.Status(500).JSON(fiber.Map{"error": "internal error"})
			return
		}
		// inform about status

		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))
		// block when limit is exceeded

		if !result.Allowed {
			retryAfter := int64(time.Until(result.ResetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 1
			}
			c.Set("Retry-After", strconv.FormatInt(retryAfter, 10))
			err = c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many requests",
				"retry_after": retryAfter,
			})
			return
		}
		// allow the request to continue
		return c.Next()
	}
}

func (rl *RateLimiter) Close() error {
	return rl.client.Close()
}
