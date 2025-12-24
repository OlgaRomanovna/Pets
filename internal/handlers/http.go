package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"petsproject/internal/usecase"

	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	uc *usecase.PetUsecase
}

func NewHTTPHandler(uc *usecase.PetUsecase) *HTTPHandler {
	return &HTTPHandler{uc: uc}
}

func (h *HTTPHandler) Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/pets", h.CreatePet).Methods(http.MethodPost)
	r.HandleFunc("/pets/{id}", h.GetPet).Methods(http.MethodGet)
	r.HandleFunc("/pets", h.ListPets).Methods(http.MethodGet)
	return r
}

func (h *HTTPHandler) CreatePet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name    string `json:"name"`
		Species string `json:"species"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	pet, err := h.uc.CreatePet(r.Context(), req.Name, req.Species)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // ← вместо 500
		return
	}
	_ = json.NewEncoder(w).Encode(pet)
}

func (h *HTTPHandler) GetPet(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// Если ID не число → 404 (не найден)
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	pet, err := h.uc.GetPet(r.Context(), id)
	if err != nil {
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(pet)
}

func (h *HTTPHandler) ListPets(w http.ResponseWriter, r *http.Request) {
	pets, err := h.uc.ListPets(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(pets)
}
