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

	//get all category
	categoryGroup.Get("/", categoryController.FindAll)

	//create category
	categoryGroup.Post("/", categoryController.Create)

	// Update category
	categoryGroup.Put("/:id", categoryController.Update)

	//Find by Id category
	categoryGroup.Get("/:id", categoryController.FindById)

	//Delete category
	categoryGroup.Delete("/:id", categoryController.Delete)

}
