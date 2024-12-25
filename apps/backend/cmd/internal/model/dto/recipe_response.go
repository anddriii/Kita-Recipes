package dto

type RecipeResponses struct {
	ID                   int64  `validate:"required" json:"id"`
	Name                 string `validate:"required" json:"name"`
	Slug                 string `json:"slug"`
	Thumbnail            string `validate:"required" json:"thumbnail"`
	About                string `validate:"required" json:"about"`
	UrlFile              string `json:"url_file"`
	UrlVideo             string `validate:"required" json:"url_video"`
	CategoryId           int
	CategoryResponse     *CategoryResponse `json:"category"`
	RecipePhotos         []*RecipePhotos
	RecipePhotosResponse []*RecipePhotosResponse `json:"photos"`
}

type RecipeResponseDetail struct {
	ID               int64                 `validate:"required" json:"id"`
	Name             string                `validate:"required" json:"name"`
	Slug             string                `json:"slug"`
	Thumbnail        string                `validate:"required" json:"thumbnail"`
	About            string                `validate:"required" json:"about"`
	UrlFile          string                `json:"url_file"`
	UrlVideo         string                `validate:"required" json:"url_video"`
	CategoryResponse *CategoryResponse     `json:"category"`
	Author           *AuthorResponses      `json:"author"`
	Ingredients      []*IngredientResponse `json:"recipe_ingredient"`
	RecipePhotos     []*RecipePhotos       `json:"photos"`
	RecipeTutorials  []*Tutorials          `json:"tutorials"`
}

type RecipeResponseCreate struct {
	ID                   int64                   `validate:"required" json:"id"`
	Name                 string                  `validate:"required" json:"name"`
	Slug                 string                  `json:"slug"`
	Thumbnail            string                  `validate:"required" json:"thumbnail"`
	About                string                  `validate:"required" json:"about"`
	UrlFile              string                  `json:"url_file"`
	UrlVideo             string                  `validate:"required" json:"url_video"`
	CategoryResponse     *CategoryResponse       `json:"category"`
	RecipePhotosResponse []*RecipePhotosResponse `json:"photos"`
}
