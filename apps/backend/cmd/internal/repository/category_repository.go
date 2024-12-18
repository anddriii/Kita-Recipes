package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Save(ctx context.Context, db *gorm.DB, category *domain.Categories) (domain.Categories, error)
	Update(ctx context.Context, db *gorm.DB, category *domain.CategoryDetail) (domain.CategoryDetail, error)
	Delete(ctx context.Context, db *gorm.DB, category *domain.CategoryDetail) error
	FindBySlug(ctx context.Context, db *gorm.DB, slug string) (domain.CategoryDetail, error)
	FindAll(ctx context.Context, db *gorm.DB) ([]domain.Categories, error)
}
