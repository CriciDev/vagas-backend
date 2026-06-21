package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const RoleAdmin = "admin"

var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Role         string
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AuthenticatedUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func newClaims(user User) Claims {
	now := time.Now()
	return Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
}
