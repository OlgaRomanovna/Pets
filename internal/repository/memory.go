package repository

import (
	"context"
	"errors"
	"sync"

	"petsproject/internal/domain"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	seq  int64
	pets map[int64]domain.Pet
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		pets: make(map[int64]domain.Pet),
	}
}

func (r *MemoryRepository) CreatePet(ctx context.Context, pet domain.Pet) (domain.Pet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	pet.ID = r.seq

	r.pets[pet.ID] = pet
	return pet, nil
}

func (r *MemoryRepository) GetPet(ctx context.Context, id int64) (domain.Pet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pet, ok := r.pets[id]
	if !ok {
		return domain.Pet{}, errors.New("pet not found")
	}
	return pet, nil
}

func (r *MemoryRepository) ListPets(ctx context.Context) ([]domain.Pet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Pet, 0, len(r.pets))
	for _, pet := range r.pets {
		result = append(result, pet)
	}
	return result, nil
}
