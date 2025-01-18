package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/repository"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db := SetupTestDB()
	ctx := context.TODO()
	validate := validator.New()

	authRepo := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepo, db, validate)

	fileContent := []byte("dummy image content")
	fileheader, err := createTestFileHeaderFromBuffer("photo.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	request := dto.RegisterDTO{
		Name:     "Elaina",
		Photo:    fileheader,
		Username: "elaina",
		Password: "rahasia",
		Email:    "elaina21@gmail.com",
		Role:     "admin",
	}

	err = authService.Register(ctx, request)
	if err != nil {
		return
	}
}

func TestLogin(t *testing.T) {
	db := SetupTestDB()
	ctx := context.TODO()
	validate := validator.New()

	authRepo := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepo, db, validate)

	request := dto.LoginDTO{
		Username: "elaina",
		Password: "rahasia",
	}

	response, err := authService.Login(ctx, request)
	assert.Nil(t, err, "Gagal login")
	fmt.Println(response)
}
