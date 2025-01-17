package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
)

type AuthRepo interface {
	CreateUSer(ctx context.Context, name string, photo string, username string, email string, hasPassword string, role string)
	GetUserByName(ctx context.Context, usernamename string) (*domain.Login, error)
}
