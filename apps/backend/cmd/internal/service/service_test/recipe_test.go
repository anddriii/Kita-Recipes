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

func TestRecipeSave(t *testing.T) {
	db := SetupTestDB()
	// defer cleanupTestDB(db)

	ctx := context.TODO()
	validate := validator.New()
	recipe := repository.NewRecipeRepository()
	recipePhotoRepository := repository.NewRecipePhotoRepository()
	categoryRepo := repository.NewCategoryRepository()
	recipeService := service.NewRecipeService(recipe, recipePhotoRepository, categoryRepo, db, validate)

	// Pastikan direktori storage/images ada sebelum test
	absolutePath, err := filepath.Abs("storage/images/categories/")
	require.NoError(t, err, "Gagal mendapatkan path absolut")
	err = os.MkdirAll(absolutePath, os.ModePerm)
	require.NoError(t, err, "Gagal membuat direktori storage/images")

	// Buat *multipart.FileHeader dari buffer
	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("photo.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Pastikan AuthorRequest menggunakan pointer
	request := &dto.RecipeRequestCreate{
		Name:           "Cara bikin Es teh bpy Furiahayaaa 16",
		Thumbnail:      fileHeader,
		About:          "Makaron adalah",
		UrlFile:        "makaron.pdf",
		UrlVideo:       "ada",
		CategoryId:     5,
		RecipeAuthorId: 8,
		RecipePhotos: []dto.RecipePhotos{
			{
				Photo: *fileHeader,
			},
			{
				Photo: *fileHeader,
			},
		},
	}

	// Jalankan service Save
	response, err := recipeService.Save(ctx, *request)
	require.NoError(t, err)
	assert.Equal(t, "Cara bikin Es teh bpy Furiahayaaa 15", response.Name)
}

func TestFindByIdRecipe(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewRecipeRepository()
	recipePhotoRepository := repository.NewRecipePhotoRepository()
	categoryRepo := repository.NewCategoryRepository()
	service := service.NewRecipeService(repo, recipePhotoRepository, categoryRepo, db, validate)

	ctx := context.Background()
	author, err := service.FindById(ctx, 77)

	assert.NoError(t, err)
	fmt.Println(author)
}

func TestUpdateRecipe(t *testing.T) {
	// Setup
	db := SetupTestDB()
	validate := validator.New()

	repo := repository.NewRecipeRepository()
	recipePhotoRepository := repository.NewRecipePhotoRepository()
	categoryRepo := repository.NewCategoryRepository()
	service := service.NewRecipeService(repo, recipePhotoRepository, categoryRepo, db, validate)

	var initialRecipe domain.Recipe
	result := db.Where("id = ?", 128).First(&initialRecipe)
	assert.Nil(t, result.Error)

	// // Seed initial data
	// initialRecipe := domain.Recipe{
	// 	Name:           "ini adalah data lama 106",
	// 	Slug:           "ini adalah data lama 106",
	// 	Thumbnail:      "original_photoasd.jpg",
	// 	About:          "ini tes lama",
	// 	CategoryId:     2,
	// 	RecipeAuthorId: 1,
	// 	UrlFile:        "data_lama.pdf",
	// 	UrlVideo:       "aidaiskd",
	// }
	// result := db.Create(&initialRecipe)
	// assert.Nil(t, result.Error)

	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("thumbnail.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Prepare update request
	updateRequest := &dto.RecipeRequestUpdate{
		ID:             initialRecipe.ID,
		Name:           "Cara bikin Jus mangga by Kafka",
		Slug:           initialRecipe.Slug,
		Thumbnail:      fileHeader,
		About:          "Semoga tidak error, pleasseeee",
		UrlFile:        "Update file",
		UrlVideo:       "update video",
		CategoryId:     3,
		RecipeAuthorId: 2,
	}

	// Call Update service
	ctx := context.Background()
	_, err = service.Update(ctx, *updateRequest)
	assert.Nil(t, err, "Service Update gagal")

	// Validate updated data in database
	// var updatedRecipe domain.Recipe
	// result = db.Where("id = ?", initialRecipe.ID).First(&updatedRecipe)
	// assert.Nil(t, result.Error)

	// Validate assertions
	// expectedSlug := utils.Slugify(updateRequest.Name)
	// assert.Equal(t, "Data setelah Update 42", updatedRecipe.Name)
	// assert.Equal(t, expectedSlug, updatedRecipe.Slug)
	// assert.Equal(t, "Update file", updatedRecipe.UrlFile)
	// assert.Equal(t, "update video", updatedRecipe.UrlVideo)
	// assert.Equal(t, updateRequest.CategoryId, updatedRecipe.CategoryId)
	// assert.Equal(t, updateRequest.RecipeAuthorId, updatedRecipe.RecipeAuthorId)

	// // Optional: Validate thumbnail
	// assert.Contains(t, updatedRecipe.Thumbnail, "thumbnail.jpg")

	// Debug logs
	// fmt.Printf("Updated Recipe: %+v\n", updatedRecipe)
}

func TestDeleteRecipe(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewRecipeRepository()
	recipePhotoRepository := repository.NewRecipePhotoRepository()
	categoryRepo := repository.NewCategoryRepository()
	service := service.NewRecipeService(repo, recipePhotoRepository, categoryRepo, db, validate)

	ctx := context.Background()
	service.Delete(ctx, 33)

	var author domain.Recipe
	result := db.First(&author, "id = ?", 33)
	fmt.Println(result)
}

func TestFindAllRecipes(t *testing.T) {
	db := SetupTestDB()

	validate := validator.New()
	repo := repository.NewRecipeRepository()
	recipePhotoRepository := repository.NewRecipePhotoRepository()
	categoryRepo := repository.NewCategoryRepository()
	service := service.NewRecipeService(repo, recipePhotoRepository, categoryRepo, db, validate)

	ctx := context.Background()
	recipes, err := service.FindAll(ctx)

	assert.NoError(t, err)
	fmt.Println(recipes)
}
