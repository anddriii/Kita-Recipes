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
	Tutorials      []Tutorials           `json:"tutorials" form:"tutorials"`
}

type RecipeRequestUpdate struct {
	ID             int64                 `json:"id" form:"id"`
	Name           string                `validate:"required" json:"name" form:"name"`
	Slug           string                `json:"slug"`
	Thumbnail      *multipart.FileHeader `json:"thumbnail" form:"thumbnail"`
	About          string                `validate:"required" json:"about" form:"about"`
	UrlFile        string                `json:"url_file" form:"url_file"`
	UrlVideo       string                `json:"url_video" form:"url_video"`
	CategoryId     int                   `json:"category_id" form:"category_id"`
	Tutorials      []Tutorials           `json:"tutorials" form:"tutorials"`
	RecipeAuthorId int                   `json:"recipe_author_id" form:"recipe_author_id"`
	RecipePhotos   []RecipePhotos        `json:"photos" form:"recipe_photos"`
}
