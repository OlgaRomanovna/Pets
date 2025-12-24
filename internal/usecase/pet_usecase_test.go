package usecase

import (
	"context"
	"testing"

	"petsproject/internal/repository"
	"petsproject/internal/service"
)

func TestPetUsecase_CreateAndGetPet(t *testing.T) {
	ctx := context.Background()

	memRepo := repository.NewMemoryRepository()
	cache := service.NewRedisCache("localhost:6379")
	petUC := NewPetUsecase(memRepo, cache)

	pet, err := petUC.CreatePet(ctx, "Buddy", "Dog")
	if err != nil {
		t.Fatalf("CreatePet failed: %v", err)
	}

	if pet.Name != "Buddy" || pet.Species != "Dog" {
		t.Errorf("unexpected pet: %+v", pet)
	}

	got, err := petUC.GetPet(ctx, pet.ID)
	if err != nil {
		t.Fatalf("GetPet failed: %v", err)
	}

	if got.ID != pet.ID || got.Name != pet.Name {
		t.Errorf("got %+v, want %+v", got, pet)
	}
}

func TestPetUsecase_CreatePetValidation(t *testing.T) {
	ctx := context.Background()
	memRepo := repository.NewMemoryRepository()
	cache := service.NewRedisCache("localhost:6379")
	petUC := NewPetUsecase(memRepo, cache)

	_, err := petUC.CreatePet(ctx, "", "")
	if err == nil {
		t.Errorf("expected error for empty name/species, got nil")
	}
}
