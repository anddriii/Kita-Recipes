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

	"github.com/anddriii/KitaRecipes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	DB                 *gorm.DB
	Validate           *validator.Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository, DB *gorm.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

// Create implements CategoryService.
func (service *CategoryServiceImpl) Create(ctx context.Context, request dto.CategoryRequest) (dto.CategoryResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	categoryIcon := domain.Categories{
		Name: request.Name,
	}

	filename := uuid.New().String() + "-" + categoryIcon.Name + "." + strings.Split(request.Icon.Filename, ".")[len(strings.Split(request.Icon.Filename, "."))-1]
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return dto.CategoryResponse{}, err
	}

	file, err := request.Icon.Open()
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	filePath := "storage/images/categories/" + filename
	fmt.Printf("Creating file: %s\n", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	defer out.Close()

	//salin foto ke penyimpanan
	_, err = io.Copy(out, file)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	request.Slug = utils.Slugify(request.Name)

	category := &domain.Categories{
		Name: request.Name,
		Slug: request.Slug,
		Icon: filename,
	}
	if _, err := service.CategoryRepository.Save(ctx, service.DB, category); err != nil {
		log.Printf("Error saving category:")
		return dto.CategoryResponse{}, err
	}

	return dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
		Icon: category.Icon,
	}, nil

}

// Update implements CategoryService.
func (service *CategoryServiceImpl) Update(ctx context.Context, request dto.CategoryRequest) (dto.CategoryResponseDetail, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.CategoryResponseDetail{}, err
	}

	category, err := service.CategoryRepository.FindById(ctx, service.DB, int(request.ID))
	if err != nil {
		return dto.CategoryResponseDetail{}, err
	}
	newSlug := utils.Slugify(request.Name)
	request.Slug = newSlug
	category.Slug = newSlug
	category.Name = request.Name

	if request.Icon != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + category.Name + "." +
			strings.Split(request.Icon.Filename, ".")[len(strings.Split(request.Icon.Filename, "."))-1]

		// Buka file yang di-upload
		file, err := request.Icon.Open()
		if err != nil {
			return dto.CategoryResponseDetail{}, err
		}
		defer file.Close()

		// Simpan file ke direktori
		filePath := "storage/images/" + filename
		out, err := os.Create(filePath)
		if err != nil {
			return dto.CategoryResponseDetail{}, err
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			return dto.CategoryResponseDetail{}, err
		}

		// Assign path foto baru ke author
		category.Icon = filename

		// err = os.Remove("storage/images" + category.Icon)
		// if err != nil {
		// 	return dto.CategoryResponseDetail{}, err
		// }
	}

	if _, err := service.CategoryRepository.Update(ctx, service.DB, &category); err != nil {
		log.Printf("Error saving category")
		return dto.CategoryResponseDetail{}, err
	}

	return dto.CategoryResponseDetail{
		ID:     category.ID,
		Name:   category.Name,
		Slug:   category.Slug,
		Recipe: []*dto.RecipeResponses{},
	}, nil

}

// Delete implements CategoryService.
func (service *CategoryServiceImpl) Delete(ctx context.Context, id int) {
	category, err := service.CategoryRepository.FindById(ctx, service.DB, id)
	if err != nil {
		panic("Data tidak ditemukan")
	}
	service.CategoryRepository.Delete(ctx, service.DB, &category)
}

// FindAll implements CategoryService.
func (service *CategoryServiceImpl) FindAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := service.CategoryRepository.FindAll(ctx, service.DB)
	if err != nil {
		panic("Data tidak ditemukan")
	}

	var categoryResponses []dto.CategoryResponse
	for _, category := range categories {
		categoryResponses = append(categoryResponses, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
			Icon: category.Icon,
		})
	}
	return categoryResponses, nil
}

// FindById implements CategoryService.
func (service *CategoryServiceImpl) FindById(ctx context.Context, id int) (dto.CategoryResponseDetail, error) {
	categoryDetail, err := service.CategoryRepository.FindById(ctx, service.DB, int(id))
	if err != nil {
		return dto.CategoryResponseDetail{}, err
	}

	response := dto.CategoryResponseDetail{
		ID:     categoryDetail.ID,
		Name:   categoryDetail.Name,
		Slug:   categoryDetail.Slug,
		Icon:   categoryDetail.Icon,
		Recipe: []*dto.RecipeResponses{},
	}

	return response, err
}
