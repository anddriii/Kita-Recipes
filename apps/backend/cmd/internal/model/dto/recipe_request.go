package dto

import "mime/multipart"

type RecipeRequestCreate struct {
	Name           string                `validate:"required" json:"name"`
	Slug           string                `json:"slug"`
	Thumbnail      *multipart.FileHeader `validate:"required" json:"thumbnail"`
	About          string                `validate:"required" json:"about"`
	UrlFile        string                `json:"url_file"`
	UrlVideo       string                `validate:"required" json:"url_video"`
	CategoryId     int                   `json:"category_id"`
	RecipeAuthorId int                   `json:"recipe_author_id"`
	RecipePhotos   []RecipePhotos        `json:"photos"`
}

type RecipeRequestUpdate struct {
	ID         int64                 ` json:"id"`
	Name       string                `validate:"required" json:"name"`
	Slug       string                `json:"slug"`
	Thumbnail  *multipart.FileHeader `json:"thumbnail"`
	About      string                `validate:"required" json:"about"`
	UrlFile    string                `json:"url_file"`
	UrlVideo   string                `json:"url_video"`
	CategoryId int                   `json:"category_id" form:"category_id"`
	// Category       *CategoryRequest `json:"category_id"`
	RecipeAuthorId int            `json:"recipe_author_id"`
	RecipePhotos   []RecipePhotos `json:"photos"`
}
