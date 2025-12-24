package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"petsproject/internal/domain"
	"petsproject/internal/repository"
	"petsproject/internal/service"
)

type PetUsecase struct {
	repo  repository.Repository
	cache service.Cache
}

func NewPetUsecase(repo repository.Repository, cache service.Cache) *PetUsecase {
	return &PetUsecase{repo: repo, cache: cache}
}

func (u *PetUsecase) CreatePet(ctx context.Context, name, species string) (domain.Pet, error) {
	if name == "" || species == "" {
		return domain.Pet{}, errors.New("name and species cannot be empty")
	}

	pet := domain.Pet{
		Name:      name,
		Species:   species,
		CreatedAt: time.Now(),
	}

	return u.repo.CreatePet(ctx, pet)
}

func (u *PetUsecase) GetPet(ctx context.Context, id int64) (domain.Pet, error) {
	var pet domain.Pet
	key := fmt.Sprintf("pet:%d", id)

	ok, err := u.cache.Get(ctx, key, &pet)
	if err != nil {
		return domain.Pet{}, err
	}
	if ok {
		return pet, nil
	}

	pet, err = u.repo.GetPet(ctx, id)
	if err != nil {
		return domain.Pet{}, err
	}

	_ = u.cache.Set(ctx, key, pet, time.Minute)
	return pet, nil
}

func (u *PetUsecase) ListPets(ctx context.Context) ([]domain.Pet, error) {
	return u.repo.ListPets(ctx)
}
