package repository

import (
	"context"
	"errors"
	"log"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type RecipeRespositoryImpl struct {
}

func NewRecipeRepository() RecipeRespository {
	return &RecipeRespositoryImpl{}
}

func (repository *RecipeRespositoryImpl) Save(ctx context.Context, db *gorm.DB, recipe *domain.Recipe) (domain.Recipe, error) {
	// Simpan recipe tanpa tutorial terlebih dahulu
	err := db.WithContext(ctx).Omit("RecipeTutorial").Create(recipe).Error
	if err != nil {
		return domain.Recipe{}, err
	}

	return *recipe, nil
}

func (repository *RecipeRespositoryImpl) Update(ctx context.Context, db *gorm.DB, recipe *domain.RecipeDetail) (domain.RecipeDetail, error) {
	// Mulai transaksi
	log.Println("== DEBUG PAYLOAD ==")
	log.Printf("ID: %d", recipe.ID)
	log.Printf("Name: %s", recipe.Name)
	log.Printf("CategoryId: %d", recipe.CategoryId)
	log.Print("File name: ", recipe.UrlFile)

	log.Printf("CategoryId from Payload: %d", recipe.CategoryId)
	tx := db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return *recipe, err
	}
	log.Printf("CategoryId di Repository sebelum update: %d", recipe.CategoryId)

	//debuging category
	log.Printf("Updating Recipe ID: %d with Category ID: %d", recipe.ID, recipe.CategoryId)
	log.Println("cek category id repository", recipe.CategoryId)

	// Validasi CategoryId di API Layer
	if recipe.CategoryId == 0 {
		return *recipe, errors.New("CategoryId cannot be zero")
	}

	// Validasi category di database
	var categoryCount int64
	if err := db.Model(&domain.Categories{}).Where("id = ?", recipe.CategoryId).Count(&categoryCount).Error; err != nil {
		return *recipe, err
	}
	if categoryCount == 0 {
		return *recipe, errors.New("CategoryId does not exist in categories table")
	}

	//mengecek jika category adalah null
	if recipe.CategoryId == 0 {
		return *recipe, errors.New("CategoryId cannot be zero or null")
	}

	// Update data utama
	err := tx.Model(&domain.RecipeDetail{}).
		Where("id = ?", recipe.ID).
		Updates(map[string]interface{}{
			"Name":           recipe.Name,
			"Slug":           recipe.Slug,
			"Thumbnail":      recipe.Thumbnail,
			"About":          recipe.About,
			"UrlFile":        recipe.UrlFile,
			"UrlVideo":       recipe.UrlVideo,
			"CategoryId":     recipe.CategoryId,
			"RecipeAuthorId": recipe.RecipeAuthorId,
		}).Error
	if err != nil {
		// Batalkan transaksi
		tx.Rollback()
		return *recipe, err
	}

	// Hapus tutorial lama
	if err := tx.WithContext(ctx).Where("recipe_id = ?", recipe.ID).Delete(&domain.RecipeTutorial{}).Error; err != nil {
		tx.Rollback()
		return *recipe, err
	}

	// Tambahkan tutorial baru
	if len(recipe.RecipeTutorial) > 0 {
		for i := range recipe.RecipeTutorial {
			recipe.RecipeTutorial[i].RecipeId = recipe.ID
		}
		if err := tx.WithContext(ctx).Create(&recipe.RecipeTutorial).Error; err != nil {
			tx.Rollback()
			return *recipe, err
		}
	}

	// Commit transaksi jika semuanya berhasil
	if err := tx.Commit().Error; err != nil {
		return *recipe, err
	}

	return *recipe, nil
}

func (repository *RecipeRespositoryImpl) Delete(ctx context.Context, db *gorm.DB, recipe *domain.RecipeDetail) error {
	err := db.WithContext(ctx).Delete(&domain.Recipe{}, recipe.ID).Error
	if err != nil {
		panic(err)
	}
	return err
}

func (repository *RecipeRespositoryImpl) FindById(ctx context.Context, db *gorm.DB, id int) (domain.RecipeDetail, error) {
	var recipe domain.RecipeDetail
	err := db.WithContext(ctx).Preload("RecipeTutorial").Preload("Ingredients").Preload("Category").Preload("RecipeAuthor").Preload("RecipePhoto").Where("id = ?", id).First(&recipe).Error
	if err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (repository *RecipeRespositoryImpl) FindAll(ctx context.Context, db *gorm.DB) ([]domain.RecipeDetail, error) {
	var recipes []domain.RecipeDetail
	err := db.WithContext(ctx).Preload("RecipePhoto").Preload("Category").Find(&recipes).Error
	if err != nil {
		panic(err)
	}
	return recipes, nil
}
