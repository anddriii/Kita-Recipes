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
	recipeRepo := repository.NewRecipeRepository()
	recipeService := service.NewRecipeService(recipeRepo, db, validator)
	recipeController := controllers.NewRecipeController(recipeService)

	// Routes
	router.RecipeRouter(app, recipeController)

	// Start Server
	app.Listen(":3000")
}
