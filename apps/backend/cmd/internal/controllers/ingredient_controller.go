package controllers

import (
	"net/http"
	"strconv"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/gofiber/fiber/v2"
)

type IngredientController struct {
	IngredientService service.IngredientService
}

func NewIngredientController(ingredientService service.IngredientService) *IngredientController {
	return &IngredientController{
		IngredientService: ingredientService,
	}
}

func (controller *IngredientController) Create(c *fiber.Ctx) error {
	var request dto.IngredientRequest
	photo, err := c.FormFile("photo")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo is required",
		})
	}

	request.Name = c.FormValue("name")
	request.Photo = photo

	response, err := controller.IngredientService.Save(c.Context(), &request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func (controller *IngredientController) FindById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	response, err := controller.IngredientService.FindById(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response)
	}

	return c.Status(http.StatusOK).JSON(response)
}
