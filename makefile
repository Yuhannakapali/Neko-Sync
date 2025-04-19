# Variables
APP_NAME := nekosync
BUILD_DIR := ./bin
SRC_DIR := ./cmd/nekosync
MIGRATIONS_DIR := ./internal/infra/db/migrations
DATABASE_URL := $(shell grep DATABASE_URL .env | cut -d '=' -f2)

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	@$(BUILD_DIR)/$(APP_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Run migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running migrations up..."
	@docker run --rm -v $(PWD)/$(MIGRATIONS_DIR):/migrations --network host migrate/migrate \
		-path=/migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	@echo "Running migrations down..."
	@docker run --rm -v $(PWD)/$(MIGRATIONS_DIR):/migrations --network host migrate/migrate \
		-path=/migrations -database "$(DATABASE_URL)" down

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Lint the code
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...