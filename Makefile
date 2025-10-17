.PHONY: help setup test-api-key test lint fmt clean docker-up docker-down docker-logs

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Setup development environment
	@echo "Setting up development environment..."
	cp backend/.env.example backend/.env
	cp frontend/.env.local.example frontend/.env.local
	@echo ""
	@echo "✅ Environment files created"
	@echo ""
	@echo "📝 Next steps:"
	@echo "  1. Edit backend/.env and add your GEMINI_API_KEY"
	@echo "  2. See docs/SETUP_API_KEY.md for how to get API key using gcloud"
	@echo "  3. Run 'make test-api-key' to verify the API key"
	@echo ""

test-api-key: ## Test Gemini API key connection
	@bash scripts/test-api-key.sh

# Backend targets
backend-test: ## Run backend tests
	cd backend && go test -v -race -cover ./...

backend-lint: ## Run backend linter
	cd backend && golangci-lint run --fix ./...

backend-fmt: ## Format backend code
	cd backend && golangci-lint fmt ./...

backend-tidy: ## Tidy backend dependencies
	cd backend && go mod tidy

backend-check: backend-fmt backend-lint backend-test ## Run all backend checks (fmt, lint, test)

# Frontend targets
frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-test: ## Run frontend tests
	cd frontend && npm run test

frontend-lint: ## Run frontend linter
	cd frontend && npm run lint

frontend-fmt: ## Format frontend code
	cd frontend && npm run format

frontend-type-check: ## Run TypeScript type checking
	cd frontend && npm run type-check

frontend-check: frontend-fmt frontend-lint frontend-type-check frontend-test ## Run all frontend checks

# Docker targets
docker-up: ## Start all services with docker-compose
	docker-compose up --build

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show docker logs
	docker-compose logs -f

docker-backend-logs: ## Show backend logs only
	docker-compose logs -f backend

docker-frontend-logs: ## Show frontend logs only
	docker-compose logs -f frontend

# Development workflow (TDD style)
test: backend-test frontend-test ## Run all tests

lint: backend-lint frontend-lint ## Run all linters

fmt: backend-fmt frontend-fmt ## Format all code

check: fmt lint test ## Run all checks (fmt, lint, test) - TDD workflow

clean: ## Clean temporary files
	rm -rf backend/tmp
	rm -rf frontend/.next
	rm -rf frontend/node_modules

.DEFAULT_GOAL := help
