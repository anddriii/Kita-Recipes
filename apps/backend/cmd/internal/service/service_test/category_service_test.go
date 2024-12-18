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

func TestCategorySave(t *testing.T) {
	db := SetupTestDB()
	// defer cleanupTestDB(db)

	ctx := context.TODO()
	validate := validator.New()
	authorRepo := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(authorRepo, db, validate)

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
	request := &dto.CategoryRequest{
		Name: "Changli istri saya",
		Icon: fileHeader,
	}

	// Jalankan service Save
	response, err := categoryService.Create(ctx, *request)
	require.NoError(t, err)
	assert.Equal(t, "Changli istri saya", response.Name)
}

func TestUpdateCategory(t *testing.T) {
	// Setup
	db := SetupTestDB()
	validate := validator.New()

	repo := repository.NewCategoryRepository()
	service := service.NewCategoryService(repo, db, validate)

	// Seed initial data
	initialAuthor := domain.Categories{
		Name: "ini tes",
		Slug: "ini tes",
		Icon: "original_photoasd.jpg",
	}
	result := db.Create(&initialAuthor)
	assert.Nil(t, result.Error)

	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("asdasd.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Prepare update request
	updateRequest := &dto.CategoryRequest{
		ID:   initialAuthor.ID,
		Name: "Zhezhi tukang gambar",
		Slug: "ini tes",
		Icon: fileHeader,
	}

	// Call Update service
	ctx := context.Background()
	_, err = service.Update(ctx, updateRequest)
	assert.Nil(t, err, "apa saja ")

	fmt.Println(initialAuthor.ID)
}

func TestFindAllCategory(t *testing.T) {
	db := SetupTestDB()

	validate := validator.New()
	repo := repository.NewCategoryRepository()
	service := service.NewCategoryService(repo, db, validate)

	ctx := context.Background()
	authors, err := service.FindAll(ctx)

	assert.NoError(t, err)
	fmt.Println(authors)
}

func TestFindByIdCategory(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewCategoryRepository()
	service := service.NewCategoryService(repo, db, validate)

	ctx := context.Background()
	author, err := service.FindById(ctx, "ini tes")

	assert.NoError(t, err)
	fmt.Println(author)
}

func TestDeleteCategory(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewCategoryRepository()
	service := service.NewCategoryService(repo, db, validate)

	ctx := context.Background()
	service.Delete(ctx, "ams")

	var author domain.RecipeAuthor
	result := db.First(&author, "ams")
	fmt.Println(result)
}
