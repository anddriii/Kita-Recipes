package service

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
)

type RecipeService interface {
	Save(ctx context.Context, request dto.RecipeRequestCreate) (dto.RecipeResponseCreate, error)
	Update(ctx context.Context, request dto.RecipeRequestUpdate) (dto.RecipeResponseUpdate, error)
	Delete(ctx context.Context, recipeId int)
	FindById(ctx context.Context, id int) (dto.RecipeResponseDetail, error)
	FindAll(ctx context.Context) ([]dto.RecipeResponses, error)
}
