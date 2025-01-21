package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func APIKEYMiddleware(c *fiber.Ctx) error {
	apikey := c.Get("X-API-KEY")
	validApiKey := os.Getenv("API_KEY")

	if validApiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "API key not set in server",
		})
	}

	if apikey != validApiKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid API key",
		})
	}
	return c.Next()
}
