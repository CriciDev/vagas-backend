package config

import "os"

type Config struct {
	HTTPPort      string
	DatabaseURL   string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
}

func Load() Config {
	return Config{
		HTTPPort:      env("HTTP_PORT", "8080"),
		DatabaseURL:   env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		JWTSecret:     env("JWT_SECRET", "local-dev-secret"),
		AdminEmail:    env("ADMIN_EMAIL", "admin@criciumadevs.local"),
		AdminPassword: env("ADMIN_PASSWORD", "admin123"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
