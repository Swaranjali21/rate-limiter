package middleware

import (
	"rate-limiter/config"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func NewRateLimiter() (r1 *RateLimiter) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.RedisAddr,
		Password: config.RedisConfig.RedisPassword,
		DB:       config.RedisConfig.RedisDB,
	})
	r1 = &RateLimiter{
		client: client,
		limit:  config.RedisConfig.RateLimit,
		window: config.RedisConfig.WindowSize,
	}
	return
}

func getClientID(c *fiber.Ctx) (value string) {
	apiKey := c.Get("X-API-Key")
	if apiKey != "" {
		value = "api:" + apiKey
		return
	}
	ipAddress := c.IP()
	if strings.Contains(ipAddress, ":") {
		ipAddress = strings.Split(ipAddress, ":")[0]
	}
	value = "ip:" + ipAddress
	return
}