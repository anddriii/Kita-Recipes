package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func CategoryRouter(app *fiber.App, categoryController *controllers.CategoryController) {
	categoryGroup := app.Group("categories")
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Contenr-Type", "multipart/form-data")
		return c.Next()
	})

	categoryGroup.Post("/", categoryController.Create)
	categoryGroup.Put("/:id", categoryController.Update)
}
