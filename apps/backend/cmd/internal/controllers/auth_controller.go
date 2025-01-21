package controllers

import (
	"net/http"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	var request dto.RegisterDTO
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	request.Name = c.FormValue("name")
	photo, _ := c.FormFile("photo")
	request.Photo = photo
	request.Username = c.FormValue("username")
	request.Password = c.FormValue("password")
	request.Email = c.FormValue("email")

	if controller.AuthService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal server error",
		})
	}

	err := controller.AuthService.Register(c.Context(), &request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Register success",
	})

}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	req.Username = c.FormValue("username")
	req.Password = c.FormValue("password")

	token, err := controller.AuthService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid username or password",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login success",
		"token":   token,
	})
}
