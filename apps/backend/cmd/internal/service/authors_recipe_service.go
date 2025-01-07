package service

import (
	"context"

	"github.com/anddriii/KitaRecipes/cmd/internal/model/dto"
)

type AuthorService interface {
	Save(ctx context.Context, request *dto.AuthorRequest) (dto.AuthorResponses, error)
	Update(ctx context.Context, request *dto.AuthorRequest) (dto.AuthorResponses, error)
	Delete(ctx context.Context, authorID int)
	FindById(ctx context.Context, id int) (dto.AuthorResponseDetail, error)
	FindAll(ctx context.Context) ([]dto.AuthorResponses, error)
}
