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
	CategoryRepository    repository.CategoryRepository
	RecipePhotoRepository repository.RecipePhotoRepository
	AuthorRepository      repository.AuthorRepository
	TutorialsRepository   repository.RecipeTutorials
	DB                    *gorm.DB
	Validate              *validator.Validate
}

func NewRecipeService(recipeRepository repository.RecipeRespository, recipePhoto repository.RecipePhotoRepository, categoryRepository repository.CategoryRepository, authorRepository repository.AuthorRepository, tutorialsRepo repository.RecipeTutorials, DB *gorm.DB, validate *validator.Validate) RecipeService {
	return &RecipeServiceImpl{
		RecipeRepository:      recipeRepository,
		CategoryRepository:    categoryRepository,
		AuthorRepository:      authorRepository,
		RecipePhotoRepository: recipePhoto,
		TutorialsRepository:   tutorialsRepo,
		DB:                    DB,
		Validate:              validate,
	}
}

// Save implements RecipeService.
func (service *RecipeServiceImpl) Save(ctx context.Context, request dto.RecipeRequestCreate) (dto.RecipeResponseCreate, error) {
	// Validasi request
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.RecipeResponseCreate{}, err
	}

	recipeName := domain.Recipe{
		Name: request.Name,
	}

	// Path direktori utama
	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return dto.RecipeResponseCreate{}, fmt.Errorf("gagal memendapatkan path absolut: %w", err)
	}

	// Proses Thumbnail
	thumbnailFilename := uuid.New().String() + "-" + recipeName.Name + "." +
		strings.Split(request.Thumbnail.Filename, ".")[len(strings.Split(request.Thumbnail.Filename, "."))-1]
	thumbnailPath := filepath.Join(basePath, "images", "recipes", "thumbnails", thumbnailFilename)

	if err := saveFile(request.Thumbnail, thumbnailPath); err != nil {
		return dto.RecipeResponseCreate{}, fmt.Errorf("failed to save thumbnail: %w", err)
	}

	// Proses RecipePhotos
	var photoFilenames []string
	for _, photo := range request.RecipePhotos {
		photoFilename := uuid.New().String() + "-" + recipeName.Name + "." +
			strings.Split(photo.Photo.Filename, ".")[len(strings.Split(photo.Photo.Filename, "."))-1]
		photoPath := filepath.Join(basePath, "images", "recipes", "recipes_photos", photoFilename)

		log.Println("Saving photo to:", photoPath) // Debug log

		if err := saveFile(&photo.Photo, photoPath); err != nil {
			return dto.RecipeResponseCreate{}, fmt.Errorf("failed to save photo: %w", err)
		}
		photoFilenames = append(photoFilenames, photoFilename)
	}

	// Simpan data Recipe ke DB
	recipe := domain.Recipe{
		Name:           request.Name,
		Slug:           utils.Slugify(request.Name),
		Thumbnail:      thumbnailFilename,
		About:          request.About,
		UrlFile:        request.UrlFile,
		UrlVideo:       request.UrlVideo,
		CategoryId:     int64(request.CategoryId),
		RecipeAuthorId: int64(request.RecipeAuthorId),
	}

	// tutorial := request.Tutorials
	// recipe.RecipeTutorial = append(recipe.RecipeTutorial, domain.RecipeTutorial{
	// 	Name: tutorial.Name,
	// })

	// simopan tutorial ke DB
	for _, t := range request.Tutorials {
		recipe.RecipeTutorial = append(recipe.RecipeTutorial, domain.RecipeTutorial{
			Name: t.Name,
		})
	}

	if _, err := service.RecipeRepository.Save(ctx, service.DB, &recipe); err != nil {
		return dto.RecipeResponseCreate{}, err
	}

	// Simpan RecipePhotos ke DB
	for _, filename := range photoFilenames {
		photo := domain.Photo{
			RecipeId: recipe.ID,
			Photo:    filename,
		}
		if err := service.RecipePhotoRepository.Save(ctx, service.DB, photo); err != nil {
			return dto.RecipeResponseCreate{}, err
		}
	}

	// menampilkan recipes photos dari DB
	photoses, err := service.RecipePhotoRepository.Show(ctx, service.DB, int(recipe.ID))
	if err != nil {
		return dto.RecipeResponseCreate{}, err
	}

	var photos []*dto.RecipePhotosResponse
	for _, photo := range photoses {
		photos = append(photos, &dto.RecipePhotosResponse{
			ID:   photo.ID,
			Name: photo.Photo,
		})
	}

	// Mapping Author
	author := &dto.AuthorResponses{
		ID:    recipe.RecipeAuthorId,
		Name:  recipe.RecipeAuthor.Name,
		Photo: recipe.RecipeAuthor.Photo,
	}

	// Mapping Category
	category, err := service.CategoryRepository.FindById(ctx, service.DB, request.CategoryId)
	if err != nil {
		return dto.RecipeResponseCreate{}, err
	}

	categoryResponse := dto.ToCategoryResponse(category)

	response := dto.RecipeResponseCreate{
		ID:                   recipe.ID,
		Name:                 recipe.Name,
		Slug:                 recipe.Slug,
		Thumbnail:            recipe.Thumbnail,
		CategoryResponse:     &categoryResponse,
		About:                recipe.About,
		UrlFile:              recipe.UrlFile,
		UrlVideo:             recipe.UrlVideo,
		RecipePhotosResponse: photos,
		Author:               author,
	}
	// Log response in JSON format
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	log.Println(string(responseJSON))

	return response, nil
}

// Update implements RecipeService.
func (service *RecipeServiceImpl) Update(ctx context.Context, request dto.RecipeRequestUpdate) (dto.RecipeResponseUpdate, error) {
	err := service.Validate.Struct(request)
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
	}

	recipe, err := service.RecipeRepository.FindById(ctx, service.DB, int(request.ID))
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
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

	// mapping tutorial
	for _, t := range request.Tutorials {
		recipe.RecipeTutorial = append(recipe.RecipeTutorial, domain.RecipeTutorial{
			Name: t.Name,
		})
	}
	// Path direktori utama
	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return dto.RecipeResponseUpdate{}, fmt.Errorf("gagal memendapatkan path absolut: %w", err)
	}

	//update foto
	if request.Thumbnail != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + recipe.Name + "." +
			strings.Split(request.Thumbnail.Filename, ".")[len(strings.Split(request.Thumbnail.Filename, "."))-1]

		// menyimpan foto update
		thumbnailPath := filepath.Join(basePath, "images", "recipes", "thumbnails", filename)

		if err := saveFile(request.Thumbnail, thumbnailPath); err != nil {
			return dto.RecipeResponseUpdate{}, fmt.Errorf("failed to save thumbnail: %w", err)
		}

		// Assign path foto baru ke author
		recipe.Thumbnail = filename

		// err = os.Remove("storage/images" + category.Icon)
		// if err != nil {
		// 	return dto.CategoryResponseDetail{}, err
		// }
	}

	if request.RecipePhotos != nil {
		// Proses RecipePhotos
		var photoFilenames []string
		for _, photo := range request.RecipePhotos {
			photoFilename := uuid.New().String() + "-" + recipe.Name + "." +
				strings.Split(photo.Photo.Filename, ".")[len(strings.Split(photo.Photo.Filename, "."))-1]
			photoPath := filepath.Join(basePath, "images", "recipes", "recipes_photos", photoFilename)

			log.Println("Saving photo to:", photoPath) // Debug log

			if err := saveFile(&photo.Photo, photoPath); err != nil {

				return dto.RecipeResponseUpdate{}, fmt.Errorf("failed to save photo: %w", err)
			}
			photoFilenames = append(photoFilenames, photoFilename)
		}

		// Simpan foto ke database sekali saja
		log.Println("Attempting to save photos in bulk")
		if err := service.RecipePhotoRepository.Update(ctx, service.DB, int(request.ID), photoFilenames); err != nil {
			log.Printf("Failed to save photos: %v", err)
			return dto.RecipeResponseUpdate{}, err
		}
		log.Println("All photos saved successfully")
	}

	log.Println("recipe category id request", request.CategoryId)

	recipeDetail, err := service.RecipeRepository.Update(ctx, service.DB, &recipe)
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
	}

	// Mapping Category
	category, err := service.CategoryRepository.FindById(ctx, service.DB, request.CategoryId)
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
	}
	// category domain to category response
	categoryResponse := dto.ToCategoryResponse(category)

	//mappin photos
	photoses, err := service.RecipePhotoRepository.Show(ctx, service.DB, int(recipe.ID))
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
	}

	var photos []*dto.RecipePhotosResponse
	for _, photo := range photoses {
		photos = append(photos, &dto.RecipePhotosResponse{
			ID:   photo.ID,
			Name: photo.Photo,
		})
	}

	// Mapping Author
	author, err := service.AuthorRepository.FindById(ctx, service.DB, request.RecipeAuthorId)
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
	}
	// author domain to author response
	authorResponse := dto.ToAuthorResponse(author)

	// Mapping Ingredients
	var ingredients []*dto.IngredientResponse
	for _, ingredient := range recipeDetail.Ingredients {
		ingredients = append(ingredients, &dto.IngredientResponse{
			ID:    ingredient.ID,
			Name:  ingredient.Name,
			Photo: ingredient.Photo,
		})
	}

	//mapping tutorials
	tutorialsDb, err := service.TutorialsRepository.Show(ctx, service.DB, int(recipe.ID))
	if err != nil {
		return dto.RecipeResponseUpdate{}, err
	}

	// tutorialResponse := dto.ToTutorialResponse(tutoriald)

	var tutorialsRes []*dto.Tutorials
	for _, tutorial := range tutorialsDb {
		tutorialsRes = append(tutorialsRes, &dto.Tutorials{
			ID:   tutorial.ID,
			Name: tutorial.Name,
		})
	}

	// Build the response
	response := dto.RecipeResponseUpdate{
		ID:               recipeDetail.ID,
		Name:             recipeDetail.Name,
		Slug:             recipeDetail.Slug,
		UrlFile:          recipeDetail.UrlFile,
		UrlVideo:         recipeDetail.UrlVideo,
		Thumbnail:        recipeDetail.Thumbnail,
		About:            recipeDetail.About,
		RecipeTutorials:  tutorialsRes,
		Ingredients:      ingredients,
		CategoryResponse: &categoryResponse,
		Author:           &authorResponse,
		RecipePhotos:     photos,
	}

	log.Printf("Recipe to save: %+v", recipe)
	log.Printf("CategoryId di Service: %d", request.CategoryId)

	// Log response in JSON format
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	log.Println(string(responseJSON))

	return response, nil
}

func saveFile(fileHeader *multipart.FileHeader, filePath string) error {
	dir := filepath.Dir(filePath)

	// Pastikan direktori ada
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	log.Println("File saved successfully:", filePath)
	return nil
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
