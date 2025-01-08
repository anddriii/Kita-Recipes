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
	err := db.WithContext(ctx).Model(&domain.Categories{}).Where("id = ?", category.ID).Updates(domain.Categories{
		Name: category.Name,
		Slug: category.Slug,
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

func (repository *CategoryRepositoryImpl) FindById(ctx context.Context, db *gorm.DB, id int) (domain.CategoryDetail, error) {
	var category domain.CategoryDetail
	err := db.WithContext(ctx).Preload("Recipes").Where("id = ?", id).Find(&category).Error
	if err != nil {
		return category, err
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

func (repository *CategoryRepositoryImpl) FindByIds(ctx context.Context, db *gorm.DB, ids []int) (map[int]domain.CategoryDetail, error) {
	var categories []domain.CategoryDetail
	err := db.WithContext(ctx).Where("id IN ?", ids).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[int]domain.CategoryDetail)
	for _, category := range categories {
		categoryMap[int(category.ID)] = category
	}
	return categoryMap, nil
}
