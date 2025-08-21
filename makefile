# =============================================================================
# Neko-Sync Makefile
# =============================================================================

# Application Info
APP_NAME := nekosync
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Directories
BUILD_DIR := ./bin
SRC_DIR := ./cmd/nekosync
MIGRATIONS_DIR := ./migrations
DOCS_DIR := ./docs

# Go build flags
GO_VERSION := $(shell go version | awk '{print $$3}')
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

# Environment
DOCKER_IMAGE := $(APP_NAME)
DOCKER_TAG := $(VERSION)
DATABASE_URL := $(shell grep -E '^DATABASE_URL=' .env 2>/dev/null | cut -d '=' -f2 || echo "postgres://user:pass@localhost:5432/nekosync?sslmode=disable")

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
NC := \033[0m # No Color

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help: ## Show this help message
	@echo "$(CYAN)Neko-Sync Development Commands$(NC)"
	@echo "$(CYAN)================================$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Usage examples:$(NC)"
	@echo "  make dev          # Start development server with hot reload"
	@echo "  make test-watch   # Run tests in watch mode"
	@echo "  make docker-up    # Start all services with Docker Compose"

# =============================================================================
# Development
# =============================================================================

.PHONY: dev
dev: ## Start development server with hot reload
	@echo "$(CYAN)Starting development server...$(NC)"
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "$(YELLOW)Installing air for hot reload...$(NC)"; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

.PHONY: dev-setup
dev-setup: ## Setup development environment
	@echo "$(CYAN)Setting up development environment...$(NC)"
	@go mod download
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)Development environment setup complete!$(NC)"

.PHONY: env
env: ## Create .env file from template
	@if [ ! -f .env ]; then \
		echo "$(CYAN)Creating .env file...$(NC)"; \
		cp .env.example .env 2>/dev/null || \
		echo -e "PORT=8080\nDATABASE_URL=postgres://nekosync:password@localhost:5432/nekosync?sslmode=disable\nJWT_SECRET=your-secret-key\nENVIRONMENT=development" > .env; \
		echo "$(GREEN).env file created! Please update with your settings.$(NC)"; \
	else \
		echo "$(YELLOW).env file already exists$(NC)"; \
	fi

# =============================================================================
# Build & Run
# =============================================================================

.PHONY: build
build: clean ## Build the application
	@echo "$(CYAN)Building $(APP_NAME) $(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

.PHONY: build-all
build-all: ## Build for multiple platforms
	@echo "$(CYAN)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(SRC_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(SRC_DIR)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(SRC_DIR)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(SRC_DIR)
	@echo "$(GREEN)Multi-platform build complete!$(NC)"

.PHONY: run
run: build ## Build and run the application
	@echo "$(CYAN)Running $(APP_NAME)...$(NC)"
	@$(BUILD_DIR)/$(APP_NAME)

.PHONY: install
install: ## Install the application to $GOPATH/bin
	@echo "$(CYAN)Installing $(APP_NAME)...$(NC)"
	@go install $(LDFLAGS) $(SRC_DIR)

# =============================================================================
# Testing
# =============================================================================

.PHONY: test
test: ## Run tests
	@echo "$(CYAN)Running tests...$(NC)"
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

.PHONY: test-watch
test-watch: ## Run tests in watch mode
	@echo "$(CYAN)Running tests in watch mode...$(NC)"
	@if command -v gotestsum > /dev/null 2>&1; then \
		gotestsum --watch ./...; \
	else \
		echo "$(YELLOW)Installing gotestsum for watch mode...$(NC)"; \
		go install gotest.tools/gotestsum@latest; \
		gotestsum --watch ./...; \
	fi

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(CYAN)Running integration tests...$(NC)"
	@go test -v -tags=integration ./tests/integration/...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "$(CYAN)Running unit tests...$(NC)"
	@go test -v -short ./...

.PHONY: bench
bench: ## Run benchmarks
	@echo "$(CYAN)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

# =============================================================================
# Code Quality
# =============================================================================

.PHONY: lint
lint: ## Run linter
	@echo "$(CYAN)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)Installing golangci-lint...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "$(CYAN)Formatting code...$(NC)"
	@go fmt ./...
	@go mod tidy

.PHONY: vet
vet: ## Run go vet
	@echo "$(CYAN)Running go vet...$(NC)"
	@go vet ./...

.PHONY: sec
sec: ## Run security checks
	@echo "$(CYAN)Running security checks...$(NC)"
	@if command -v gosec > /dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)Installing gosec...$(NC)"; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

.PHONY: check
check: fmt vet lint sec test ## Run all checks (format, vet, lint, security, test)

# =============================================================================
# Database
# =============================================================================

.PHONY: db-up
db-up: ## Start database with Docker
	@echo "$(CYAN)Starting database...$(NC)"
	@docker run --name nekosync-postgres -d \
		-e POSTGRES_DB=nekosync \
		-e POSTGRES_USER=nekosync \
		-e POSTGRES_PASSWORD=password \
		-p 5432:5432 \
		postgres:15-alpine || echo "$(YELLOW)Database container may already be running$(NC)"

.PHONY: db-down
db-down: ## Stop and remove database container
	@echo "$(CYAN)Stopping database...$(NC)"
	@docker stop nekosync-postgres 2>/dev/null || true
	@docker rm nekosync-postgres 2>/dev/null || true

.PHONY: db-reset
db-reset: db-down db-up ## Reset database (stop, remove, start fresh)
	@echo "$(GREEN)Database reset complete!$(NC)"

.PHONY: migrate-install
migrate-install: ## Install migrate tool
	@echo "$(CYAN)Installing migrate tool...$(NC)"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "$(RED)Error: NAME is required. Usage: make migrate-create NAME=migration_name$(NC)"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATIONS_DIR)
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
	@echo "$(GREEN)Migration created in $(MIGRATIONS_DIR)$(NC)"

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	@echo "$(CYAN)Running migrations up...$(NC)"
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Run database migrations down (usage: make migrate-down STEPS=1)
	@echo "$(CYAN)Running migrations down...$(NC)"
	@if [ -n "$(STEPS)" ]; then \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down $(STEPS); \
	else \
		migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1; \
	fi

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Error: VERSION is required. Usage: make migrate-force VERSION=1$(NC)"; \
		exit 1; \
	fi
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(VERSION)

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

# =============================================================================
# Docker
# =============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(CYAN)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

.PHONY: docker-run
docker-run: ## Run application in Docker
	@echo "$(CYAN)Running Docker container...$(NC)"
	@docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-up
docker-up: ## Start all services with Docker Compose
	@echo "$(CYAN)Starting services with Docker Compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Services started! API available at http://localhost:8080$(NC)"

.PHONY: docker-down
docker-down: ## Stop all services
	@echo "$(CYAN)Stopping services...$(NC)"
	@docker-compose down

.PHONY: docker-logs
docker-logs: ## Show Docker Compose logs
	@docker-compose logs -f

.PHONY: docker-clean
docker-clean: ## Clean Docker images and containers
	@echo "$(CYAN)Cleaning Docker resources...$(NC)"
	@docker system prune -f
	@docker image prune -f

# =============================================================================
# Documentation
# =============================================================================

.PHONY: docs
docs: ## Generate documentation
	@echo "$(CYAN)Generating documentation...$(NC)"
	@mkdir -p $(DOCS_DIR)
	@if command -v godoc > /dev/null 2>&1; then \
		echo "Documentation server: http://localhost:6060/pkg/nekosync/"; \
		godoc -http=:6060; \
	else \
		echo "$(YELLOW)Installing godoc...$(NC)"; \
		go install golang.org/x/tools/cmd/godoc@latest; \
		echo "Documentation server: http://localhost:6060/pkg/nekosync/"; \
		godoc -http=:6060; \
	fi

.PHONY: docs-gen
docs-gen: ## Generate static documentation
	@echo "$(CYAN)Generating static documentation...$(NC)"
	@mkdir -p $(DOCS_DIR)
	@go doc -all ./... > $(DOCS_DIR)/api.txt
	@echo "$(GREEN)Documentation generated in $(DOCS_DIR)/$(NC)"

# =============================================================================
# Utilities
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts and caches
	@echo "$(CYAN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean -cache -testcache -modcache
	@echo "$(GREEN)Clean complete!$(NC)"

.PHONY: deps
deps: ## Download dependencies
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "$(CYAN)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy

.PHONY: version
version: ## Show version information
	@echo "$(CYAN)Version Information:$(NC)"
	@echo "App Version: $(VERSION)"
	@echo "Go Version:  $(GO_VERSION)"
	@echo "Git Commit:  $(GIT_COMMIT)"
	@echo "Build Time:  $(BUILD_TIME)"

.PHONY: info
info: ## Show project information
	@echo "$(CYAN)Project Information:$(NC)"
	@echo "Name:        $(APP_NAME)"
	@echo "Version:     $(VERSION)"
	@echo "Go Version:  $(GO_VERSION)"
	@echo "Build Dir:   $(BUILD_DIR)"
	@echo "Source Dir:  $(SRC_DIR)"
	@echo "Database:    $(DATABASE_URL)"

# =============================================================================
# CI/CD
# =============================================================================

.PHONY: ci
ci: deps check test-coverage ## Run CI pipeline
	@echo "$(GREEN)CI pipeline completed successfully!$(NC)"

.PHONY: release
release: clean check test build-all ## Prepare release
	@echo "$(GREEN)Release build completed!$(NC)"
	@echo "$(CYAN)Built artifacts:$(NC)"
	@ls -la $(BUILD_DIR)/

.PHONY: tag
tag: ## Create and push a new git tag (usage: make tag VERSION=v1.0.0)
ifndef VERSION
	@echo "$(RED)Error: VERSION is required. Usage: make tag VERSION=v1.0.0$(NC)"
	@exit 1
endif
	@echo "$(CYAN)Creating tag $(VERSION)...$(NC)"
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "$(GREEN)Tag $(VERSION) created and pushed!$(NC)"

.PHONY: release-dry-run
release-dry-run: ## Show what would be released
	@echo "$(CYAN)Release dry run:$(NC)"
	@echo "Current version: $(VERSION)"
	@echo "Git commit: $(GIT_COMMIT)"
	@echo "Build time: $(BUILD_TIME)"
	@echo ""
	@echo "$(CYAN)Would build for platforms:$(NC)"
	@echo "- Linux AMD64"
	@echo "- Linux ARM64" 
	@echo "- macOS AMD64"
	@echo "- macOS ARM64"
	@echo "- Windows AMD64"
	@echo ""
	@echo "$(CYAN)Docker images would be published to:$(NC)"
	@echo "- ghcr.io/$(shell echo $(GIT_REMOTE) | sed 's/.*github.com[:\/]//' | sed 's/\.git$$//' | tr '[:upper:]' '[:lower:]')"

.PHONY: release-notes
release-notes: ## Generate release notes for the current version
	@echo "$(CYAN)Generating release notes...$(NC)"
	@PREV_TAG=$$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "Initial release"); \
	echo "# Release Notes"; \
	echo ""; \
	echo "## Changes since $$PREV_TAG"; \
	echo ""; \
	git log --pretty=format:"- %s (%h)" --no-merges $$PREV_TAG..HEAD || git log --pretty=format:"- %s (%h)" --no-merges

.PHONY: check-release-ready
check-release-ready: ## Check if repository is ready for release
	@echo "$(CYAN)Checking release readiness...$(NC)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)Error: Working directory is not clean$(NC)"; \
		git status --short; \
		exit 1; \
	fi
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "main" ]; then \
		echo "$(YELLOW)Warning: Not on main branch (current: $$(git rev-parse --abbrev-ref HEAD))$(NC)"; \
	fi
	@echo "$(GREEN)Repository is ready for release!$(NC)"

# =============================================================================
# Default target
# =============================================================================

.DEFAULT_GOAL := help

# Make sure intermediate files are not deleted
.PRECIOUS: %/.git