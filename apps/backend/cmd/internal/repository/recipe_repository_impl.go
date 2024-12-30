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
	err := db.WithContext(ctx).Create(recipe).Error
	if err != nil {
		return domain.Recipe{}, err
	}

	return *recipe, nil
}

func (repository *RecipeRespositoryImpl) Update(ctx context.Context, db *gorm.DB, recipe *domain.RecipeDetail) (domain.RecipeDetail, error) {
	// Mulai transaksi
	tx := db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return *recipe, err
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
			"Ingredients":    recipe.Ingredients,
		}).Error
	if err != nil {
		// Batalkan transaksi
		tx.Rollback()
		return *recipe, err
	}

	log.Println("cek category id", recipe.CategoryId)

	// Hapus relasi lama
	err = tx.Where("recipe_id = ?", recipe.ID).Delete(&domain.Photo{}).Error
	if err != nil {
		tx.Rollback()
		return *recipe, err
	}

	err = tx.Where("recipe_id = ?", recipe.ID).Delete(&domain.RecipeTutorial{}).Error
	if err != nil {
		tx.Rollback()
		return *recipe, err
	}

	var count int64
	if err := db.Model(&domain.Categories{}).Where("id = ?", recipe.CategoryId).Count(&count).Error; err != nil {
		return *recipe, err
	}
	if count == 0 {
		return *recipe, errors.New("invalid category_id: category does not exist")
	}

	//debuging category
	log.Printf("Updating Recipe ID: %d with Category ID: %d", recipe.ID, recipe.CategoryId)

	//mengecek jika category adalah null
	if recipe.CategoryId == 0 {
		return *recipe, errors.New("CategoryId cannot be zero or null")
	}

	// Tambahkan relasi baru
	for _, photo := range recipe.RecipePhoto {
		photo.RecipeId = recipe.ID
		if err := tx.Create(&photo).Error; err != nil {
			tx.Rollback()
			return *recipe, err
		}
	}

	for _, tutorial := range recipe.RecipeTutorial {
		tutorial.RecipeId = recipe.ID
		if err := tx.Create(&tutorial).Error; err != nil {
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
