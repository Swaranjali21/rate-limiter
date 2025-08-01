package main

import (
	"log"
	config "rate-limiter/config"
	middleware "rate-limiter/middleware"
	router "rate-limiter/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)

	}
	rateLimiter := middleware.NewRateLimiter()
	defer rateLimiter.Close()

	app := fiber.New()
	router.Configure(app, rateLimiter)

	addr := config.RedisConfig.ServerPort
	log.Printf("Fiber listening on %s (limit %d/%v)", addr, config.RedisConfig.RateLimit, config.RedisConfig.WindowSize)
	err = app.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}

}