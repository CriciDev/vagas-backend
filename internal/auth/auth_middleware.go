package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CriciumaDevJobs/backend/handlers"
	"github.com/CriciumaDevJobs/backend/internal/devs"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserIDKey = "user_id"
)

var (
	JwtSecretKey = getJWTKey()
)

type ContextKey string

func GenerateJwtToken(dev *devs.FindDevByEmailRow, expiration time.Time) (string, *handlers.ErrorResponse) {

	claims := jwt.MapClaims{
		"user_id":    dev.ID,
		"expires_at": expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(JwtSecretKey))

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
			handlers.ResponseWithHttpError(w, http.StatusUnauthorized, "Token ausente ou malformatado")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			handlers.ResponseWithHttpError(w, http.StatusUnauthorized, "Token Inválido ou Expirado!")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ContextKey := UserIDKey
			ctx := context.WithValue(r.Context(), ContextKey, claims[UserIDKey])
			r = r.WithContext(ctx)
		}

		next(w, r)
	}
}

func getJWTKey() string {

	key, ok := os.LookupEnv("JWT_SECRET_KEY")

	if !ok {
		log.Fatalf("ERRO: Variavel de ambiente para a Key JWT não foi definida")
	}

	return key
}
