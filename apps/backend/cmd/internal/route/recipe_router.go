package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	middleware "github.com/anddriii/KitaRecipes/cmd/internal/middlewares"
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
	recipeGroup.Get("/:id", middleware.APIKEYMiddleware, recipeController.FindById)

	//delete recipe
	recipeGroup.Delete("/:id", middleware.APIKEYMiddleware, recipeController.Delete)

	//download file recipes
	recipeGroup.Get("/file/:filename", recipeController.DownloadFile)
}
