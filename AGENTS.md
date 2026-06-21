# AGENTS.md

## Project Shape
- Go backend module: `github.com/CriciumaDevJobs/backend`; `go.mod` declares Go `1.25.4`.
- Runtime entrypoint is `cmd/server/main.go`; it is intentionally the only application code in the scaffold.
- Keep the scaffold dry: only add directories/packages when there is real behavior; do not create empty `doc.go`, `controller`, `repository`, `usecase`, `sql`, or `infra` placeholders.

## Commands
- Run the app: `make run` or `go run ./cmd/server`.
- Run all tests: `make test` or `go test ./...`.
- Format packages: `make fmt` (`go fmt ./...`).
- Main verification: `make check` runs `fmt`, `vet`, `test`, and `build`.
- Build output is `bin/server`; remove it with `make clean`.

## Current Gotchas
- There is no committed infra, database, codegen, CI, or HTTP router yet; do not add tool-specific commands or docs until the corresponding files exist.
- Prefer `Makefile` and `go.mod` over prose when commands disagree.
