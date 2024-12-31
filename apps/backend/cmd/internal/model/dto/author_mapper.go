package dto

import "github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

func ToAuthorResponse(author domain.RecipeAuthorDetail) AuthorResponses {
	return AuthorResponses{
		ID:    author.ID,
		Name:  author.Name,
		Photo: author.Photo,
	}
}
