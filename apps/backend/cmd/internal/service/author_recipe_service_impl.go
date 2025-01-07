package service

import (
	"context"
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

type AuthorServiceImpl struct {
	AuthorRepository repository.AuthorRepository
	RecipeRepo       repository.RecipeRespository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewAuthorService(authorRepository repository.AuthorRepository, recipeRepo repository.RecipeRespository, DB *gorm.DB, validate *validator.Validate) AuthorService {
	return &AuthorServiceImpl{
		AuthorRepository: authorRepository,
		RecipeRepo:       recipeRepo,
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

	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	authorPhoto := domain.RecipeAuthor{
		Name: request.Name,
	}
	filename := uuid.New().String() + "-" + authorPhoto.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]
	photoPath := filepath.Join(basePath, "author_photo", filename)

	if err := utils.SaveFile(request.Photo, photoPath); err != nil {
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
func (a *AuthorServiceImpl) Update(ctx context.Context, request *dto.AuthorRequest) (dto.AuthorResponses, error) {
	err := a.Validate.Struct(request)
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	author, err := a.AuthorRepository.FindById(ctx, a.DB, request.ID)
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return dto.AuthorResponses{}, err
	}

	author.Name = request.Name

	if request.Photo != nil {
		// Generate nama file baru
		filename := uuid.New().String() + "-" + author.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]

		photoPath := filepath.Join(basePath, "author_photo", filename)

		if err := utils.SaveFile(request.Photo, photoPath); err != nil {
			return dto.AuthorResponses{}, err
		}

		// Assign path foto baru ke author
		author.Photo = filename

		// err = os.Remove("storage/images" + author.Photo)
		// if err != nil {
		// 	return dto.AuthorResponses{}, err
		// }

	}

	if _, err := a.AuthorRepository.Update(ctx, a.DB, &author); err != nil {
		log.Printf("Error saving author:")
		return dto.AuthorResponses{}, err
	}

	return dto.AuthorResponses{
		ID:    author.ID,
		Name:  author.Name,
		Photo: author.Photo,
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

	recipeRepo, err := a.RecipeRepo.RecipeAuthor(ctx, a.DB, int(authorDetail.ID))
	if err != nil {
		return dto.AuthorResponseDetail{}, err
	}

	var recipes []*dto.RecipeResponseDetail
	for _, recipe := range recipeRepo {
		recipes = append(recipes, &dto.RecipeResponseDetail{
			ID:        recipe.ID,
			Name:      recipe.Name,
			Thumbnail: recipe.Thumbnail,
			About:     recipe.About,
		})
	}

	response := dto.AuthorResponseDetail{
		ID:      authorDetail.ID,
		Name:    authorDetail.Name,
		Photo:   authorDetail.Photo,
		Recipes: recipes,
	}
	return response, nil
}
