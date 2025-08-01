package router

import (
	handlers "rate-limiter/handlers"
	middleware "rate-limiter/middleware"

	"github.com/gofiber/fiber/v2"
)

// registers all routes on the Fiber app, applying middleware as needed
func Configure(app *fiber.App, rl *middleware.RateLimiter) {
	// Register the rate limiter middleware on the entire app or per group
	app.Use(rl.FiberMiddleware())

	app.Get("/health", handlers.Health)
	app.Get("/api/users", handlers.Users)
	app.All("/api/data", handlers.Data)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the Fiber Rate Limited API!",
			"endpoints": []string{
				"GET /health",
				"GET /api/users",
				"GET,POST /api/data",
			},
		})
	})
}
