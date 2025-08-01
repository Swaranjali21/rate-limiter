package config

import (
	"log"
	"os"
	"strconv"
	"time"

	godotenv "github.com/joho/godotenv"
)

func Init() (err error) {
	err = godotenv.Load("env/.config")
	if err != nil {
		log.Println("Warning: .env file not found or not loaded")
	}
	limit := getEnv("RATE_LIMIT", "5")
	rateLimit, err := strconv.Atoi(limit)
	if err != nil {
		return
	}
	window := getEnv("WINDOW_SECONDS", "30")
	windowSeconds, err := strconv.Atoi(window)
	if err != nil {
		return
	}
	configRedisDB := getEnv("REDIS_DB", "0")
	redisDB, err := strconv.Atoi(configRedisDB)
	if err != nil {
		return
	}
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	serverPort := getEnv("SERVER_PORT", ":8080")
	RedisConfig.RedisAddr = redisAddr
	RedisConfig.RedisPassword = redisPassword
	RedisConfig.RedisDB = redisDB
	RedisConfig.ServerPort = serverPort
	RedisConfig.RateLimit = rateLimit
	RedisConfig.WindowSize = time.Duration(windowSeconds) * time.Second
	return
}

func getEnv(key string, defaultVal string) (result string) {
	result = os.Getenv(key)
	if result != "" {
		return
	}
	result = defaultVal
	return
}
