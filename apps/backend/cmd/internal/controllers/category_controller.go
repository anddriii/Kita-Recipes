package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	CategoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) *CategoryController {
	return &CategoryController{
		CategoryService: categoryService,
	}
}

func (controller *CategoryController) Create(c *fiber.Ctx) error {
	var request dto.CategoryRequest
	icon, err := c.FormFile("icon")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Icon is required",
		})
	}

	request.Name = c.FormValue("name")
	request.Slug = c.FormValue("slug")
	request.Icon = icon

	response, err := controller.CategoryService.Create(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func (controller *CategoryController) Update(c *fiber.Ctx) error {
	icon, err := c.FormFile("icon")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Icon is required",
		})
	}

	var request dto.CategoryRequest
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	request.Icon = icon

	request.ID = int64(id)
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error body parser: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response, err := controller.CategoryService.Update(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}
