package devs

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	router.GET("/developers", handler.List)
	router.GET("/developers/:id", handler.FindByID)

	protected := router.Group("/developers", authMiddleware, adminMiddleware)
	protected.POST("", handler.Create)
	protected.PUT("/:id", handler.Update)
	protected.DELETE("/:id", handler.Delete)
}

func (handler *Handler) Create(ctx *gin.Context) {
	var request SaveDeveloperRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid developer payload"})
		return
	}

	developer, err := handler.service.Create(ctx.Request.Context(), request)
	writeDeveloperResult(ctx, developer, err, http.StatusCreated)
}

func (handler *Handler) List(ctx *gin.Context) {
	developers, err := handler.service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "developer listing failed"})
		return
	}

	ctx.JSON(http.StatusOK, developers)
}

func (handler *Handler) FindByID(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "developer not found"})
		return
	}

	developer, err := handler.service.FindByID(ctx.Request.Context(), id)
	writeDeveloperResult(ctx, developer, err, http.StatusOK)
}

func (handler *Handler) Update(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "developer not found"})
		return
	}

	var request SaveDeveloperRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid developer payload"})
		return
	}

	developer, err := handler.service.Update(ctx.Request.Context(), id, request)
	writeDeveloperResult(ctx, developer, err, http.StatusOK)
}

func (handler *Handler) Delete(ctx *gin.Context) {
	id, ok := idParam(ctx)
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "developer not found"})
		return
	}

	if err := handler.service.Delete(ctx.Request.Context(), id); errors.Is(err, ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "developer not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "developer deletion failed"})
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

func writeDeveloperResult(ctx *gin.Context, developer Developer, err error, status int) {
	if errors.Is(err, ErrValidation) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid developer payload"})
		return
	}
	if errors.Is(err, ErrNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "developer not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "developer operation failed"})
		return
	}

	ctx.JSON(status, developer)
}
