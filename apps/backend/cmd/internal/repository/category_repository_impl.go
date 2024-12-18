package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (repository *CategoryRepositoryImpl) Save(ctx context.Context, db *gorm.DB, category *domain.Categories) (domain.Categories, error) {
	err := db.WithContext(ctx).Create(category).Error
	if err != nil {
		panic(err)
	}
	return *category, nil
}

func (repository *CategoryRepositoryImpl) Update(ctx context.Context, db *gorm.DB, category *domain.CategoryDetail) (domain.CategoryDetail, error) {
	err := db.WithContext(ctx).Model(&domain.Categories{}).Where("slug = ?", category.Slug).Updates(domain.Categories{
		Name: category.Name,
		Icon: category.Icon,
	}).Error
	if err != nil {
		panic("Error Repo")
	}

	return *category, nil
}

func (repository *CategoryRepositoryImpl) Delete(ctx context.Context, db *gorm.DB, category *domain.CategoryDetail) error {
	err := db.WithContext(ctx).Delete(&domain.Categories{}, category.ID).Error
	if err != nil {
		panic(err)
	}
	return err
}

func (repository *CategoryRepositoryImpl) FindBySlug(ctx context.Context, db *gorm.DB, slug string) (domain.CategoryDetail, error) {
	var category domain.CategoryDetail
	err := db.WithContext(ctx).Preload("Recipes").Where("slug = ?", slug).First(&category).Error
	if err != nil {
		panic(err)
	}
	return category, nil
}

func (repository *CategoryRepositoryImpl) FindAll(ctx context.Context, db *gorm.DB) ([]domain.Categories, error) {
	var categories []domain.Categories
	err := db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		panic(err)
	}
	return categories, nil
}
