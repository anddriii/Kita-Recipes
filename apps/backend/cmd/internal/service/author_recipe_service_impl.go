package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthorServiceImpl struct {
	AuthorRepository repository.AuthorRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewAuthorService(authorRepository repository.AuthorRepository, DB *gorm.DB, validate *validator.Validate) AuthorService {
	return &AuthorServiceImpl{
		AuthorRepository: authorRepository,
		DB:               DB,
		Validate:         validate,
	}
}

// Create implements AuthorService.
func (a *AuthorServiceImpl) Save(ctx context.Context, request *dto.AuthorRequest) (dto.AuthorResponses, error) {
	//validasi
	err := a.Validate.Struct(request)
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	authorPhoto := domain.RecipeAuthor{
		Name: request.Name,
	}
	filename := uuid.New().String() + "-" + authorPhoto.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return dto.AuthorResponses{}, err
	}

	file, err := request.Photo.Open()
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	filePath := "storage/images/" + filename
	fmt.Printf("Creating file: %s\n", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.AuthorResponses{}, err
	}
	defer out.Close()

	//salin foto ke penyimpanan
	_, err = io.Copy(out, file)
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	author := &domain.RecipeAuthor{
		Name:  request.Name,
		Photo: filename,
	}

	if _, err := a.AuthorRepository.Save(ctx, a.DB, author); err != nil {
		log.Printf("Error saving author:")
		return dto.AuthorResponses{}, err
	}
	log.Printf("Saved author with ID: %d", author.ID)

	return dto.AuthorResponses{
		ID:    author.ID,
		Name:  author.Name,
		Photo: author.Photo,
	}, nil

}

// Update implements AuthorService.
func (a *AuthorServiceImpl) Update(ctx context.Context, request *dto.AuthorRequest) (dto.AuthorResponseDetail, error) {
	err := a.Validate.Struct(request)
	if err != nil {
		return dto.AuthorResponseDetail{}, err
	}

	author, err := a.AuthorRepository.FindById(ctx, a.DB, request.ID)
	if err != nil {
		return dto.AuthorResponseDetail{}, err
	}

	author.Name = request.Name

	if request.Photo != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + author.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]

		// Buka file yang di-upload
		file, err := request.Photo.Open()
		if err != nil {
			log.Printf("Gagal open poto")
			panic(err)
		}
		defer file.Close()

		// Simpan file ke direktori
		filePath := "storage/images/" + filename
		fmt.Printf("Creating file: %s\n", filePath)
		out, err := os.Create(filePath)
		if err != nil {
			return dto.AuthorResponseDetail{}, err
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			return dto.AuthorResponseDetail{}, err
		}

		// Assign path foto baru ke author
		author.Photo = filename

		// err = os.Remove("storage/images" + author.Photo)
		// if err != nil {
		// 	return dto.AuthorResponseDetail{}, err
		// }

	}

	if _, err := a.AuthorRepository.Update(ctx, a.DB, &author); err != nil {
		log.Printf("Error saving author:")
		return dto.AuthorResponseDetail{}, err
	}

	return dto.AuthorResponseDetail{
		ID:      author.ID,
		Name:    author.Name,
		Photo:   author.Photo,
		Recipes: []*dto.RecipeResponses{},
	}, nil
}

// Delete implements AuthorService.
func (a *AuthorServiceImpl) Delete(ctx context.Context, authorID int) {
	author, err := a.AuthorRepository.FindById(ctx, a.DB, authorID)
	if err != nil {
		panic("Data tidak ditemukan")
	}
	a.AuthorRepository.Delete(ctx, a.DB, &author)
}

// FindAll implements AuthorService.
func (a *AuthorServiceImpl) FindAll(ctx context.Context) ([]dto.AuthorResponses, error) {
	authors, err := a.AuthorRepository.FindAll(ctx, a.DB)
	if err != nil {
		panic("Data tidak ditemukan")
	}

	var authorResponses []dto.AuthorResponses
	for _, author := range authors {
		authorResponses = append(authorResponses, dto.AuthorResponses{
			ID:    author.ID,
			Name:  author.Name,
			Photo: author.Photo,
		})
	}

	return authorResponses, nil
}

// FindById implements AuthorService.
func (a *AuthorServiceImpl) FindById(ctx context.Context, id int) (dto.AuthorResponseDetail, error) {
	authorDetail, err := a.AuthorRepository.FindById(ctx, a.DB, id)
	if err != nil {
		return dto.AuthorResponseDetail{}, err
	}

	response := dto.AuthorResponseDetail{
		ID:      authorDetail.ID,
		Name:    authorDetail.Name,
		Photo:   authorDetail.Photo,
		Recipes: []*dto.RecipeResponses{},
	}
	return response, nil
}
