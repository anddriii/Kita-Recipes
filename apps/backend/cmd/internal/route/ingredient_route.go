package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func IngredientRouter(app *fiber.App, ingredientController *controllers.IngredientController) {
	ingredientGroup := app.Group("ingredient")
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "multipart/form-data")
		return c.Next()
	})

	//create Ingredient
	ingredientGroup.Post("/", ingredientController.Create)

	//update ingredient
	ingredientGroup.Put("/:id", ingredientController.Update)

	//find By Id
	ingredientGroup.Get("/:id", ingredientController.FindById)

	//get all
	ingredientGroup.Get("/", ingredientController.FindAll)

	//delete ingredient
	ingredientGroup.Delete("/:id", ingredientController.Delete)
}
