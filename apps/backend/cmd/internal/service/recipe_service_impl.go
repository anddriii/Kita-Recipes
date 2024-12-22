package service

import (
	"context"
	"encoding/json"
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

type RecipeServiceImpl struct {
	RecipeRepository repository.RecipeRespository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewRecipeService(recipeRepository repository.RecipeRespository, DB *gorm.DB, validate *validator.Validate) RecipeService {
	return &RecipeServiceImpl{
		RecipeRepository: recipeRepository,
		DB:               DB,
		Validate:         validate,
	}
}

// Save implements RecipeService.
func (service *RecipeServiceImpl) Save(ctx context.Context, request dto.RecipeRequestCreate) (dto.RecipeResponses, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.RecipeResponses{}, err
	}

	recipeThumbnailName := domain.Recipe{
		Name: request.Name,
	}

	if request.CategoryId == 0 {
		return dto.RecipeResponses{}, fmt.Errorf("INVALID CATEGORY ID: %d", request.CategoryId)
	}

	if request.RecipeAuthorId == 0 {
		return dto.RecipeResponses{}, fmt.Errorf("INVALID RECIPE AUTHOR ID: %d", request.RecipeAuthorId)
	}

	filename := uuid.New().String() + "-" + recipeThumbnailName.Name + "." + strings.Split(request.Thumbnail.Filename, ".")[len(strings.Split(request.Thumbnail.Filename, "."))-1]
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return dto.RecipeResponses{}, err
	}

	file, err := request.Thumbnail.Open()
	if err != nil {
		return dto.RecipeResponses{}, err
	}

	basePath, err := os.Getwd()
	if err != nil {
		return dto.RecipeResponses{}, err
	}

	uploadPath := filepath.Join(basePath, "assets", "images", "recipes")
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return dto.RecipeResponses{}, err
	}

	filePath := filepath.Join(uploadPath, filename)
	fmt.Printf("Creating file: %s\n", filePath)
	out, err := os.Create(filePath)
	if err != nil {
		return dto.RecipeResponses{}, err
	}
	defer out.Close()

	//salin foto ke penyimpanan
	_, err = io.Copy(out, file)
	if err != nil {
		return dto.RecipeResponses{}, err
	}

	request.Slug = utils.Slugify(request.Name)

	recipe := domain.Recipe{
		Name:           request.Name,
		Slug:           request.Slug,
		Thumbnail:      filename,
		About:          request.About,
		UrlFile:        request.UrlFile,
		UrlVideo:       request.UrlVideo,
		CategoryId:     int64(request.CategoryId),
		RecipeAuthorId: int64(request.RecipeAuthorId),
		RecipePhoto:    []domain.Photo{},
	}

	if _, err := service.RecipeRepository.Save(ctx, service.DB, &recipe); err != nil {
		log.Printf("Error saving category:")
		return dto.RecipeResponses{}, err
	}

	// 📝 Simpan RecipePhotos
	photoDir := "../../../../assets/images/recipe_photos"
	if err := os.MkdirAll(photoDir, os.ModePerm); err != nil {
		return dto.RecipeResponses{}, fmt.Errorf("failed to create recipe photos directory: %v", err)
	}

	for _, photo := range request.RecipePhotos {
		photoFilename := uuid.New().String() + "-" + photo.Filename
		photoPath := filepath.Join(photoDir, photoFilename)

		photoFile, err := photo.File.Open()
		if err != nil {
			return dto.RecipeResponses{}, fmt.Errorf("failed to open recipe photo file: %v", err)
		}
		defer photoFile.Close()

		outPhoto, err := os.Create(photoPath)
		if err != nil {
			return dto.RecipeResponses{}, fmt.Errorf("failed to create recipe photo file: %v", err)
		}
		defer outPhoto.Close()

		if _, err := io.Copy(outPhoto, photoFile); err != nil {
			return dto.RecipeResponses{}, fmt.Errorf("failed to save recipe photo file: %v", err)
		}

		recipePhoto := domain.Photo{
			RecipeId: recipe.ID,
			Photo:    photoFilename,
		}

		if err := service.DB.Create(&recipePhoto).Error; err != nil {
			return dto.RecipeResponses{}, fmt.Errorf("failed to save recipe photo: %v", err)
		}
	}

	return dto.RecipeResponses{
		ID:         recipe.ID,
		Name:       recipe.Name,
		Slug:       recipe.Slug,
		Thumbnail:  recipe.Thumbnail,
		About:      recipe.About,
		UrlFile:    recipe.UrlFile,
		UrlVideo:   recipe.UrlVideo,
		CategoryId: int(recipe.CategoryId),
	}, err
}

// Update implements RecipeService.
func (service *RecipeServiceImpl) Update(ctx context.Context, request dto.RecipeRequestUpdate) (dto.RecipeResponseDetail, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.RecipeResponseDetail{}, err
	}

	recipe, err := service.RecipeRepository.FindById(ctx, service.DB, int(request.ID))
	if err != nil {
		return dto.RecipeResponseDetail{}, err
	}

	newSlug := utils.Slugify(request.Name)
	request.Slug = newSlug
	recipe.Slug = newSlug
	recipe.Name = request.Name
	recipe.About = request.About
	recipe.UrlFile = request.UrlFile
	recipe.UrlVideo = request.UrlVideo
	recipe.CategoryId = int64(request.CategoryId)
	recipe.RecipeAuthorId = int64(request.RecipeAuthorId)

	//update foto
	if request.Thumbnail != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + recipe.Name + "." +
			strings.Split(request.Thumbnail.Filename, ".")[len(strings.Split(request.Thumbnail.Filename, "."))-1]

		// Buka file yang di-upload
		file, err := request.Thumbnail.Open()
		if err != nil {
			return dto.RecipeResponseDetail{}, err
		}
		defer file.Close()

		// Simpan file ke direktori
		filePath := "storage/images/recipes/" + filename
		out, err := os.Create(filePath)
		if err != nil {
			return dto.RecipeResponseDetail{}, err
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			return dto.RecipeResponseDetail{}, err
		}

		// Assign path foto baru ke author
		recipe.Thumbnail = filename

		// err = os.Remove("storage/images" + category.Icon)
		// if err != nil {
		// 	return dto.CategoryResponseDetail{}, err
		// }
	}

	recipeDetail, err := service.RecipeRepository.Update(ctx, service.DB, &recipe)
	if err != nil {
		return dto.RecipeResponseDetail{}, err
	}

	// Mapping Category
	category := &dto.CategoryResponse{
		ID:   recipeDetail.Category.ID,
		Name: recipeDetail.Category.Name,
		Slug: recipeDetail.Category.Slug,
		Icon: recipeDetail.Category.Icon,
	}

	// Mapping Author
	author := &dto.AuthorResponses{
		ID:    recipeDetail.RecipeAuthor.ID,
		Name:  recipeDetail.RecipeAuthor.Name,
		Photo: recipeDetail.RecipeAuthor.Photo,
	}

	// Mapping Ingredients
	var ingredients []*dto.IngredientResponse
	for _, ingredient := range recipeDetail.Ingredients {
		ingredients = append(ingredients, &dto.IngredientResponse{
			ID:    ingredient.ID,
			Name:  ingredient.Name,
			Photo: ingredient.Photo,
		})
	}

	var tutorials []*dto.Tutorials
	for _, tutorial := range recipeDetail.RecipeTutorial {
		tutorials = append(tutorials, &dto.Tutorials{
			ID:   tutorial.ID,
			Name: tutorial.Name,
		})
	}

	var photos []*dto.RecipePhotos
	for _, photo := range recipeDetail.RecipePhoto {
		photos = append(photos, &dto.RecipePhotos{
			ID:    photo.ID,
			Photo: photo.Photo,
		})
	}

	// Build the response
	response := dto.RecipeResponseDetail{
		ID:               recipeDetail.ID,
		Name:             recipeDetail.Name,
		Slug:             recipeDetail.Slug,
		UrlFile:          recipeDetail.UrlFile,
		UrlVideo:         recipeDetail.UrlVideo,
		Thumbnail:        recipeDetail.Thumbnail,
		About:            recipeDetail.About,
		RecipeTutorials:  tutorials,
		Ingredients:      ingredients,
		CategoryResponse: category,
		Author:           author,
		RecipePhotos:     photos,
	}

	log.Printf("Recipe to save: %+v", recipe)

	// Log response in JSON format
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	log.Println(string(responseJSON))

	return response, nil
}

// Delete implements RecipeService.
func (service *RecipeServiceImpl) Delete(ctx context.Context, recipeId int) {
	recipe, err := service.RecipeRepository.FindById(ctx, service.DB, recipeId)
	if err != nil {
		log.Printf("Data tidak ditemukan")
	}

	service.RecipeRepository.Delete(ctx, service.DB, &recipe)

}

// FindAll implements RecipeService.
func (service *RecipeServiceImpl) FindAll(ctx context.Context) ([]dto.RecipeResponses, error) {
	recipes, err := service.RecipeRepository.FindAll(ctx, service.DB)
	if err != nil {
		return []dto.RecipeResponses{}, err
	}

	var recipeResponses []dto.RecipeResponses

	for _, recipe := range recipes {
		category := &dto.CategoryResponse{
			ID:   recipe.Category.ID,
			Name: recipe.Category.Name,
			Slug: recipe.Category.Slug,
			Icon: recipe.Category.Icon,
		}
		// Ambil semua RecipePhoto dari elemen ini
		photos := []dto.PhotoUpload{}
		for _, photo := range recipe.RecipePhoto {
			photos = append(photos, dto.PhotoUpload{
				ID:       photo.ID,
				Filename: photo.Photo,
			})
		}

		// Buat objek RecipeResponses
		recipeResponses = append(recipeResponses, dto.RecipeResponses{
			ID:               recipe.ID,
			Name:             recipe.Name,
			Slug:             recipe.Slug,
			UrlFile:          recipe.UrlFile,
			UrlVideo:         recipe.UrlVideo,
			Thumbnail:        recipe.Thumbnail,
			About:            recipe.About,
			RecipePhotos:     photos,
			CategoryResponse: category,
		})
	}
	// Log response in JSON format
	responseJSON, _ := json.MarshalIndent(recipeResponses, "", "  ")
	log.Println(string(responseJSON))

	return recipeResponses, nil
}

func (service *RecipeServiceImpl) FindById(ctx context.Context, id int) (dto.RecipeResponseDetail, error) {
	recipeDetail, err := service.RecipeRepository.FindById(ctx, service.DB, id)
	if err != nil {
		return dto.RecipeResponseDetail{}, err
	}

	// Mapping Category
	category := &dto.CategoryResponse{
		ID:   recipeDetail.Category.ID,
		Name: recipeDetail.Category.Name,
		Slug: recipeDetail.Category.Slug,
		Icon: recipeDetail.Category.Icon,
	}

	// Mapping Author
	author := &dto.AuthorResponses{
		ID:    recipeDetail.RecipeAuthor.ID,
		Name:  recipeDetail.RecipeAuthor.Name,
		Photo: recipeDetail.RecipeAuthor.Photo,
	}

	// Mapping Ingredients
	var ingredients []*dto.IngredientResponse
	for _, ingredient := range recipeDetail.Ingredients {
		ingredients = append(ingredients, &dto.IngredientResponse{
			ID:    ingredient.ID,
			Name:  ingredient.Name,
			Photo: ingredient.Photo,
		})
	}

	var tutorials []*dto.Tutorials
	for _, tutorial := range recipeDetail.RecipeTutorial {
		tutorials = append(tutorials, &dto.Tutorials{
			ID:   tutorial.ID,
			Name: tutorial.Name,
		})
	}

	var photos []*dto.RecipePhotos
	for _, photo := range recipeDetail.RecipePhoto {
		photos = append(photos, &dto.RecipePhotos{
			ID:    photo.ID,
			Photo: photo.Photo,
		})
	}

	// Build the response
	response := dto.RecipeResponseDetail{
		ID:               recipeDetail.ID,
		Name:             recipeDetail.Name,
		Slug:             recipeDetail.Slug,
		UrlFile:          recipeDetail.UrlFile,
		UrlVideo:         recipeDetail.UrlVideo,
		Thumbnail:        recipeDetail.Thumbnail,
		About:            recipeDetail.About,
		RecipeTutorials:  tutorials,
		Ingredients:      ingredients,
		CategoryResponse: category,
		Author:           author,
		RecipePhotos:     photos,
	}

	// Log response in JSON format
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	log.Println(string(responseJSON))

	return response, nil
}

func (service *RecipeServiceImpl) HandleRecipePhotos(ctx context.Context, recipeID int, photos []dto.PhotoUpload) ([]domain.Photo, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	uploadPath := filepath.Join(basePath, "assets", "images", "recipes")
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	var savedPhotos []domain.Photo
	for _, photoUpload := range photos {
		file, err := photoUpload.File.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open photo file: %w", err)
		}
		defer file.Close()

		filename := uuid.New().String() + "-" + photoUpload.File.Filename
		filePath := filepath.Join(uploadPath, filename)

		out, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create photo file: %w", err)
		}
		defer out.Close()

		if _, err := io.Copy(out, file); err != nil {
			return nil, fmt.Errorf("failed to save photo file: %w", err)
		}

		savedPhotos = append(savedPhotos, domain.Photo{
			RecipeId: int64(recipeID),
			Photo:    filename,
		})
	}

	return savedPhotos, nil
}
