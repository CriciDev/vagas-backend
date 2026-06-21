## 1. Project Setup

- [x] 1.1 Create or switch to git branch `feature/crud-developer` before implementation work.
- [x] 1.2 Add Gin and required backend dependencies to `go.mod` and `go.sum`.
- [x] 1.3 Add application configuration loading from environment variables for HTTP port, database URL, JWT secret, admin email, and admin password.
- [x] 1.4 Initialize Gin in `cmd/server/main.go` with health, auth, and developer route registration.

## 2. Database Bootstrap

- [x] 2.1 Add local database bootstrap SQL for `developers` and `users` tables.
- [x] 2.2 Ensure the local admin user can be created or seeded from Docker environment variables.
- [x] 2.3 Implement database connection setup with startup failure handling.

## 3. Auth And Roles

- [x] 3.1 Implement login request handling in `internal/auth`.
- [x] 3.2 Implement password verification for the seeded admin user.
- [x] 3.3 Implement bearer token creation and validation with role claims.
- [x] 3.4 Implement Gin middleware for authentication and required-role authorization.
- [x] 3.5 Add auth tests for valid login, invalid login, missing token, invalid token, and forbidden role access.

## 4. Developer CRUD

- [x] 4.1 Define the developer model and create/update request validation in `internal/devs`.
- [x] 4.2 Implement PostgreSQL repository methods to create, list, find by ID, update, and delete developers.
- [x] 4.3 Implement usecase logic for developer CRUD with validation and not-found handling.
- [x] 4.4 Implement Gin handlers for developer CRUD endpoints.
- [x] 4.5 Register public read routes and admin-protected write routes.
- [x] 4.6 Add developer tests covering create, list, lookup, update, delete, validation errors, and authorization failures.

## 5. Docker Local Environment

- [x] 5.1 Add a backend Dockerfile for local API execution.
- [x] 5.2 Update `infra/docker-compose.yml` to run API and PostgreSQL together.
- [x] 5.3 Wire environment variables for database connection, JWT secret, and seeded admin credentials.
- [x] 5.4 Ensure Docker startup initializes the local database schema on a fresh volume.
- [x] 5.5 Verify the health endpoint works after Docker Compose startup.

## 6. Documentation And Verification

- [x] 6.1 Update `README.md` with Docker Compose startup, health check, login, and developer CRUD examples.
- [x] 6.2 Run `go test ./...` and fix any failures.
- [x] 6.3 Run the documented Docker Compose flow and verify API health.
- [x] 6.4 Verify protected developer write endpoints require a valid admin token.
