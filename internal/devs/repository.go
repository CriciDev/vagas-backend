package devs

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Repository interface {
	Create(ctx context.Context, developer Developer) (Developer, error)
	List(ctx context.Context) ([]Developer, error)
	FindByID(ctx context.Context, id int64) (Developer, error)
	Update(ctx context.Context, developer Developer) (Developer, error)
	Delete(ctx context.Context, id int64) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) Create(ctx context.Context, developer Developer) (Developer, error) {
	row := repo.db.QueryRowContext(ctx, `
		INSERT INTO developers (name, email, skills, available, bio)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, email, skills, available, bio, created_at, updated_at
	`, developer.Name, developer.Email, pq.Array(developer.Skills), developer.Available, developer.Bio)

	return scanDeveloper(row)
}

func (repo *PostgresRepository) List(ctx context.Context) ([]Developer, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT id, name, email, skills, available, bio, created_at, updated_at
		FROM developers
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	developers := []Developer{}
	for rows.Next() {
		developer, err := scanDeveloper(rows)
		if err != nil {
			return nil, err
		}
		developers = append(developers, developer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return developers, nil
}

func (repo *PostgresRepository) FindByID(ctx context.Context, id int64) (Developer, error) {
	row := repo.db.QueryRowContext(ctx, `
		SELECT id, name, email, skills, available, bio, created_at, updated_at
		FROM developers
		WHERE id = $1
	`, id)

	developer, err := scanDeveloper(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Developer{}, ErrNotFound
	}
	if err != nil {
		return Developer{}, err
	}

	return developer, nil
}

func (repo *PostgresRepository) Update(ctx context.Context, developer Developer) (Developer, error) {
	row := repo.db.QueryRowContext(ctx, `
		UPDATE developers
		SET name = $2, email = $3, skills = $4, available = $5, bio = $6, updated_at = now()
		WHERE id = $1
		RETURNING id, name, email, skills, available, bio, created_at, updated_at
	`, developer.ID, developer.Name, developer.Email, pq.Array(developer.Skills), developer.Available, developer.Bio)

	updated, err := scanDeveloper(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Developer{}, ErrNotFound
	}
	if err != nil {
		return Developer{}, err
	}

	return updated, nil
}

func (repo *PostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := repo.db.ExecContext(ctx, `DELETE FROM developers WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanDeveloper(row scanner) (Developer, error) {
	var developer Developer
	var skills pq.StringArray
	err := row.Scan(&developer.ID, &developer.Name, &developer.Email, &skills, &developer.Available, &developer.Bio, &developer.CreatedAt, &developer.UpdatedAt)
	if err != nil {
		return Developer{}, err
	}
	developer.Skills = []string(skills)
	return developer, nil
}
