package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type RecipeIngredientRepository interface {
	Create(ctx context.Context, db *gorm.DB, recipeId int, ingredientId int) error
	GetIngredientsByRecipeId(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.Ingredient, error)
	DeleteByRecipeID(ctx context.Context, db *gorm.DB, recipeID int) error
}
