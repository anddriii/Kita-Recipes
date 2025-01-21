package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	middleware "github.com/anddriii/KitaRecipes/cmd/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authController *controllers.AuthController) {
	auth := app.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	// Semua route di bawah ini butuh login (JWT)
	auth.Get("/profile", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success", "message": "Authenticated"})
	})
}
