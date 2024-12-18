package service

import (
	"context"
	"errors"
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

type IngredientServiceImpl struct {
	IngredientRepository repository.IngredientRepository
	DB                   *gorm.DB
	Validate             *validator.Validate
}

func NewIngredientService(ingredientRepository repository.IngredientRepository, DB *gorm.DB, validate *validator.Validate) IngredientService {
	return &IngredientServiceImpl{
		IngredientRepository: ingredientRepository,
		DB:                   DB,
		Validate:             validate,
	}
}

// Save implements IngredientService.
func (service *IngredientServiceImpl) Save(ctx context.Context, request *dto.IngredientRequest) (dto.IngredientResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return dto.IngredientResponse{}, err
	}

	ingredientPhoto := domain.Ingredient{
		Name: request.Name,
	}
	filename := uuid.New().String() + "-" + ingredientPhoto.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return dto.IngredientResponse{}, err
	}

	file, err := request.Photo.Open()
	if err != nil {
		return dto.IngredientResponse{}, err
	}

	filePath := "storage/images/ingredients/" + filename
	if err := os.MkdirAll(filepath.Dir("storage/images/ingredients/"), os.ModePerm); err != nil {
		return dto.IngredientResponse{}, err
	}
	fmt.Printf("Creating file: %s\n", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.IngredientResponse{}, err
	}
	defer out.Close()

	//salin foto ke penyimpanan
	_, err = io.Copy(out, file)
	if err != nil {
		return dto.IngredientResponse{}, err
	}
	if request.Photo == nil {
		return dto.IngredientResponse{}, errors.New("PHOTO IS REQUIRED")
	}

	ingredient := &domain.Ingredient{
		Name:  request.Name,
		Photo: filename,
	}

	if _, err := service.IngredientRepository.Save(ctx, service.DB, ingredient); err != nil {
		log.Printf("Error saving ingredient")
		return dto.IngredientResponse{}, err
	}
	log.Printf("Saved ingredient with ID: %d", ingredient.ID)

	return dto.IngredientResponse{
		ID:    ingredient.ID,
		Name:  ingredient.Name,
		Photo: ingredient.Photo,
	}, nil

}

// Update implements IngredientService.
func (service *IngredientServiceImpl) Update(ctx context.Context, request *dto.IngredientRequest) (dto.IngredientResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.IngredientResponse{}, err
	}

	log.Printf("Retrieving ingredient with ID: %d", request.ID)
	ingredient, err := service.IngredientRepository.FindById(ctx, service.DB, request.ID)
	if err != nil {
		log.Printf("Failed to find ingredient: %v", err)
		return dto.IngredientResponse{}, err
	}
	log.Printf("Updating ingredient: %+v", ingredient)
	ingredient.Name = request.Name

	if request.Photo != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + ingredient.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]

		// Buka file yang di-upload
		file, err := request.Photo.Open()
		if err != nil {
			log.Printf("Gagal open poto")
			panic(err)
		}
		defer file.Close()

		// Simpan file ke direktori
		filePath := "storage/images/ingredients/" + filename
		fmt.Printf("Creating file: %s\n", filePath)
		out, err := os.Create(filePath)
		if err != nil {
			return dto.IngredientResponse{}, err
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			return dto.IngredientResponse{}, err
		}

		// Assign path foto baru ke author
		ingredient.Photo = filename

		// err = os.Remove("storage/images" + author.Photo)
		// if err != nil {
		// 	return dto.AuthorResponseDetail{}, err
		// }
	}
	if _, err := service.IngredientRepository.Update(ctx, service.DB, &ingredient); err != nil {
		log.Printf("Error update Ingredient")
		return dto.IngredientResponse{}, err
	}

	return dto.IngredientResponse{
		ID:    ingredient.ID,
		Name:  ingredient.Name,
		Photo: ingredient.Photo,
	}, nil
}

// Delete implements IngredientService.
func (service *IngredientServiceImpl) Delete(ctx context.Context, ingredientId int) {
	ingredient, err := service.IngredientRepository.FindById(ctx, service.DB, ingredientId)
	if err != nil {
		log.Println("Error Deleting Ingredient")
	}
	service.IngredientRepository.Delete(ctx, service.DB, ingredient)
}

// FindById implements IngredientService.
func (service *IngredientServiceImpl) FindById(ctx context.Context, ingredientId int) (dto.IngredientResponse, error) {
	ingredient, err := service.IngredientRepository.FindById(ctx, service.DB, ingredientId)
	if err != nil {
		return dto.IngredientResponse{}, err
	}

	response := dto.IngredientResponse{
		ID:    ingredient.ID,
		Name:  ingredient.Name,
		Photo: ingredient.Photo,
	}

	return response, nil
}

// FindAll implements IngredientService.
func (service *IngredientServiceImpl) FindAll(ctx context.Context) ([]dto.IngredientResponse, error) {
	ingredients, err := service.IngredientRepository.FindAll(ctx, service.DB)
	if err != nil {
		return []dto.IngredientResponse{}, err
	}

	var ingredientResponses []dto.IngredientResponse
	for _, ingredient := range ingredients {
		ingredientResponses = append(ingredientResponses, dto.IngredientResponse{
			ID:    ingredient.ID,
			Name:  ingredient.Name,
			Photo: ingredient.Photo,
		})
	}

	return ingredientResponses, nil
}
