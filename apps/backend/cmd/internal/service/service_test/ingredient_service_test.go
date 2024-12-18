package service_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/repository"
	"github.com/anddriii/KitaRecipes/cmd/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngredientSave(t *testing.T) {
	db := SetupTestDB()
	// defer cleanupTestDB(db)

	ctx := context.TODO()
	validate := validator.New()
	ingredientRepo := repository.NewIngredientRepository()
	InredientService := service.NewIngredientService(ingredientRepo, db, validate)

	// Pastikan direktori storage/images ada sebelum test
	absolutePath, err := filepath.Abs("storage/images/ingredients/")
	require.NoError(t, err, "Gagal mendapatkan path absolut")
	err = os.MkdirAll(absolutePath, os.ModePerm)
	require.NoError(t, err, "Gagal membuat direktori storage/images")

	// Buat *multipart.FileHeader dari buffer
	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("photo.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Pastikan AuthorRequest menggunakan pointer
	request := &dto.IngredientRequest{
		Name:  "Changli istri saya",
		Photo: fileHeader,
	}

	// Jalankan service Save
	response, err := InredientService.Save(ctx, request)
	require.NoError(t, err)
	assert.Equal(t, "Changli istri saya", response.Name)
}

func TestUpdateIngredient(t *testing.T) {
	// Setup
	db := SetupTestDB()
	validate := validator.New()

	repo := repository.NewIngredientRepository()
	service := service.NewIngredientService(repo, db, validate)

	// Seed initial data
	ingredient := domain.Ingredient{
		Name:  "ini tes",
		Photo: "foto lama",
	}
	result := db.Create(&ingredient)
	assert.Nil(t, result.Error)

	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("asdasd.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Prepare update request
	updateRequest := &dto.IngredientRequest{
		ID:    int(ingredient.ID),
		Name:  "Zhezhi tukang",
		Photo: fileHeader,
	}

	// Call Update service
	ctx := context.Background()
	_, err = service.Update(ctx, updateRequest)
	assert.Nil(t, err, "Zhezhi tukang gambar")

	fmt.Println(ingredient.ID)
	fmt.Println(ingredient)
}

func TestFindByIdIngredient(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewIngredientRepository()
	service := service.NewIngredientService(repo, db, validate)

	ctx := context.Background()
	author, err := service.FindById(ctx, 7)

	assert.NoError(t, err)
	fmt.Println(author)
}

func TestDeleteIngredient(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewIngredientRepository()
	service := service.NewIngredientService(repo, db, validate)

	ctx := context.Background()
	service.Delete(ctx, 15)

	var ingredient domain.Ingredient
	result := db.First(&ingredient, 15)
	fmt.Println(result)
}

func TestFindAllIngredient(t *testing.T) {
	db := SetupTestDB()

	validate := validator.New()
	repo := repository.NewIngredientRepository()
	service := service.NewIngredientService(repo, db, validate)

	ctx := context.Background()
	ingredients, err := service.FindAll(ctx)

	assert.NoError(t, err)
	fmt.Println(ingredients)
}
