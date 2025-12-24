package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"petsproject/internal/repository"
	"petsproject/internal/service"
	"petsproject/internal/usecase"
)

func setupTestHandler() *HTTPHandler {
	memRepo := repository.NewMemoryRepository()
	cache := service.NewRedisCache("localhost:6379")
	petUC := usecase.NewPetUsecase(memRepo, cache)
	return NewHTTPHandler(petUC)
}

func TestCreatePetHandler(t *testing.T) {
	h := setupTestHandler()

	body := `{"name":"Fluffy","species":"Cat"}`
	req := httptest.NewRequest("POST", "/pets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreatePet(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var pet map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pet); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if pet["name"] != "Fluffy" || pet["species"] != "Cat" {
		t.Errorf("unexpected pet returned: %+v", pet)
	}
}

func TestCreatePetHandlerValidation(t *testing.T) {
	h := setupTestHandler()

	body := `{"name":"","species":""}`
	req := httptest.NewRequest("POST", "/pets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreatePet(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetPetHandlerNotFound(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest("GET", "/pets/9999", nil) // ID, которого нет
	w := httptest.NewRecorder()

	h.GetPet(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestListPetsHandler_Empty(t *testing.T) {
	h := setupTestHandler()

	req := httptest.NewRequest("GET", "/pets", nil)
	w := httptest.NewRecorder()

	h.ListPets(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var pets []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pets); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(pets) != 0 {
		t.Errorf("expected empty list, got %+v", pets)
	}
}

func TestListPetsHandler_WithPets(t *testing.T) {
	h := setupTestHandler()

	// Создаём питомцев через HTTP POST
	createPet := func(name, species string) {
		body := fmt.Sprintf(`{"name":"%s","species":"%s"}`, name, species)
		req := httptest.NewRequest("POST", "/pets", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.CreatePet(w, req)
		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("failed to create pet %s: status %d", name, w.Result().StatusCode)
		}
	}

	createPet("Fluffy", "Cat")
	createPet("Buddy", "Dog")

	// Получаем список
	req := httptest.NewRequest("GET", "/pets", nil)
	w := httptest.NewRecorder()
	h.ListPets(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var pets []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&pets); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(pets) != 2 {
		t.Errorf("expected 2 pets, got %d", len(pets))
	}

	names := []string{pets[0]["name"].(string), pets[1]["name"].(string)}
	if !(contains(names, "Fluffy") && contains(names, "Buddy")) {
		t.Errorf("unexpected pet names: %v", names)
	}
}

// вспомогательная функция
func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
