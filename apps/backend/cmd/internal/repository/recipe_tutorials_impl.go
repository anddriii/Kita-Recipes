package repository

import (
	"context"
	"log"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type NewRecipeTutorialsImpl struct {
}

func NewRecipeTutorialsRepository() RecipeTutorials {
	return &NewRecipeTutorialsImpl{}
}

// Save implements RecipeTutorials.
func (n *NewRecipeTutorialsImpl) Save(ctx context.Context, db *gorm.DB, recipeTutorial domain.RecipeTutorial) error {
	if err := db.Create(&recipeTutorial).Error; err != nil {
		return err
	}

	return nil
}

// Updadate implements RecipeTutorials.
func (n *NewRecipeTutorialsImpl) Update(ctx context.Context, db *gorm.DB, recipeId int, tutorials []string) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//hapus tutorials lama
	if err := tx.WithContext(ctx).Where("recipe_id = ?", recipeId).Delete(&domain.RecipeTutorial{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting old tutorials: %v", err)
		return err
	}

	// //masukkan tutorials baru
	// for _, tutorial := range tutorials {
	// 	newTutorial := domain.RecipeTutorial{
	// 		RecipeId:  int64(recipeId),
	// 		Name:      tutorial,
	// 		CreatedAt: time.Now(),
	// 		UpdatedAt: time.Now(),
	// 	}
	// 	if err := tx.Create(&newTutorial).Error; err != nil {
	// 		tx.Rollback()
	// 		log.Printf("Error inserting new tutorial: %v", err)
	// 		return err
	// 	}
	// }

	if err := tx.Commit().Error; err != nil {
		log.Printf("Error commiting transaction: %v", err)
		return err
	}

	return nil
}

// Show implements RecipeTutorials.
func (n *NewRecipeTutorialsImpl) Show(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.RecipeTutorial, error) {
	var tutorials []domain.RecipeTutorial
	err := db.WithContext(ctx).Where("recipe_id = ?", recipeId).Find(&tutorials).Error
	return tutorials, err
}
