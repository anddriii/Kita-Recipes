package service

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
)

type CategoryService interface {
	Create(ctx context.Context, request dto.CategoryRequest) (dto.CategoryResponse, error)
	Update(ctx context.Context, request dto.CategoryRequest) (dto.CategoryResponseDetail, error)
	Delete(ctx context.Context, id int)
	FindById(ctx context.Context, id int) (dto.CategoryResponseDetail, error)
	FindAll(ctx context.Context) ([]dto.CategoryResponse, error)
}
