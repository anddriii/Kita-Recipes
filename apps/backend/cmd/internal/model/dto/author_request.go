package dto

import "mime/multipart"

type AuthorRequest struct {
	ID    int                   `json:"id"`
	Name  string                `json:"name"`
	Photo *multipart.FileHeader `json:"photo"`
}
