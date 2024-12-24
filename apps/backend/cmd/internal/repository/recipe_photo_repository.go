package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type RecipePhotoRepository interface {
	Save(ctx context.Context, db *gorm.DB, recipePhoto domain.Photo) error
	Show(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.Photo, error)
}
