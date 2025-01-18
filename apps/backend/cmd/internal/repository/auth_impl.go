package repository

import (
	"context"
	"errors"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
}

func NewAuthRepository() AuthRepo {
	return &AuthRepositoryImpl{}
}

// CreateUSer implements AuthRepo.
func (repo *AuthRepositoryImpl) CreateUSer(ctx context.Context, db *gorm.DB, name string, photo string, username string, email string, hasPassword string, role string) error {
	tx := db.Begin()

	//insert ke tabel author
	author := domain.RecipeAuthor{Name: name, Photo: photo}
	if err := tx.WithContext(ctx).Create(&author).Error; err != nil {
		tx.Rollback()
		return err
	}

	//insert ke tabel logins
	login := domain.Login{
		AuthorId: int(author.ID),
		Username: username,
		Email:    email,
		Password: hasPassword,
		Role:     role,
	}

	if err := tx.Create(&login).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetUserByName implements AuthRepo.
func (repo *AuthRepositoryImpl) GetUserByName(ctx context.Context, db *gorm.DB, username string) (*domain.Login, error) {
	var user domain.Login
	err := db.WithContext(ctx).Preload("Author").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
