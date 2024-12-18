package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

	"gorm.io/gorm"
)

type IngredientRepository interface {
	Save(ctx context.Context, db *gorm.DB, ingredient *domain.Ingredient) (domain.Ingredient, error)
	Update(ctx context.Context, db *gorm.DB, ingredient *domain.Ingredient) (domain.Ingredient, error)
	Delete(ctx context.Context, db *gorm.DB, ingredient domain.Ingredient) error
	FindById(ctx context.Context, db *gorm.DB, id int) (domain.Ingredient, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]domain.Ingredient, error)
}
