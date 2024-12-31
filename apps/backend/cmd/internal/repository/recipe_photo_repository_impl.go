package repository

import (
	"context"
	"log"
	"time"

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
	log.Printf("Saving photos recipes: %+v", recipePhoto)
	if err := db.Create(&recipePhoto).Error; err != nil {
		log.Printf("Error saving photos recipes: %v", err)
		return err
	}
	return nil
}

func (repository *RecipePhotoRepositoryImpl) Update(ctx context.Context, db *gorm.DB, recipeID int, photos []string) error {
	tx := db.Begin() // Mulai transaksi
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback jika terjadi panic
		}
	}()

	// Hapus foto lama
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&domain.Photo{}).Error; err != nil {
		tx.Rollback() // Rollback jika error
		log.Printf("Error deleting old photos: %v", err)
		return err
	}

	// Masukkan foto baru
	for _, photo := range photos {
		newPhoto := domain.Photo{
			RecipeId:  int64(recipeID),
			Photo:     photo,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := tx.Create(&newPhoto).Error; err != nil {
			tx.Rollback() // Rollback jika error
			log.Printf("Error inserting new photo: %v", err)
			return err
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v", err)
		return err
	}

	log.Println("Transaction committed successfully")
	return nil
}

func (r *RecipePhotoRepositoryImpl) Show(ctx context.Context, db *gorm.DB, recipeId int) ([]domain.Photo, error) {
	var photos []domain.Photo
	err := db.WithContext(ctx).Where("recipe_id = ?", recipeId).Find(&photos).Error
	return photos, err
}
