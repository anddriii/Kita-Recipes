package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

	"gorm.io/gorm"
)

type RecipeRespository interface {
	Save(ctx context.Context, db *gorm.DB, recipe *domain.Recipe) (domain.Recipe, error)
	Update(ctx context.Context, db *gorm.DB, recipe *domain.RecipeDetail) (domain.RecipeDetail, error)
	Delete(ctx context.Context, db *gorm.DB, recipe *domain.RecipeDetail) error
	FindById(ctx context.Context, db *gorm.DB, id int) (domain.RecipeDetail, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]domain.RecipeDetail, error)
	RecipeAuthor(ctx context.Context, db *gorm.DB, id int) ([]domain.Recipe, error)
}
