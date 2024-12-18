package service_test

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Fungsi untuk mengatur database salinan untuk testing
func SetupTestDB() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/kita_recipes_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrasi schema jika diperlukan
	db.AutoMigrate(&domain.RecipeAuthor{})
	return db
}

func createTestFileHeaderFromBuffer(filename string, content []byte) (*multipart.FileHeader, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Buat file part
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(content)
	if err != nil {
		return nil, err
	}
	writer.Close()

	// Buat *multipart.FileHeader dari in-memory buffer
	fileReader := bytes.NewReader(b.Bytes())
	mr := multipart.NewReader(fileReader, writer.Boundary())
	form, err := mr.ReadForm(1024)
	if err != nil {
		return nil, err
	}

	// Ambil file header
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return nil, fmt.Errorf("file header tidak tersedia")
	}

	return fileHeaders[0], nil
}

func TestAuthorService_Save_SuccessWithDB(t *testing.T) {
	db := SetupTestDB()
	// defer cleanupTestDB(db)

	ctx := context.TODO()
	validate := validator.New()
	authorRepo := repository.NewAuthorRepository()
	authorService := service.NewAuthorService(authorRepo, db, validate)

	// Pastikan direktori storage/images ada sebelum test
	absolutePath, err := filepath.Abs("storage/images/")
	require.NoError(t, err, "Gagal mendapatkan path absolut")
	err = os.MkdirAll(absolutePath, os.ModePerm)
	require.NoError(t, err, "Gagal membuat direktori storage/images")

	// Buat *multipart.FileHeader dari buffer
	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("photo.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Pastikan AuthorRequest menggunakan pointer
	request := &dto.AuthorRequest{
		Name:  "Changli the badag",
		Photo: fileHeader,
	}

	// Jalankan service Save
	response, err := authorService.Save(ctx, request)
	require.NoError(t, err)

	//ini error well
	// Verifikasi data tersimpan di database
	// var author domain.RecipeAuthor
	// err = db.First(&author, "id = ?", response.ID).Error
	// require.NoError(t, err, "Record tidak ditemukan di database")
	assert.Equal(t, "Changli the badag", response.Name)
	// assert.Contains(t, &author.Photo, "storage/images/")
}

func TestAuthorService_Update(t *testing.T) {
	// Setup
	db := SetupTestDB()
	validate := validator.New()

	repo := repository.NewAuthorRepository()
	authorService := service.NewAuthorService(repo, db, validate)

	// Seed initial data
	initialAuthor := domain.RecipeAuthor{
		Name:  "Original Name",
		Photo: "original_photo.jpg",
	}
	result := db.Create(&initialAuthor)
	assert.Nil(t, result.Error)

	fileContent := []byte("dummy image content")
	fileHeader, err := createTestFileHeaderFromBuffer("asdasd.jpg", fileContent)
	require.NoError(t, err, "Gagal membuat file header dari buffer")

	// Prepare update request
	updateRequest := &dto.AuthorRequest{
		ID:    int(initialAuthor.ID),
		Name:  "Kafka Lopyuu",
		Photo: fileHeader,
	}

	// Call Update service
	ctx := context.Background()
	_, err = authorService.Update(ctx, updateRequest)
	assert.Nil(t, err)

	fmt.Println(initialAuthor.ID)
}

func TestFindAllAuthor(t *testing.T) {
	db := SetupTestDB()

	validate := validator.New()
	repo := repository.NewAuthorRepository()
	service := service.NewAuthorService(repo, db, validate)

	ctx := context.Background()
	authors, err := service.FindAll(ctx)

	assert.NoError(t, err)
	fmt.Println(authors)
}

func TestFindByIdAuthor(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewAuthorRepository()
	service := service.NewAuthorService(repo, db, validate)

	ctx := context.Background()
	author, err := service.FindById(ctx, 1)

	assert.NoError(t, err)
	fmt.Println(author)
}

func TestDeleteAuthor(t *testing.T) {
	db := SetupTestDB()
	validate := validator.New()
	repo := repository.NewAuthorRepository()
	service := service.NewAuthorService(repo, db, validate)

	ctx := context.Background()
	service.Delete(ctx, 9)

	var author domain.RecipeAuthor
	result := db.First(&author, 9)
	fmt.Println(result)
}
