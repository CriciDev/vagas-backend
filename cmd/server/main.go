package main

import (
	"context"
	"log"
	"time"

	"github.com/CriciumaDevJobs/backend/internal/auth"
	"github.com/CriciumaDevJobs/backend/internal/config"
	"github.com/CriciumaDevJobs/backend/internal/database"
	"github.com/CriciumaDevJobs/backend/internal/devs"
	"github.com/CriciumaDevJobs/backend/internal/health"
	"github.com/CriciumaDevJobs/backend/internal/middleware"
	"github.com/CriciumaDevJobs/backend/internal/opportunities"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	if err := database.EnsureSchema(ctx, db); err != nil {
		log.Fatalf("database schema setup failed: %v", err)
	}

	if err := database.SeedAdmin(ctx, db, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("admin seed failed: %v", err)
	}

	router := gin.New()
	router.Use(middleware.RequestID(), gin.LoggerWithFormatter(middleware.AccessLogFormatter), gin.Recovery())
	health.RegisterRoutes(router)

	api := router.Group("/api")
	authRepo := auth.NewPostgresRepository(db)
	authService := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(api)

	devRepo := devs.NewPostgresRepository(db)
	devService := devs.NewService(devRepo)
	devHandler := devs.NewHandler(devService)
	devHandler.RegisterRoutes(api, auth.Authenticate(authService), auth.RequireRole(auth.RoleAdmin))

	opportunityRepo := opportunities.NewPostgresRepository(db)
	opportunityService := opportunities.NewService(opportunityRepo)
	opportunityHandler := opportunities.NewHandler(opportunityService)
	opportunityHandler.RegisterRoutes(api, auth.OptionalAuthenticate(authService), auth.Authenticate(authService), auth.RequireRole(auth.RoleAdmin))

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
