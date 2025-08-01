package config

import "time"

// RedisConfig holds environment-based settings
var (
	RedisConfig struct {
		RedisAddr     string
		RedisPassword string
		RedisDB       int
		ServerPort    string
		RateLimit     int
		WindowSize    time.Duration
	}
)
