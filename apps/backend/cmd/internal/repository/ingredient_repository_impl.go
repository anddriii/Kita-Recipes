package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type IngredientRepositoryImpl struct {
}

func NewIngredientRepository() IngredientRepository {
	return &IngredientRepositoryImpl{}
}

func (repository *IngredientRepositoryImpl) Save(ctx context.Context, db *gorm.DB, ingredient *domain.Ingredient) (domain.Ingredient, error) {
	err := db.WithContext(ctx).Create(ingredient).Error
	if err != nil {
		return *ingredient, err
	}
	return *ingredient, nil
}

func (repository *IngredientRepositoryImpl) Update(ctx context.Context, db *gorm.DB, ingredient *domain.Ingredient) (domain.Ingredient, error) {
	// Lakukan update menggunakan context
	result := db.WithContext(ctx).Model(&domain.Ingredient{}).Where("id = ?", ingredient.ID).Updates(domain.Ingredient{
		Name:  ingredient.Name,
		Photo: ingredient.Photo,
	})

	// Periksa apakah ada error
	if result.Error != nil {
		return domain.Ingredient{}, result.Error
	}

	// Kembalikan data yang telah diperbarui
	return *ingredient, nil
}

func (repository *IngredientRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, ingredient domain.Ingredient) error {
	err := db.WithContext(ctx).Delete(&domain.Ingredient{}, ingredient.ID).Error
	if err != nil {
		return err
	}
	return err
}

func (repository *IngredientRepositoryImpl) FindById(ctx context.Context, db *gorm.DB, id int) (domain.Ingredient, error) {
	var ingredient domain.Ingredient
	err := db.WithContext(ctx).First(&ingredient, id).Error
	if err != nil {
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}

func (repository *IngredientRepositoryImpl) FindAll(ctx context.Context, db *gorm.DB) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	err := db.WithContext(ctx).Find(&ingredients).Error
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}
