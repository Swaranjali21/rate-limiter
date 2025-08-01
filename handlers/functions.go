package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) (err error) {
	return c.Status(200).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   "Fiber API is running.",
	})
}

func Users(c *fiber.Ctx) (err error) {
	users := []fiber.Map{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	return c.JSON(fiber.Map{"users": users, "total": len(users)})
}

func Data(c *fiber.Ctx) (err error) {
	if c.Method() == "POST" {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		return c.Status(201).JSON(fiber.Map{
			"message": "Data created",
			"data":    data,
		})
	}
	return c.JSON(fiber.Map{
		"data":  []string{"item1", "item2"},
		"count": 2,
	})
}
