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

	//get all recipes
	recipeGroup.Get("/", recipeController.FindAll)

	//create recipe
	recipeGroup.Post("/", recipeController.Create)

	//update recipe
	recipeGroup.Put("/:id", recipeController.Update)

	//get recipe by id
	recipeGroup.Get("/:id", recipeController.FindById)

	//delete recipe
	recipeGroup.Delete("/:id", recipeController.Delete)

	//download file recipes
	recipeGroup.Get("/file/:filename", recipeController.DownloadFile)
}
