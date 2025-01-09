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

func (controller *IngredientController) Update(c *fiber.Ctx) error {
	photo, _ := c.FormFile("photo")

	var request dto.IngredientRequest
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	request.Photo = photo

	request.ID = id
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response, err := controller.IngredientService.Update(c.Context(), &request)
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

func (controller *IngredientController) FindAll(c *fiber.Ctx) error {
	response, err := controller.IngredientService.FindAll(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (controller *IngredientController) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	controller.IngredientService.Delete(c.Context(), id)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Ingredient deleted succesfully",
	})
}
