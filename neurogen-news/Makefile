.PHONY: all build run dev test lint clean docker-build docker-up docker-down migrate frontend backend help

# Variables
APP_NAME := neurogen-news
GO := go
NPM := npm
DOCKER_COMPOSE := docker compose

# Colors
GREEN := \033[0;32m
NC := \033[0m

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

all: build ## Build everything

# Development
dev: ## Start development environment
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml up -d postgres redis meilisearch
	@echo "Starting backend..."
	cd backend && $(GO) run ./cmd/server &
	@echo "Starting frontend..."
	cd frontend && $(NPM) run dev

dev-docker: ## Start development with Docker
	$(DOCKER_COMPOSE) -f docker-compose.yml -f docker-compose.dev.yml up

# Frontend
frontend-install: ## Install frontend dependencies
	cd frontend && $(NPM) ci

frontend-dev: ## Start frontend dev server
	cd frontend && $(NPM) run dev

frontend-build: ## Build frontend for production
	cd frontend && $(NPM) run build

frontend-lint: ## Lint frontend code
	cd frontend && $(NPM) run lint

# Backend
backend-dev: ## Start backend dev server
	cd backend && $(GO) run ./cmd/server

backend-build: ## Build backend binary
	cd backend && CGO_ENABLED=0 $(GO) build -ldflags="-w -s" -o ../bin/$(APP_NAME) ./cmd/server

backend-test: ## Run backend tests
	cd backend && $(GO) test -v ./...

backend-lint: ## Lint backend code
	cd backend && golangci-lint run

# Build
build: frontend-build backend-build ## Build everything for production
	@echo "Build complete! Binary at ./bin/$(APP_NAME)"

# Docker
docker-build: ## Build Docker image
	docker build -t $(APP_NAME):latest .

docker-up: ## Start Docker containers
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-clean: ## Clean Docker volumes
	$(DOCKER_COMPOSE) down -v

# Database
migrate: ## Run database migrations
	@echo "Running migrations..."
	psql $(DATABASE_URL) -f backend/migrations/001_init.sql

migrate-down: ## Rollback migrations (careful!)
	@echo "This would rollback migrations - implement as needed"

db-seed: ## Seed database with test data
	cd backend && $(GO) run ./cmd/seed

# Testing
test: backend-test ## Run all tests

test-coverage: ## Run tests with coverage
	cd backend && $(GO) test -coverprofile=coverage.out ./...
	cd backend && $(GO) tool cover -html=coverage.out -o coverage.html

# Linting
lint: frontend-lint backend-lint ## Lint all code

# Cleaning
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf frontend/dist/
	rm -rf backend/vendor/
	cd frontend && rm -rf node_modules/

# Production
deploy: ## Deploy to production
	@echo "Building and pushing Docker image..."
	docker build -t $(APP_NAME):latest .
	docker tag $(APP_NAME):latest your-registry/$(APP_NAME):latest
	docker push your-registry/$(APP_NAME):latest
	@echo "Deploy to your infrastructure..."

# Utilities
generate-key: ## Generate a random JWT secret
	@openssl rand -base64 32

check-deps: ## Check for outdated dependencies
	cd frontend && $(NPM) outdated || true
	cd backend && $(GO) list -u -m all

