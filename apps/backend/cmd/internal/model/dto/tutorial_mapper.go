package dto

import "github.com/anddriii/KitaRecipes/cmd/internal/model/domain"

func ToTutorialResponse(tutorial domain.RecipeTutorial) Tutorials {
	return Tutorials{
		Name: tutorial.Name,
	}
}
