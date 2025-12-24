package repository

import (
	"context"
	"petsproject/internal/domain"
)

type Repository interface {
	CreatePet(ctx context.Context, pet domain.Pet) (domain.Pet, error)
	GetPet(ctx context.Context, id int64) (domain.Pet, error)
	ListPets(ctx context.Context) ([]domain.Pet, error)
}
