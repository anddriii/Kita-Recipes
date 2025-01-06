package controllers

import (
	"net/http"

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
