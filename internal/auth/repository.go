package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type (
	Repository interface {
		FindByEmail(ctx context.Context, email string) (User, error)
	}

	PostgresRepository struct {
		db *sql.DB
	}
)

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repo *PostgresRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	row := repo.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`, email)

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}
