package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type RecipeTutorials interface {
	Save(ctx context.Context, db *gorm.DB, recipeTutorial domain.RecipeTutorial) error
	Update(ctx context.Context, db *gorm.DB, recipeID int, tutorials []string) error
	Show(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.RecipeTutorial, error)
}
