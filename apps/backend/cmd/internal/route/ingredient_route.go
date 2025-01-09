package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func IngredientRouter(app *fiber.App, ingredientController *controllers.IngredientController) {
	ingredinetGroup := app.Group("ingredient")
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "multipart/form-data")
		return c.Next()
	})

	//create Ingredient
	ingredinetGroup.Post("/", ingredientController.Create)

	//update ingredient
	ingredinetGroup.Put("/:id", ingredientController.Update)

	//find By Id
	ingredinetGroup.Get("/:id", ingredientController.FindById)
}
