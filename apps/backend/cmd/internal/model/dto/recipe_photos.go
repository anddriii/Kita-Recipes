package dto

import "mime/multipart"

type RecipePhotos struct {
	ID    int64 `json:"id"`
	Photo multipart.FileHeader
}

type RecipePhotosResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"photo"`
}

// type PhotoUpload struct {
// 	ID       int64  `json:"id"`
// 	Filename string `json:"photo"`
// 	File     multipart.FileHeader
// }
