package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

	"gorm.io/gorm"
)

type AuthorRepositoryImpl struct {
}

func NewAuthorRepository() AuthorRepository {
	return &AuthorRepositoryImpl{}
}

func (repository *AuthorRepositoryImpl) Save(ctx context.Context, db *gorm.DB, author *domain.RecipeAuthor) (domain.RecipeAuthor, error) {
	err := db.WithContext(ctx).Create(author).Error
	if err != nil {
		return domain.RecipeAuthor{}, err
	}
	return *author, nil
}

func (repository *AuthorRepositoryImpl) Update(ctx context.Context, db *gorm.DB, author *domain.RecipeAuthorDetail) (domain.RecipeAuthorDetail, error) {
	err := db.WithContext(ctx).Model(&domain.RecipeAuthor{}).Where("id = ?", author.ID).Updates(domain.RecipeAuthor{
		Name:  author.Name,
		Photo: author.Photo,
	}).Error
	if err != nil {
		return domain.RecipeAuthorDetail{}, err
	}

	return *author, nil
}

func (repository *AuthorRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, author *domain.RecipeAuthorDetail) error {
	err := db.WithContext(ctx).Delete(&domain.RecipeAuthor{}, author.ID).Error
	if err != nil {
		return err
	}
	return err
}

func (repository *AuthorRepositoryImpl) FindById(ctx context.Context, db *gorm.DB, id int) (domain.RecipeAuthorDetail, error) {
	var author domain.RecipeAuthorDetail
	err := db.WithContext(ctx).Preload("Recipes").First(&author, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.RecipeAuthorDetail{}, fmt.Errorf("author with ID %d not found", id)
		}
		return domain.RecipeAuthorDetail{}, err
	}
	return author, nil
}

func (repository *AuthorRepositoryImpl) FindAll(ctx context.Context, db *gorm.DB) ([]domain.RecipeAuthor, error) {
	var authors []domain.RecipeAuthor
	err := db.WithContext(ctx).Find(&authors).Error
	if err != nil {
		return nil, err
	}
	return authors, nil
}
