package controllers

import (
	"net/http"
	"strconv"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/gofiber/fiber/v2"
)

type RecipeController struct {
	RecipeService service.RecipeService
}

func NewRecipeController(recipeService service.RecipeService) *RecipeController {
	return &RecipeController{
		RecipeService: recipeService,
	}
}

// Create Recipe
func (controller *RecipeController) Create(c *fiber.Ctx) error {
	var request dto.RecipeRequestCreate

	// Parsing fields from multipart form
	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Thumbnail is required",
		})
	}

	urlVideo := c.FormValue("url_video")
	if urlVideo == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "UrlVideo is required",
		})
	}

	// Parse other string fields
	request.Name = c.FormValue("name")
	request.Slug = c.FormValue("slug")
	request.About = c.FormValue("about")
	request.UrlFile = c.FormValue("url_file")
	request.UrlVideo = urlVideo

	// Parse integer fields
	categoryId, _ := strconv.Atoi(c.FormValue("category_id"))
	request.CategoryId = categoryId

	recipeAuthorId, _ := strconv.Atoi(c.FormValue("recipe_author_id"))
	request.RecipeAuthorId = recipeAuthorId

	// Assign Thumbnail
	request.Thumbnail = thumbnail

	// Parse "photos" (array of files)
	form, err := c.MultipartForm() // Inisialisasi variabel form
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}

	photos := form.File["photos"] // Get array of photos
	if len(photos) > 0 {
		var photoUploads []dto.RecipePhotos
		for _, photo := range photos {
			photoUploads = append(photoUploads, dto.RecipePhotos{
				Photo: *photo,
			})
		}
		request.RecipePhotos = photoUploads
	}

	// Call Service Layer
	response, err := controller.RecipeService.Save(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}

// Update Recipe
func (controller *RecipeController) Update(c *fiber.Ctx) error {
	var request dto.RecipeRequestUpdate
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	request.ID = int64(id)
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response, err := controller.RecipeService.Update(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

// Get Recipe By ID
func (controller *RecipeController) FindById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	response, err := controller.RecipeService.FindById(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

// Get All Recipes
func (controller *RecipeController) FindAll(c *fiber.Ctx) error {
	response, err := controller.RecipeService.FindAll(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

// Delete Recipe
func (controller *RecipeController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	controller.RecipeService.Delete(c.Context(), id)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Recipe deleted successfully",
	})
}
