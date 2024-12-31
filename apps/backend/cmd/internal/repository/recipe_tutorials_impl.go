package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type NewRecipeTutorialsImpl struct {
}

// Save implements RecipeTutorials.
func (n *NewRecipeTutorialsImpl) Save(ctx context.Context, db *gorm.DB, recipeTutorial domain.RecipeTutorial) error {
	if err := db.Create(&recipeTutorial).Error; err != nil {
		return err
	}

	return nil
}

// Show implements RecipeTutorials.
func (n *NewRecipeTutorialsImpl) Show(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.RecipeTutorial, error) {
	panic("unimplemented")
}

// Updadate implements RecipeTutorials.
func (n *NewRecipeTutorialsImpl) Updadate(ctx context.Context, db *gorm.DB, recipeID int, tutorials []string) error {
	panic("unimplemented")
}

func NewRecipeTutorialsRepository() RecipeTutorials {
	return &NewRecipeTutorialsImpl{}
}
