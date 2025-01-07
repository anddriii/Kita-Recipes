package router

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthorRouter(app *fiber.App, authorController *controllers.AuthorController) {
	authorGroup := app.Group("author")
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "multipart/form-data")
		return c.Next()
	})

	// Create Author
	authorGroup.Post("/", authorController.Create)

	// Update Author
	authorGroup.Put("/:id", authorController.Update)

	//Find by Id
	authorGroup.Get("/:id", authorController.FindById)
}
