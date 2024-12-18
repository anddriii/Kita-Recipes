package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

	"gorm.io/gorm"
)

type AuthorRepository interface {
	Save(ctx context.Context, db *gorm.DB, author *domain.RecipeAuthor) (domain.RecipeAuthor, error)
	Update(ctx context.Context, db *gorm.DB, author *domain.RecipeAuthorDetail) (domain.RecipeAuthorDetail, error)
	Delete(ctx context.Context, db *gorm.DB, author *domain.RecipeAuthorDetail) error
	FindById(ctx context.Context, db *gorm.DB, id int) (domain.RecipeAuthorDetail, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]domain.RecipeAuthor, error)
}
