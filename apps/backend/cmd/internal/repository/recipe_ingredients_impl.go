package repository

import (
	"context"
	"fmt"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type RecipeIngredientIRepositoryImpl struct {
}

func NewRecipeIngredients() RecipeIngredientRepository {
	return &RecipeIngredientIRepositoryImpl{}
}

// Create implements RecipeIngredientRepository.
func (r *RecipeIngredientIRepositoryImpl) Create(ctx context.Context, db *gorm.DB, recipeId int, ingredientId int) error {
	// Validasi ID
	if recipeId == 0 || ingredientId == 0 {
		return fmt.Errorf("invalid recipe_id or ingredient_id")
	}

	// Insert data ke junction table
	return db.Exec(`
		INSERT INTO recipe_ingredients (recipe_id, ingredient_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE recipe_id = recipe_id
	`, recipeId, ingredientId).Error
}

// DeleteByRecipeID implements RecipeIngredientRepository.
func (r *RecipeIngredientIRepositoryImpl) DeleteByRecipeID(ctx context.Context, db *gorm.DB, recipeID int) error {
	return db.Where("recipe_id = ?", recipeID).Delete(&domain.RecipeIngredient{}).Error
}

// GetIngredientsByRecipeId implements RecipeIngredientRepository.
func (r *RecipeIngredientIRepositoryImpl) GetIngredientsByRecipeId(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	err := db.Joins("JOIN recipe_ingredients ri ON ri.ingredient_id = ingredients.id").
		Where("ri.recipe_id = ?", recipeId).
		Find(&ingredients).Error
	return ingredients, err
}
