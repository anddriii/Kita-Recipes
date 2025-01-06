package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthorController struct {
	AuthorService service.AuthorService
}

func NewAuthorController(authorService service.AuthorService) *AuthorController {
	return &AuthorController{
		AuthorService: authorService,
	}
}

func (controller *AuthorController) Create(c *fiber.Ctx) error {
	var request dto.AuthorRequest
	photo, err := c.FormFile("photo")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Photo is required",
		})
	}

	request.Name = c.FormValue("name")
	request.Photo = photo

	response, err := controller.AuthorService.Save(c.Context(), &request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func (controller *AuthorController) Update(c *fiber.Ctx) error {
	photo, err := c.FormFile("photo")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "photo is required",
		})
	}

	var request dto.AuthorRequest
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	request.Photo = photo

	request.ID = id

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error body parser: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response, err := controller.AuthorService.Update(c.Context(), &request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(response)
}
