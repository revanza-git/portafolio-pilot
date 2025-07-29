.PHONY: help run dev test lint migrate seed clean docker-up docker-down generate

# Default target
.DEFAULT_GOAL := help

# Variables
DB_URL := postgresql://defi:defi123@localhost:5432/defi_dashboard?sslmode=disable
DOCKER_COMPOSE := docker-compose
GO := go
GOTEST := $(GO) test
GOLINT := golangci-lint
MIGRATE := migrate
SQLC := sqlc

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
run: ## Run the application
	$(GO) run ./cmd/api

run-worker: ## Run the worker
	$(GO) run ./cmd/worker

dev: ## Run with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

dev-worker: ## Run worker with hot reload
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air -c .air.worker.toml

# Testing
test: ## Run tests
	$(GOTEST) -v -cover ./...

test-coverage: ## Run tests with coverage report
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Linting
lint: ## Run linter
	@which $(GOLINT) > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	$(GOLINT) run

# Database
migrate-up: ## Run database migrations
	$(MIGRATE) -path db/migrations -database "$(DB_URL)" up

migrate-down: ## Rollback database migrations
	$(MIGRATE) -path db/migrations -database "$(DB_URL)" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	$(MIGRATE) create -ext sql -dir db/migrations -seq $(name)

seed: ## Seed the database
	$(GO) run ./scripts/seed.go

# Code generation
generate: ## Generate code (sqlc, oapi-codegen)
	@echo "Generating sqlc code..."
	@which $(SQLC) > /dev/null || (echo "Installing sqlc..." && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
	$(SQLC) generate
	@echo "Generating OpenAPI types..."
	@which oapi-codegen > /dev/null || (echo "Installing oapi-codegen..." && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest)
	oapi-codegen -package api -generate types,server,spec ../spec/openapi.yaml > internal/api/openapi.gen.go

# Docker
docker-up: ## Start all services with Docker Compose
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop all services
	$(DOCKER_COMPOSE) down

docker-build: ## Build Docker images
	$(DOCKER_COMPOSE) build

docker-logs: ## View logs from all services
	$(DOCKER_COMPOSE) logs -f

docker-clean: ## Clean up Docker resources
	$(DOCKER_COMPOSE) down -v
	docker system prune -f

# Build
build: ## Build the application
	$(GO) build -o bin/api ./cmd/api
	$(GO) build -o bin/worker ./cmd/worker

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/api-linux ./cmd/api
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/worker-linux ./cmd/worker

# Dependencies
deps: ## Download dependencies
	$(GO) mod download
	$(GO) mod tidy

# Clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install tools
install-tools: ## Install development tools
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Quick setup
setup: install-tools deps docker-up migrate-up generate seed ## Complete setup for development
	@echo "Setup complete! Run 'make dev' to start the development server"