.DEFAULT_GOAL := help

.PHONY: help run run-go setup-env test test-verbose test-bruno run-testes fmt vet build clean check

APP_NAME ?= server
MAIN_PACKAGE := ./cmd/server
BUILD_DIR ?= bin

help: ## Mostra os comandos disponíveis
	@printf "Uso:\n  make <alvo>\n\nAlvos:\n"
	@printf "  run            Roda a aplicação\n"
	@printf "  setup-env      Cria arquivos .env a partir dos exemplos\n"
	@printf "  test           Roda os testes\n"
	@printf "  test-verbose   Roda os testes com saída verbosa\n"
	@printf "  test-bruno     Roda a coleção Bruno local\n"
	@printf "  fmt            Formata os pacotes Go\n"
	@printf "  vet            Roda go vet\n"
	@printf "  build          Compila o binário em $(BUILD_DIR)/$(APP_NAME)\n"
	@printf "  clean          Remove artefatos de build\n"
	@printf "  check          Roda fmt, vet, test e build\n"

run: ## Roda a aplicação
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi; go run $(MAIN_PACKAGE)

setup-env: ## Cria arquivos .env a partir dos exemplos
	@if [ ! -f .env ] && [ -f .env.example ]; then cp .env.example .env; fi
	@if [ ! -f infra/.env ] && [ -f infra/.env.example ]; then cp infra/.env.example infra/.env; fi

run-go: run ## Alias compatível com o Makefile antigo

test: ## Roda os testes
	go test ./...

test-verbose: ## Roda os testes com saída verbosa
	go test -v ./...

test-bruno: ## Roda a coleção Bruno local
	cd docs/bruno-vagas && bru run --env-file environments/local.json

run-testes: test-verbose ## Alias compatível com o Makefile antigo

fmt: ## Formata os pacotes Go
	go fmt ./...

vet: ## Roda go vet
	go vet ./...

build: ## Compila o binário
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PACKAGE)

clean: ## Remove artefatos de build
	rm -rf $(BUILD_DIR) $(APP_NAME)

check: fmt vet test build ## Roda as validações principais
