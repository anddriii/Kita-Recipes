package dto

type AuthorResponses struct {
	ID    int64  `validate:"required" json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type AuthorResponseDetail struct {
	ID      int64                   `validate:"required" json:"id"`
	Name    string                  `json:"name"`
	Photo   string                  `json:"photo"`
	Recipes []*RecipeResponseDetail `json:"recipes"`
}
