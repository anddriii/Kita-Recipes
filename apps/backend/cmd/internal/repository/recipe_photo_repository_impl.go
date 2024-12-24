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

func (r *RecipePhotoRepositoryImpl) Show(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.Photo, error) {
	var photos []domain.Photo
	err := db.WithContext(ctx).Where("recipe_id = ?", recipeId).Find(&photos).Error
	return photos, err
}
