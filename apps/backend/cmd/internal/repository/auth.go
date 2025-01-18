package repository

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type AuthRepo interface {
	CreateUSer(ctx context.Context, db *gorm.DB, name string, photo string, username string, email string, hasPassword string, role string) error
	GetUserByName(ctx context.Context, db *gorm.DB, username string) (*domain.Login, error)
}
