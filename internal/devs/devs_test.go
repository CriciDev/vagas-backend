package devs

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type memoryDeveloperRepository struct {
	nextID int64
	items  map[int64]Developer
}

func newMemoryDeveloperRepository() *memoryDeveloperRepository {
	return &memoryDeveloperRepository{nextID: 1, items: map[int64]Developer{}}
}

func (repo *memoryDeveloperRepository) Create(ctx context.Context, developer Developer) (Developer, error) {
	developer.ID = repo.nextID
	repo.nextID++
	repo.items[developer.ID] = developer
	return developer, nil
}

func (repo *memoryDeveloperRepository) List(ctx context.Context) ([]Developer, error) {
	developers := []Developer{}
	for _, developer := range repo.items {
		developers = append(developers, developer)
	}
	return developers, nil
}

func (repo *memoryDeveloperRepository) FindByID(ctx context.Context, id int64) (Developer, error) {
	developer, ok := repo.items[id]
	if !ok {
		return Developer{}, ErrNotFound
	}
	return developer, nil
}

func (repo *memoryDeveloperRepository) Update(ctx context.Context, developer Developer) (Developer, error) {
	if _, ok := repo.items[developer.ID]; !ok {
		return Developer{}, ErrNotFound
	}
	repo.items[developer.ID] = developer
	return developer, nil
}

func (repo *memoryDeveloperRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := repo.items[id]; !ok {
		return ErrNotFound
	}
	delete(repo.items, id)
	return nil
}

func TestDeveloperCRUDRoutes(t *testing.T) {
	router := newDeveloperRouter(newMemoryDeveloperRepository(), passMiddleware(), passMiddleware())

	created := requestJSON(t, router, http.MethodPost, "/developers", SaveDeveloperRequest{
		Name:      "Ada Lovelace",
		Email:     "ada@test.local",
		Skills:    []string{"Go", "Angular"},
		Available: true,
		Bio:       "Backend developer",
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, created.Code)
	}

	var developer Developer
	if err := json.Unmarshal(created.Body.Bytes(), &developer); err != nil {
		t.Fatalf("expected developer response, got %v", err)
	}
	if developer.ID == 0 {
		t.Fatal("expected generated developer id")
	}

	listed := requestJSON(t, router, http.MethodGet, "/developers", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listed.Code)
	}

	found := requestJSON(t, router, http.MethodGet, "/developers/1", nil)
	if found.Code != http.StatusOK {
		t.Fatalf("expected lookup status %d, got %d", http.StatusOK, found.Code)
	}

	updated := requestJSON(t, router, http.MethodPut, "/developers/1", SaveDeveloperRequest{
		Name:      "Ada Byron",
		Email:     "ada@test.local",
		Skills:    []string{"Go"},
		Available: false,
		Bio:       "Community developer",
	})
	if updated.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d", http.StatusOK, updated.Code)
	}

	deleted := requestJSON(t, router, http.MethodDelete, "/developers/1", nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("expected delete status %d, got %d", http.StatusNoContent, deleted.Code)
	}

	missing := requestJSON(t, router, http.MethodGet, "/developers/1", nil)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("expected not found status %d, got %d", http.StatusNotFound, missing.Code)
	}
}

func TestDeveloperValidationError(t *testing.T) {
	router := newDeveloperRouter(newMemoryDeveloperRepository(), passMiddleware(), passMiddleware())

	recorder := requestJSON(t, router, http.MethodPost, "/developers", SaveDeveloperRequest{
		Name:  "",
		Email: "invalid",
	})
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected validation status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestDeveloperWriteAuthorizationFailure(t *testing.T) {
	router := newDeveloperRouter(newMemoryDeveloperRepository(), abortMiddleware(http.StatusUnauthorized), passMiddleware())

	recorder := requestJSON(t, router, http.MethodPost, "/developers", SaveDeveloperRequest{
		Name:      "Ada Lovelace",
		Email:     "ada@test.local",
		Skills:    []string{"Go"},
		Available: true,
	})
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func newDeveloperRouter(repo Repository, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	service := NewService(repo)
	handler := NewHandler(service)
	handler.RegisterRoutes(&router.RouterGroup, authMiddleware, adminMiddleware)
	return router
}

func requestJSON(t *testing.T, router *gin.Engine, method string, path string, payload any) *httptest.ResponseRecorder {
	t.Helper()

	body := bytes.NewBuffer(nil)
	if payload != nil {
		if err := json.NewEncoder(body).Encode(payload); err != nil {
			t.Fatalf("expected encoded payload, got %v", err)
		}
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, body)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	return recorder
}

func passMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func abortMiddleware(status int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatus(status)
	}
}
