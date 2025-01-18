package dto

import "mime/multipart"

type RegisterDTO struct {
	Name     string                `json:"name" binding:"required"`
	Photo    *multipart.FileHeader `json:"photo"`
	Username string                `json:"username" binding:"required"`
	Password string                `json:"password" binding:"required"`
	Email    string                `json:"email" binding:"required,email"`
	Role     string                `json:"role"`
}

type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
