package dto

import "mime/multipart"

type CategoryRequest struct {
	ID   int64                 `json:"id"`
	Name string                `validate:"required" json:"name"`
	Slug string                `json:"slug"`
	Icon *multipart.FileHeader `json:"icon"`
}

// type CategoryUpdateRequest struct {
// 	ID   int64  `validate:"required" json:"id"`
// 	Name string `validate:"required" json:"name"`
// 	Slug string
// 	Icon *multipart.FileHeader `json:"icon"`
// }
