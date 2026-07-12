package opportunities

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/CriciumaDevJobs/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type (
	Handler struct {
		service *Service
	}
)

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup, optionalAuthMiddleware gin.HandlerFunc, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	router.GET("/opportunities", handler.List)
	router.GET("/opportunities/:id", optionalAuthMiddleware, handler.FindByID)

	protected := router.Group("/opportunities", authMiddleware, adminMiddleware)
	protected.POST("", handler.Create)
	protected.PUT("/:id", handler.Update)
	protected.DELETE("/:id", handler.Delete)
}

func (handler *Handler) Create(ctx *gin.Context) {
	var request SaveOpportunityRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity payload"})
		return
	}

	opportunity, err := handler.service.Create(ctx.Request.Context(), request)
	writeOpportunityResult(ctx, opportunity, err, http.StatusCreated)
}

func (handler *Handler) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		pageSize = 0
	}

	filters := ListFilters{
		Type:     ctx.Query("type"),
		WorkMode: ctx.Query("work_mode"),
		Location: ctx.Query("location"),
	}

	result, err := handler.service.List(ctx.Request.Context(), filters, NewPagination(page, pageSize))
	if errors.Is(err, ErrValidation) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity filters"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "opportunity listing failed"})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (handler *Handler) FindByID(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}

	includeUnpublished := false
	if user, ok := auth.UserFromContext(ctx); ok && user.Role == auth.RoleAdmin {
		includeUnpublished = true
	}

	opportunity, err := handler.service.FindByID(ctx.Request.Context(), id, includeUnpublished)
	writeOpportunityResult(ctx, opportunity, err, http.StatusOK)
}

func (handler *Handler) Update(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}

	var request SaveOpportunityRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity payload"})
		return
	}

	opportunity, err := handler.service.Update(ctx.Request.Context(), id, request)
	writeOpportunityResult(ctx, opportunity, err, http.StatusOK)
}

func (handler *Handler) Delete(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}

	if err := handler.service.Delete(ctx.Request.Context(), id); errors.Is(err, ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "opportunity deletion failed"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func idParam(ctx *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func writeOpportunityResult(ctx *gin.Context, opportunity Opportunity, err error, status int) {
	if errors.Is(err, ErrValidation) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity payload"})
		return
	}
	if errors.Is(err, ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "opportunity operation failed"})
		return
	}

	ctx.JSON(status, opportunity)
}
