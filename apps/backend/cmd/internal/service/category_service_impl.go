package service

import (
	"context"
	"fmt"
	"log"
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
	RecipeRepo         repository.RecipeRespository
	DB                 *gorm.DB
	Validate           *validator.Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository, recipeRepo repository.RecipeRespository, DB *gorm.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		RecipeRepo:         recipeRepo,
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
	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	categoryIcon := domain.Categories{
		Name: request.Name,
	}

	filename := uuid.New().String() + "-" + categoryIcon.Name + "." +
		strings.Split(request.Icon.Filename, ".")[len(strings.Split(request.Icon.Filename, "."))-1]
	iconPath := filepath.Join(basePath, "category_icon", filename)

	if err := utils.SaveFile(request.Icon, iconPath); err != nil {
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
func (service *CategoryServiceImpl) Update(ctx context.Context, request dto.CategoryRequest) (dto.CategoryResponse, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	category, err := service.CategoryRepository.FindById(ctx, service.DB, int(request.ID))
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	newSlug := utils.Slugify(request.Name)
	request.Slug = newSlug
	category.Slug = newSlug
	category.Name = request.Name

	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return dto.CategoryResponse{}, fmt.Errorf("gagal mendapatkan path: %w", err)
	}

	if request.Icon != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + category.Name + "." +
			strings.Split(request.Icon.Filename, ".")[len(strings.Split(request.Icon.Filename, "."))-1]

		iconPath := filepath.Join(basePath, "category_icon", filename)

		if err := utils.SaveFile(request.Icon, iconPath); err != nil {
			return dto.CategoryResponse{}, err
		}

		// Assign path foto baru ke author
		category.Icon = filename

		// err = os.Remove("storage/images" + category.Icon)
		// if err != nil {
		// 	return dto.CategoryResponseDetail{}, err
		// }
		//
	}

	if _, err := service.CategoryRepository.Update(ctx, service.DB, &category); err != nil {
		log.Printf("Error saving category")
		return dto.CategoryResponse{}, err
	}

	return dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
		Icon: category.Icon,
	}, nil

}

// Delete implements CategoryService.
func (service *CategoryServiceImpl) Delete(ctx context.Context, id int) {
	category, err := service.CategoryRepository.FindById(ctx, service.DB, id)
	if err != nil {
		return
	}
	service.CategoryRepository.Delete(ctx, service.DB, &category)
}

// FindAll implements CategoryService.
func (service *CategoryServiceImpl) FindAll(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := service.CategoryRepository.FindAll(ctx, service.DB)
	if err != nil {
		return []dto.CategoryResponse{}, err
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

	//mapping recipes
	recipeRepo, err := service.RecipeRepo.GetRecipeCategory(ctx, service.DB, int(categoryDetail.ID))
	if err != nil {
		return dto.CategoryResponseDetail{}, err
	}

	var recipes []*dto.RecipeResponseCategory
	for _, recipe := range recipeRepo {
		recipes = append(recipes, &dto.RecipeResponseCategory{
			ID:        recipe.ID,
			Name:      recipe.Name,
			Slug:      recipe.Slug,
			Thumbnail: recipe.Thumbnail,
			About:     recipe.About,
		})
	}

	response := dto.CategoryResponseDetail{
		ID:     categoryDetail.ID,
		Name:   categoryDetail.Name,
		Slug:   categoryDetail.Slug,
		Icon:   categoryDetail.Icon,
		Recipe: recipes,
	}

	return response, err
}
