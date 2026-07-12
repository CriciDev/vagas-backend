package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const userContextKey = "auth.user"

type (
	Handler struct {
		service *Service
	}
)

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/auth/login", handler.Login)
}

func (handler *Handler) Login(ctx *gin.Context) {
	var request LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid login payload"})
		return
	}

	response, err := handler.service.Login(ctx.Request.Context(), request)
	if errors.Is(err, ErrInvalidCredentials) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func Authenticate(service *Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := bearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		user, err := service.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			return
		}

		ctx.Set(userContextKey, user)
		ctx.Next()
	}
}

func OptionalAuthenticate(service *Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			return
		}

		token, ok := bearerToken(header)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		user, err := service.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			return
		}

		ctx.Set(userContextKey, user)
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := UserFromContext(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}
		if user.Role != role {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		ctx.Next()
	}
}

func UserFromContext(ctx *gin.Context) (AuthenticatedUser, bool) {
	value, exists := ctx.Get(userContextKey)
	if !exists {
		return AuthenticatedUser{}, false
	}

	user, ok := value.(AuthenticatedUser)
	return user, ok
}
