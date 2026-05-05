package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CriciumaDevJobs/backend/handlers"
	"github.com/CriciumaDevJobs/backend/internal/devs"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(dev *devs.FindDevByEmailRow, expiration time.Time) (string, *handlers.ErrorResponse) {

	claims := jwt.MapClaims{
		"user_id":    dev.ID,
		"expires_at": expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte("placeholder-key"))

	if err != nil {
		log.Printf("ERRO: Erro ao assinar token JWT! Message %s", err.Error())
		return "", handlers.NewError(http.StatusInternalServerError, "Erro interno!")
	}

	return signedToken, nil
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Token ausente ou malformatado", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("placeholder-key"), nil
		})

		if err != nil || !token.Valid {
			handlers.ResponseWithHttpError(w, http.StatusUnauthorized, "Token Inválido ou Expirado!")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			r = r.WithContext(ctx)
		}

		next(w, r)
	}
}
