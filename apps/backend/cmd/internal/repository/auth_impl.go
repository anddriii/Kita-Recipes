package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
)

type AuthRepositoryImpl struct {
}

func NewAuthRepository() AuthRepo {
	return &AuthRepositoryImpl{}
}

// CreateUSer implements AuthRepo.
func (a *AuthRepositoryImpl) CreateUSer(ctx context.Context, name string, photo string, username string, email string, hasPassword string, role string) {
	panic("unimplemented")
}

// GetUserByName implements AuthRepo.
func (a *AuthRepositoryImpl) GetUserByName(ctx context.Context, usernamename string) (*domain.Login, error) {
	panic("unimplemented")
}
