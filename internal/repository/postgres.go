package repository

import (
	"context"
	"petsproject/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreatePet(ctx context.Context, pet domain.Pet) (domain.Pet, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO pets (name, species, created_at)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		pet.Name,
		pet.Species,
		pet.CreatedAt,
	).Scan(&pet.ID)

	if err != nil {
		return domain.Pet{}, err
	}

	return pet, nil
}

func (r *PostgresRepository) GetPet(ctx context.Context, id int64) (domain.Pet, error) {
	var pet domain.Pet

	row := r.db.QueryRow(ctx,
		`SELECT id, name, species, created_at FROM pets WHERE id = $1`,
		id,
	)

	err := row.Scan(&pet.ID, &pet.Name, &pet.Species, &pet.CreatedAt)
	if err != nil {
		return domain.Pet{}, err
	}

	return pet, nil
}

func (r *PostgresRepository) ListPets(ctx context.Context) ([]domain.Pet, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, species, created_at FROM pets`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []domain.Pet
	for rows.Next() {
		var pet domain.Pet
		if err := rows.Scan(&pet.ID, &pet.Name, &pet.Species, &pet.CreatedAt); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}
	return pets, nil
}
