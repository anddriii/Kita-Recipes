package dto

type RecipeIngredient struct {
	RecipeId     int64 `json:"recipe_id"`
	IngredientId int64 `json:"ingredient_id"`
}
