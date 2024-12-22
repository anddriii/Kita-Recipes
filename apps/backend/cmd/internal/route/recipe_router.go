package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func RecipeRouter(app *fiber.App, recipeController *controllers.RecipeController) {
	recipeGroup := app.Group("/recipes")
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "multipart/form-data")
		return c.Next()
	})

	recipeGroup.Post("/", recipeController.Create)
	recipeGroup.Put("/:id", recipeController.Update)
	recipeGroup.Get("/:id", recipeController.FindById)
	recipeGroup.Get("/", recipeController.FindAll)
	recipeGroup.Delete("/:id", recipeController.Delete)
}
