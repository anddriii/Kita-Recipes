package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/domain"
	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
	"github.com/anddriii/KitaRecipes/cmd/internal/repository"
	"github.com/anddriii/KitaRecipes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	AuthRepo repository.AuthRepo
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewAuthService(authRepo repository.AuthRepo, db *gorm.DB, validate *validator.Validate) AuthService {
	return &AuthServiceImpl{
		AuthRepo: authRepo,
		DB:       db,
		Validate: validate,
	}
}

// generate token JWT
func generateToken(user *domain.Login) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"author_id": user.AuthorId,
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), //exp 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

// Login implements AuthService.
func (a *AuthServiceImpl) Login(ctx context.Context, request *dto.LoginDTO) (string, error) {
	user, err := a.AuthRepo.GetUserByName(ctx, a.DB, request.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	//generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// Register implements AuthService.
func (a *AuthServiceImpl) Register(ctx context.Context, request *dto.RegisterDTO) error {
	err := a.Validate.Struct(request)
	if err != nil {
		return err
	}

	basePath, err := filepath.Abs("../../../assets/")
	if err != nil {
		return err
	}

	authorPhoto := domain.RecipeAuthor{
		Name: request.Name,
	}

	filename := uuid.New().String() + "-" + authorPhoto.Name + "." + strings.Split(request.Photo.Filename, ".")[len(strings.Split(request.Photo.Filename, "."))-1]
	photoPath := filepath.Join(basePath, "author_photo", filename)
	if err := utils.SaveFile(request.Photo, photoPath); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return a.AuthRepo.CreateUSer(ctx, a.DB, request.Name, filename, request.Username, request.Email, string(hashedPassword), request.Role)
}
