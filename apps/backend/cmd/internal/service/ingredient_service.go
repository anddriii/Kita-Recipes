package service

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
)

type IngredientService interface {
	Save(ctx context.Context, request *dto.IngredientRequest) (dto.IngredientResponse, error)
	Update(ctx context.Context, request *dto.IngredientRequest) (dto.IngredientResponse, error)
	Delete(ctx context.Context, ingredientId int)
	FindById(ctx context.Context, ingredientId int) (dto.IngredientResponse, error)
	FindAll(ctx context.Context) ([]dto.IngredientResponse, error)
}
