package main

import (
	"github.com/anddriii/KitaRecipes/cmd/internal/controllers"
	"github.com/anddriii/KitaRecipes/cmd/internal/model"
	"github.com/anddriii/KitaRecipes/cmd/internal/repository"
	router "github.com/anddriii/KitaRecipes/cmd/internal/route"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Database & Dependency Injection
	db := model.OpenConnection() // Fungsi untuk koneksi database
	validator := validator.New()
	recipePhotoRepository := repository.NewRecipePhotoRepository()
	recipeRepo := repository.NewRecipeRepository()
	categoryRepo := repository.NewCategoryRepository()
	authorRepo := repository.NewAuthorRepository()
	tutorialRepo := repository.NewRecipeTutorialsRepository()
	recipeService := service.NewRecipeService(recipeRepo, recipePhotoRepository, categoryRepo, authorRepo, tutorialRepo, db, validator)
	recipeController := controllers.NewRecipeController(recipeService)
	categoryService := service.NewCategoryService(categoryRepo, db, validator)
	categoryController := controllers.NewCategoryController(categoryService)
	authorService := service.NewAuthorService(authorRepo, db, validator)
	authorController := controllers.NewAuthorController(authorService)

	// Routes
	router.RecipeRouter(app, recipeController)
	router.CategoryRouter(app, categoryController)
	router.AuthorRouter(app, authorController)

	// Start Server
	app.Listen(":3000")
}
