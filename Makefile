# DeFi Portfolio Makefile
.PHONY: help dev dev-backend dev-frontend build test clean migrate seed logs

# Default target
help:
	@echo "Available targets:"
	@echo "  make dev          - Start full stack development environment"
	@echo "  make dev-backend  - Start backend services only"
	@echo "  make dev-frontend - Start frontend development server"
	@echo "  make build        - Build all Docker images"
	@echo "  make test         - Run all tests"
	@echo "  make clean        - Stop and remove all containers"
	@echo "  make migrate      - Run database migrations"
	@echo "  make seed         - Seed database with test data"
	@echo "  make logs         - Show logs from all services"

# Development commands
dev:
	@echo "Starting full stack development environment..."
	docker-compose --profile dev up --build

dev-backend:
	@echo "Starting backend services..."
	docker-compose up postgres redis migrate api worker

dev-frontend:
	@echo "Starting frontend development server..."
	cd . && npm run dev

# Build commands
build:
	@echo "Building all Docker images..."
	docker-compose build

build-api-client:
	@echo "Generating TypeScript API client..."
	cd packages/api-client && npm run generate && npm run build

# Testing commands
test: test-backend test-frontend test-e2e

test-backend:
	@echo "Running backend tests..."
	cd backend && go test ./...

test-frontend:
	@echo "Running frontend tests..."
	npm test

test-e2e:
	@echo "Running E2E tests..."
	npx playwright test

test-integration:
	@echo "Running integration tests..."
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Database commands
migrate:
	@echo "Running database migrations..."
	docker-compose run --rm migrate

migrate-down:
	@echo "Rolling back last migration..."
	docker-compose run --rm migrate -path /migrations -database "postgresql://defi:defi123@postgres:5432/defi_dashboard?sslmode=disable" down 1

seed:
	@echo "Seeding database..."
	docker-compose exec api go run scripts/seed.go

# Utility commands
clean:
	@echo "Stopping and removing all containers..."
	docker-compose down -v

logs:
	docker-compose logs -f

logs-api:
	docker-compose logs -f api

logs-worker:
	docker-compose logs -f worker

logs-web:
	docker-compose logs -f web

# Production commands
prod-build:
	@echo "Building for production..."
	docker-compose -f docker-compose.prod.yml build

prod-up:
	@echo "Starting production services..."
	docker-compose -f docker-compose.prod.yml up -d

prod-down:
	@echo "Stopping production services..."
	docker-compose -f docker-compose.prod.yml down

# Development shortcuts
db:
	docker-compose up postgres redis migrate

api:
	docker-compose up postgres redis migrate api

worker:
	docker-compose up postgres redis migrate worker

web:
	docker-compose up web

# Install dependencies
install:
	@echo "Installing all dependencies..."
	npm install
	cd backend && go mod download
	cd packages/api-client && npm install

# Format code
fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...
	npm run format

# Lint code
lint:
	@echo "Linting code..."
	cd backend && golangci-lint run
	npm run lint

# ==================
# Local Development Commands (No Docker for API/Frontend)
# ==================

# Start infrastructure services only (postgres, redis)
infra-up:
	@echo "Starting infrastructure services..."
	docker-compose -f docker-compose.dev.yml up -d postgres redis

# Run database migrations on local postgres
infra-migrate:
	@echo "Running database migrations..."
	docker-compose -f docker-compose.dev.yml run --rm migrate

# Stop infrastructure services
infra-down:
	@echo "Stopping infrastructure services..."
	docker-compose -f docker-compose.dev.yml down

# Start backend API locally (hot reload)
api-local:
	@echo "Starting API server locally on http://localhost:3000..."
	cd backend && go run cmd/api/main.go

# Start backend worker locally
worker-local:
	@echo "Starting worker locally..."
	cd backend && go run cmd/worker/main.go

# Start frontend locally with Vite (hot reload)
web-local:
	@echo "Starting frontend locally on http://localhost:5173..."
	npm run dev

# Full local development setup
local-dev: infra-up infra-migrate
	@echo "===================="
	@echo "Infrastructure ready!"
	@echo "===================="
	@echo ""
	@echo "Now run in separate terminals:"
	@echo "  1. make api-local     # Backend API (http://localhost:3000)"
	@echo "  2. make web-local     # Frontend (http://localhost:5173)"
	@echo "  3. make worker-local  # Worker (optional)"
	@echo ""
	@echo "To stop infrastructure: make infra-down"

# Copy environment variables for local development
env-local:
	@echo "Setting up local development environment..."
	cp .env.local .env
	@echo "Environment ready! Edit .env to add your API keys."