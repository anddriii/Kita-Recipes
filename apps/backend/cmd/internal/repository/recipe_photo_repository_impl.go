package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type RecipePhotoRepositoryImpl struct {
}

func NewRecipePhotoRepository() RecipePhotoRepository {
	return &RecipePhotoRepositoryImpl{}
}

// Save implements RecipePhotoRepository.
func (r *RecipePhotoRepositoryImpl) Save(ctx context.Context, db *gorm.DB, recipePhoto domain.Photo) error {
	return db.WithContext(ctx).Create(&recipePhoto).Error
}
