package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service struct {
		repo      Repository
		jwtSecret []byte
	}
)

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (service *Service) Login(ctx context.Context, request LoginRequest) (LoginResponse, error) {
	user, err := service.repo.FindByEmail(ctx, request.Email)
	if errors.Is(err, ErrUserNotFound) {
		return LoginResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		return LoginResponse{}, ErrInvalidCredentials
	}

	token, err := service.CreateToken(user)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{Token: token}, nil
}

func (service *Service) CreateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims(user))
	return token.SignedString(service.jwtSecret)
}

func (service *Service) ParseToken(rawToken string) (AuthenticatedUser, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(rawToken, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidCredentials
		}
		return service.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return AuthenticatedUser{}, ErrInvalidCredentials
	}

	return AuthenticatedUser{ID: claims.UserID, Email: claims.Email, Role: claims.Role}, nil
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}
