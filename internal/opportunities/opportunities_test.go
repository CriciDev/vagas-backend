package opportunities

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/CriciumaDevJobs/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type (
	memoryOpportunityRepository struct {
		nextID int64
		items  map[int64]Opportunity
	}
)

func newMemoryOpportunityRepository() *memoryOpportunityRepository {
	return &memoryOpportunityRepository{nextID: 1, items: map[int64]Opportunity{}}
}

func (repo *memoryOpportunityRepository) Create(ctx context.Context, opportunity Opportunity) (Opportunity, error) {
	opportunity.ID = repo.nextID
	repo.nextID++
	repo.items[opportunity.ID] = opportunity
	return opportunity, nil
}

func (repo *memoryOpportunityRepository) List(ctx context.Context, filters ListFilters, pagination Pagination) ([]Opportunity, int, error) {
	matched := []Opportunity{}
	for _, opportunity := range repo.items {
		if opportunity.Status != StatusPublished {
			continue
		}
		if filters.Type != "" && opportunity.Type != filters.Type {
			continue
		}
		if filters.WorkMode != "" && opportunity.WorkMode != filters.WorkMode {
			continue
		}
		if filters.Location != "" && !strings.Contains(strings.ToLower(opportunity.Location), strings.ToLower(filters.Location)) {
			continue
		}
		matched = append(matched, opportunity)
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].ID > matched[j].ID
	})

	total := len(matched)
	start := pagination.Offset()
	if start > total {
		start = total
	}
	end := start + pagination.Limit()
	if end > total {
		end = total
	}

	page := append([]Opportunity{}, matched[start:end]...)
	return page, total, nil
}

func (repo *memoryOpportunityRepository) FindByID(ctx context.Context, id int64) (Opportunity, error) {
	opportunity, ok := repo.items[id]
	if !ok {
		return Opportunity{}, ErrNotFound
	}
	return opportunity, nil
}

func (repo *memoryOpportunityRepository) FindPublishedByID(ctx context.Context, id int64) (Opportunity, error) {
	opportunity, ok := repo.items[id]
	if !ok || opportunity.Status != StatusPublished {
		return Opportunity{}, ErrNotFound
	}
	return opportunity, nil
}

func (repo *memoryOpportunityRepository) Update(ctx context.Context, opportunity Opportunity) (Opportunity, error) {
	if _, ok := repo.items[opportunity.ID]; !ok {
		return Opportunity{}, ErrNotFound
	}
	repo.items[opportunity.ID] = opportunity
	return opportunity, nil
}

func (repo *memoryOpportunityRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := repo.items[id]; !ok {
		return ErrNotFound
	}
	delete(repo.items, id)
	return nil
}

func TestOpportunityCRUDRoutes(t *testing.T) {
	router := newOpportunityRouter(newMemoryOpportunityRepository(), passMiddleware(), passMiddleware(), passMiddleware())

	created := requestJSON(t, router, http.MethodPost, "/opportunities", validOpportunityRequest())
	if created.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, created.Code)
	}

	var opportunity Opportunity
	if err := json.Unmarshal(created.Body.Bytes(), &opportunity); err != nil {
		t.Fatalf("expected opportunity response, got %v", err)
	}
	if opportunity.ID == 0 {
		t.Fatal("expected generated opportunity id")
	}
	if opportunity.Status != StatusPublished {
		t.Fatalf("expected default status %q, got %q", StatusPublished, opportunity.Status)
	}

	var payload map[string]any
	if err := json.Unmarshal(created.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected opportunity JSON object, got %v", err)
	}
	if _, ok := payload["organization_name"]; !ok {
		t.Fatal("expected organization_name response field")
	}
	if _, ok := payload["created_at"]; !ok {
		t.Fatal("expected created_at response field")
	}
	if _, ok := payload["updated_at"]; !ok {
		t.Fatal("expected updated_at response field")
	}
	if _, ok := payload["organizationName"]; ok {
		t.Fatal("expected no organizationName response field")
	}
	if _, ok := payload["createdAt"]; ok {
		t.Fatal("expected no createdAt response field")
	}

	listed := requestJSON(t, router, http.MethodGet, "/opportunities?type=full_time", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listed.Code)
	}

	found := requestJSON(t, router, http.MethodGet, "/opportunities/1", nil)
	if found.Code != http.StatusOK {
		t.Fatalf("expected lookup status %d, got %d", http.StatusOK, found.Code)
	}

	updateRequest := validOpportunityRequest()
	updateRequest.Title = "Senior Go Engineer"
	updated := requestJSON(t, router, http.MethodPut, "/opportunities/1", updateRequest)
	if updated.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d", http.StatusOK, updated.Code)
	}

	deleted := requestJSON(t, router, http.MethodDelete, "/opportunities/1", nil)
	if deleted.Code != http.StatusNoContent {
		t.Fatalf("expected delete status %d, got %d", http.StatusNoContent, deleted.Code)
	}

	missing := requestJSON(t, router, http.MethodGet, "/opportunities/1", nil)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("expected not found status %d, got %d", http.StatusNotFound, missing.Code)
	}
}

func TestOpportunityPublicReadsExcludeUnpublished(t *testing.T) {
	repo := newMemoryOpportunityRepository()
	router := newOpportunityRouter(repo, passMiddleware(), passMiddleware(), passMiddleware())

	published := validOpportunityRequest()
	draft := validOpportunityRequest()
	draft.Status = StatusDraft

	requestJSON(t, router, http.MethodPost, "/opportunities", published)
	requestJSON(t, router, http.MethodPost, "/opportunities", draft)

	listed := requestJSON(t, router, http.MethodGet, "/opportunities", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listed.Code)
	}

	page := decodeOpportunityPage(t, listed)
	if len(page.Data) != 1 {
		t.Fatalf("expected one published opportunity, got %d", len(page.Data))
	}
	if page.Meta.Total != 1 {
		t.Fatalf("expected meta total 1, got %d", page.Meta.Total)
	}

	unpublished := requestJSON(t, router, http.MethodGet, "/opportunities/2", nil)
	if unpublished.Code != http.StatusNotFound {
		t.Fatalf("expected unpublished lookup status %d, got %d", http.StatusNotFound, unpublished.Code)
	}
}

func TestOpportunityAdminCanReadDraft(t *testing.T) {
	repo := newMemoryOpportunityRepository()
	router := newOpportunityRouter(repo, adminOptionalAuthMiddleware(t), passMiddleware(), passMiddleware())

	draft := validOpportunityRequest()
	draft.Status = StatusDraft

	created := requestJSON(t, router, http.MethodPost, "/opportunities", draft)
	if created.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, created.Code)
	}

	found := requestJSON(t, router, http.MethodGet, "/opportunities/1", nil)
	if found.Code != http.StatusOK {
		t.Fatalf("expected admin draft lookup status %d, got %d", http.StatusOK, found.Code)
	}

	var opportunity Opportunity
	if err := json.Unmarshal(found.Body.Bytes(), &opportunity); err != nil {
		t.Fatalf("expected opportunity response, got %v", err)
	}
	if opportunity.Status != StatusDraft {
		t.Fatalf("expected draft status, got %q", opportunity.Status)
	}
}

func TestOpportunityUpdatePreservesDraftStatusWhenStatusOmitted(t *testing.T) {
	repo := newMemoryOpportunityRepository()
	router := newOpportunityRouter(repo, passMiddleware(), passMiddleware(), passMiddleware())

	draft := validOpportunityRequest()
	draft.Status = StatusDraft
	created := requestJSON(t, router, http.MethodPost, "/opportunities", draft)
	if created.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d", http.StatusCreated, created.Code)
	}

	updateRequest := validOpportunityRequest()
	updateRequest.Title = "Updated draft"
	updated := requestJSON(t, router, http.MethodPut, "/opportunities/1", updateRequest)
	if updated.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d", http.StatusOK, updated.Code)
	}

	var opportunity Opportunity
	if err := json.Unmarshal(updated.Body.Bytes(), &opportunity); err != nil {
		t.Fatalf("expected opportunity response, got %v", err)
	}
	if opportunity.Status != StatusDraft {
		t.Fatalf("expected draft status after update without status, got %q", opportunity.Status)
	}
}

func TestOpportunityValidationError(t *testing.T) {
	router := newOpportunityRouter(newMemoryOpportunityRepository(), passMiddleware(), passMiddleware(), passMiddleware())

	request := validOpportunityRequest()
	request.Type = "invalid"

	recorder := requestJSON(t, router, http.MethodPost, "/opportunities", request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected validation status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestOpportunityWriteAuthorizationFailure(t *testing.T) {
	router := newOpportunityRouter(newMemoryOpportunityRepository(), passMiddleware(), abortMiddleware(http.StatusUnauthorized), passMiddleware())

	recorder := requestJSON(t, router, http.MethodPost, "/opportunities", validOpportunityRequest())
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestOpportunityListPaginationDefault(t *testing.T) {
	router := newOpportunityRouter(newMemoryOpportunityRepository(), passMiddleware(), passMiddleware(), passMiddleware())

	for i := 0; i < 3; i++ {
		requestJSON(t, router, http.MethodPost, "/opportunities", validOpportunityRequest())
	}

	listed := requestJSON(t, router, http.MethodGet, "/opportunities", nil)
	if listed.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listed.Code)
	}

	page := decodeOpportunityPage(t, listed)
	if len(page.Data) != 3 {
		t.Fatalf("expected 3 opportunities, got %d", len(page.Data))
	}
	if page.Meta.Page != 1 || page.Meta.PageSize != DefaultPageSize {
		t.Fatalf("expected default page 1 and size %d, got page %d size %d", DefaultPageSize, page.Meta.Page, page.Meta.PageSize)
	}
	if page.Meta.Total != 3 || page.Meta.TotalPages != 1 {
		t.Fatalf("expected total 3 and total_pages 1, got total %d total_pages %d", page.Meta.Total, page.Meta.TotalPages)
	}
}

func TestOpportunityListPaginationCustom(t *testing.T) {
	router := newOpportunityRouter(newMemoryOpportunityRepository(), passMiddleware(), passMiddleware(), passMiddleware())

	for i := 0; i < 5; i++ {
		requestJSON(t, router, http.MethodPost, "/opportunities", validOpportunityRequest())
	}

	listed := requestJSON(t, router, http.MethodGet, "/opportunities?page=2&page_size=2", nil)
	page := decodeOpportunityPage(t, listed)

	if len(page.Data) != 2 {
		t.Fatalf("expected 2 opportunities on page 2, got %d", len(page.Data))
	}
	if page.Meta.Page != 2 || page.Meta.PageSize != 2 || page.Meta.Total != 5 || page.Meta.TotalPages != 3 {
		t.Fatalf("unexpected meta: %+v", page.Meta)
	}
	if page.Data[0].ID != 3 || page.Data[1].ID != 2 {
		t.Fatalf("expected ids [3 2] on page 2, got [%d %d]", page.Data[0].ID, page.Data[1].ID)
	}
}

func TestOpportunityListPaginationClampsInvalidValues(t *testing.T) {
	router := newOpportunityRouter(newMemoryOpportunityRepository(), passMiddleware(), passMiddleware(), passMiddleware())

	for i := 0; i < 3; i++ {
		requestJSON(t, router, http.MethodPost, "/opportunities", validOpportunityRequest())
	}

	oversized := requestJSON(t, router, http.MethodGet, "/opportunities?page=0&page_size=999", nil)
	if oversized.Code != http.StatusOK {
		t.Fatalf("expected status %d for clamped request, got %d", http.StatusOK, oversized.Code)
	}
	oversizedPage := decodeOpportunityPage(t, oversized)
	if oversizedPage.Meta.Page != 1 || oversizedPage.Meta.PageSize != MaxPageSize {
		t.Fatalf("expected clamped page 1 and size %d, got page %d size %d", MaxPageSize, oversizedPage.Meta.Page, oversizedPage.Meta.PageSize)
	}

	nonNumeric := requestJSON(t, router, http.MethodGet, "/opportunities?page=abc&page_size=-5", nil)
	nonNumericPage := decodeOpportunityPage(t, nonNumeric)
	if nonNumericPage.Meta.Page != 1 || nonNumericPage.Meta.PageSize != DefaultPageSize {
		t.Fatalf("expected fallback page 1 and size %d, got page %d size %d", DefaultPageSize, nonNumericPage.Meta.Page, nonNumericPage.Meta.PageSize)
	}

	outOfRange := requestJSON(t, router, http.MethodGet, "/opportunities?page=99999999999999999999&page_size=99999999999999999999", nil)
	if outOfRange.Code != http.StatusOK {
		t.Fatalf("expected status %d for out-of-range values, got %d", http.StatusOK, outOfRange.Code)
	}
	outOfRangePage := decodeOpportunityPage(t, outOfRange)
	if outOfRangePage.Meta.Page != 1 || outOfRangePage.Meta.PageSize != DefaultPageSize {
		t.Fatalf("expected clamped page 1 and default size for out-of-range values, got page %d size %d", outOfRangePage.Meta.Page, outOfRangePage.Meta.PageSize)
	}
	if outOfRangePage.Meta.Total != 3 {
		t.Fatalf("expected total 3, got %d", outOfRangePage.Meta.Total)
	}
}

func decodeOpportunityPage(t *testing.T, recorder *httptest.ResponseRecorder) OpportunityPage {
	t.Helper()

	var page OpportunityPage
	if err := json.Unmarshal(recorder.Body.Bytes(), &page); err != nil {
		t.Fatalf("expected opportunity page response, got %v", err)
	}
	return page
}

func newOpportunityRouter(repo Repository, optionalAuthMiddleware gin.HandlerFunc, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	service := NewService(repo)
	handler := NewHandler(service)
	handler.RegisterRoutes(&router.RouterGroup, optionalAuthMiddleware, authMiddleware, adminMiddleware)
	return router
}

func validOpportunityRequest() SaveOpportunityRequest {
	return SaveOpportunityRequest{
		Title:            "Go Developer",
		Description:      "Build community tools",
		OrganizationName: "Criciuma Devs",
		OrganizationURL:  "https://criciumadevs.local",
		Type:             TypeFullTime,
		WorkMode:         WorkModeRemote,
		Location:         "Criciuma",
		SalaryRange:      "R$ 8k - R$ 12k",
		Seniority:        "senior",
		Skills:           []string{"Go", "PostgreSQL"},
		ContactEmail:     "jobs@criciumadevs.local",
	}
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

func adminOptionalAuthMiddleware(t *testing.T) gin.HandlerFunc {
	t.Helper()

	service := auth.NewService(nil, "test-secret")
	token, err := service.CreateToken(auth.User{ID: 1, Email: "admin@test.local", Role: auth.RoleAdmin})
	if err != nil {
		t.Fatalf("expected admin token, got %v", err)
	}

	middleware := auth.OptionalAuthenticate(service)
	return func(ctx *gin.Context) {
		ctx.Request.Header.Set("Authorization", "Bearer "+token)
		middleware(ctx)
	}
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
