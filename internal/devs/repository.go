package devs

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type DevRepository struct {
	DB *sql.DB
}

func NewDevRepository(db *sql.DB) *DevRepository {
	repository := DevRepository{
		DB: db,
	}

	return &repository
}

func (repository *DevRepository) CreateDev(ctx context.Context, dev *Dev) error {
	query := "INSERT INTO devs (name, email, skills, bio, availability, socials) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	return repository.DB.QueryRowContext(ctx, query, dev.Name, dev.Email, dev.Skills, dev.Bio, dev.Availability, pq.Array(dev.Socials)).Scan(&dev.ID)
}
