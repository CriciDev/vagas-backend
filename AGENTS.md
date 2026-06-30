# AGENTS.md

## Project Shape
- Go backend module: `github.com/CriciumaDevJobs/backend`; `go.mod` declares Go `1.25.4`.
- Runtime entrypoint is `cmd/server/main.go`.
- The backend has real behavior for health checks, auth, developer CRUD, database setup, and Docker local environment.
- Keep the project dry: only add directories/packages when there is real behavior; do not create empty `doc.go`, `controller`, `repository`, `usecase`, `sql`, or `infra` placeholders.

## Commands
- Run the app: `make run` or `go run ./cmd/server`.
- Run all tests: `make test` or `go test ./...`.
- Format packages: `make fmt` (`go fmt ./...`).
- Main verification: `make check` runs `fmt`, `vet`, `test`, and `build`.
- Build output is `bin/server`; remove it with `make clean`.

## Current Gotchas
- There is committed Gin, PostgreSQL, auth, Docker, and developer CRUD behavior.
- Prefer `Makefile` and `go.mod` over prose when commands disagree.

## Code Style
- Não use comentários em blocos de código a menos que isso seja expressamente solicitado pelo usuário.
- Agrupe declarações relacionadas de `var` e `type` usando blocos `var (...)` e `type (...)`, como em `internal/devs/devs.go`, para manter a leitura mais limpa.
- Em structs Go, mantenha nomes de campos no padrão da linguagem (`PascalCase` para exportados e `camelCase` para não exportados) e use tags `json` em `snake_case`, especialmente para campos compostos como `organization_name`, `work_mode`, `created_at` e `updated_at`.
