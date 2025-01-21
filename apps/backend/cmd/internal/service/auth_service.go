package service

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
)

type AuthService interface {
	Register(ctx context.Context, request *dto.RegisterDTO) error
	Login(ctx context.Context, request *dto.LoginDTO) (string, error)
}
