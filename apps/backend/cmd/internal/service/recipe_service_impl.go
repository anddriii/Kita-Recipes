package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
	RecipeRepository      repository.RecipeRespository
	RecipePhotoRepository repository.RecipePhotoRepository
	DB                    *gorm.DB
	Validate              *validator.Validate
}

func NewRecipeService(recipeRepository repository.RecipeRespository, recipePhoto repository.RecipePhotoRepository, DB *gorm.DB, validate *validator.Validate) RecipeService {
	return &RecipeServiceImpl{
		RecipeRepository:      recipeRepository,
		RecipePhotoRepository: recipePhoto,
		DB:                    DB,
		Validate:              validate,
	}
}

// Save implements RecipeService.
func (service *RecipeServiceImpl) Save(ctx context.Context, request dto.RecipeRequestCreate) (dto.RecipeResponses, error) {
	// Validasi request
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.RecipeResponses{}, err
	}

	recipeName := domain.Recipe{
		Name: request.Name,
	}

	// Proses Thumbnail
	thumbnailFilename := uuid.New().String() + "-" + recipeName.Name + "." +
		strings.Split(request.Thumbnail.Filename, ".")[len(strings.Split(request.Thumbnail.Filename, "."))-1]
	thumbnailPath := "../../../../../../assets/images/recipes/" + thumbnailFilename
	if err := saveFile(request.Thumbnail, thumbnailPath); err != nil {
		return dto.RecipeResponses{}, fmt.Errorf("failed to save thumbnail: %w", err)
	}

	// Proses RecipePhotos
	var photoFilenames []string
	for _, photo := range request.RecipePhotos {
		photoFilename := uuid.New().String() + "-" + recipeName.Name + "." +
			strings.Split(request.Thumbnail.Filename, ".")[len(strings.Split(request.Thumbnail.Filename, "."))-1]
		photoPath := "../../../../../../assets/images/recipes/recipes_photos/" + photoFilename
		// Buat direktori jika belum ada
		// if err := os.MkdirAll(photoPath, os.ModePerm); err != nil {
		// 	return dto.RecipeResponseDetail{}, fmt.Errorf("failed to create directory: %w", err)
		// }

		if err := saveFile(&photo.Photo, photoPath); err != nil {
			return dto.RecipeResponses{}, fmt.Errorf("failed to save photo: %w", err)
		}

		photoFilenames = append(photoFilenames, photoFilename)
	}

	// Simpan data Recipe ke DB
	recipe := domain.RecipeDetail{
		Name:           request.Name,
		Slug:           utils.Slugify(request.Name),
		Thumbnail:      thumbnailFilename,
		About:          request.About,
		UrlFile:        request.UrlFile,
		UrlVideo:       request.UrlVideo,
		CategoryId:     int64(request.CategoryId),
		RecipeAuthorId: int64(request.RecipeAuthorId),
	}

	if _, err := service.RecipeRepository.Save(ctx, service.DB, &recipe); err != nil {
		return dto.RecipeResponses{}, err
	}

	// Simpan RecipePhotos ke DB
	for _, filename := range photoFilenames {
		photo := domain.Photo{
			RecipeId: recipe.ID,
			Photo:    filename,
		}
		if err := service.RecipePhotoRepository.Save(ctx, service.DB, photo); err != nil {
			return dto.RecipeResponses{}, err
		}
	}

	var photos []*dto.RecipePhotos
	for _, photo := range recipe.RecipePhoto {
		photos = append(photos, &dto.RecipePhotos{
			ID: photo.ID,
		})
	}

	// Debugging dari postmant, untuk mengecek file yang diupload
	for i, photo := range request.RecipePhotos {
		log.Printf("Processing photo %d: %s\n", i, photo.Photo.Filename)
	}

	// Mapping Category
	category := &dto.CategoryResponse{
		ID:   recipe.Category.ID,
		Name: recipe.Category.Name,
		Slug: recipe.Category.Slug,
		Icon: recipe.Category.Icon,
	}

	response := dto.RecipeResponses{
		ID:               recipe.ID,
		Name:             recipe.Name,
		Slug:             recipe.Slug,
		Thumbnail:        recipe.Thumbnail,
		CategoryResponse: category,
		About:            recipe.About,
		UrlFile:          recipe.UrlFile,
		UrlVideo:         recipe.UrlVideo,
		RecipePhotos:     photos,
	}
	// Log response in JSON format
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	log.Println(string(responseJSON))

	return response, nil
}

func saveFile(fileHeader *multipart.FileHeader, filePath string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
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
			ID: photo.ID,
			// Photo: photo.Photo,
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
		photos := []*dto.RecipePhotos{}
		for _, photo := range recipe.RecipePhoto {
			photos = append(photos, &dto.RecipePhotos{
				ID: photo.ID,
				// Photo: photo.Photo,
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
			ID: photo.ID,
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
