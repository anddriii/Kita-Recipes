package dto

import "github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

func ToCategoryResponse(category domain.CategoryDetail) CategoryResponse {
	return CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Icon: category.Icon,
	}
}
