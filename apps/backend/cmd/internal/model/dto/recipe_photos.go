package dto

import "mime/multipart"

type RecipePhotos struct {
	ID    int64  `json:"id"`
	Photo string `json:"photo"`
}

type PhotoUpload struct {
	ID       int64  `json:"id"`
	Filename string `json:"photo"`
	File     multipart.FileHeader
}
