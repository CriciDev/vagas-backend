package devs

import (
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("developer not found")
	ErrValidation = errors.New("developer validation failed")
)

type (
	Developer struct {
		ID        int64     `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Skills    []string  `json:"skills"`
		Available bool      `json:"available"`
		Bio       string    `json:"bio"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	SaveDeveloperRequest struct {
		Name      string   `json:"name"`
		Email     string   `json:"email"`
		Skills    []string `json:"skills"`
		Available bool     `json:"available"`
		Bio       string   `json:"bio"`
	}
)
