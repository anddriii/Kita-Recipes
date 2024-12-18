package dto

import "mime/multipart"

type IngredientRequest struct {
	ID    int
	Name  string                `json:"name"`
	Photo *multipart.FileHeader `json:"photo"`
}
