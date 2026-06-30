package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type (
	fakeAuthRepository struct {
		users map[string]User
	}
)

func (repo fakeAuthRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	user, ok := repo.users[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func TestLoginWithValidCredentials(t *testing.T) {
	service := newTestAuthService(t)

	response, err := service.Login(context.Background(), LoginRequest{Email: "admin@test.local", Password: "secret"})
	if err != nil {
		t.Fatalf("expected valid login, got %v", err)
	}
	if response.Token == "" {
		t.Fatal("expected token")
	}

	user, err := service.ParseToken(response.Token)
	if err != nil {
		t.Fatalf("expected parseable token, got %v", err)
	}
	if user.Role != RoleAdmin {
		t.Fatalf("expected role %q, got %q", RoleAdmin, user.Role)
	}
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	service := newTestAuthService(t)

	_, err := service.Login(context.Background(), LoginRequest{Email: "admin@test.local", Password: "wrong"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	service := newTestAuthService(t)
	router := newProtectedRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidToken(t *testing.T) {
	service := newTestAuthService(t)
	router := newProtectedRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer invalid")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestOptionalAuthMiddlewareAllowsMissingToken(t *testing.T) {
	service := newTestAuthService(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/public", OptionalAuthenticate(service), func(ctx *gin.Context) {
		if _, ok := UserFromContext(ctx); ok {
			t.Fatal("expected no authenticated user")
		}
		ctx.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/public", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestOptionalAuthMiddlewareSetsValidUser(t *testing.T) {
	service := newTestAuthService(t)
	token, err := service.CreateToken(User{ID: 1, Email: "admin@test.local", Role: RoleAdmin})
	if err != nil {
		t.Fatalf("expected token, got %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/public", OptionalAuthenticate(service), func(ctx *gin.Context) {
		user, ok := UserFromContext(ctx)
		if !ok {
			t.Fatal("expected authenticated user")
		}
		if user.Role != RoleAdmin {
			t.Fatalf("expected role %q, got %q", RoleAdmin, user.Role)
		}
		ctx.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/public", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}

func TestRoleMiddlewareRejectsForbiddenRole(t *testing.T) {
	service := newTestAuthService(t)
	router := newProtectedRouter(service)
	token, err := service.CreateToken(User{ID: 2, Email: "dev@test.local", Role: "developer"})
	if err != nil {
		t.Fatalf("expected token, got %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func newTestAuthService(t *testing.T) *Service {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("expected password hash, got %v", err)
	}

	repo := fakeAuthRepository{users: map[string]User{
		"admin@test.local": {ID: 1, Email: "admin@test.local", PasswordHash: string(hash), Role: RoleAdmin},
	}}

	return NewService(repo, "test-secret")
}

func newProtectedRouter(service *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", Authenticate(service), RequireRole(RoleAdmin), func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	return router
}
