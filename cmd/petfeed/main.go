package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"petsproject/internal/handlers"
	"petsproject/internal/repository"
	"petsproject/internal/service"
	"petsproject/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewPostgresRepository(db)
	cache := service.NewRedisCache(os.Getenv("REDIS_ADDR"))

	uc := usecase.NewPetUsecase(repo, cache)
	h := handlers.NewHTTPHandler(uc)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Router()))
}
