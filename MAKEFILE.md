# Neko-Sync Makefile Guide

This document explains how to use the comprehensive Makefile provided with the Neko-Sync project.

## Quick Start

```bash
# Show all available commands
make help

# Setup development environment
make dev-setup

# Create .env file from template
make env

# Start development server with hot reload
make dev

# Build and run the application
make run

# Run tests
make test
```

## Command Categories

### 🚀 Development Commands

- `make dev` - Start development server with hot reload (uses Air)
- `make dev-setup` - Setup development environment (install tools)
- `make env` - Create .env file from template
- `make run` - Build and run the application
- `make install` - Install the application to $GOPATH/bin

### 🔨 Build Commands

- `make build` - Build the application for current platform
- `make build-all` - Build for multiple platforms (Linux, macOS, Windows)
- `make clean` - Clean build artifacts and caches

### 🧪 Testing Commands

- `make test` - Run all tests
- `make test-unit` - Run unit tests only
- `make test-integration` - Run integration tests
- `make test-coverage` - Run tests with coverage report
- `make test-watch` - Run tests in watch mode
- `make bench` - Run benchmarks

### 🔍 Code Quality Commands

- `make fmt` - Format code with go fmt
- `make vet` - Run go vet
- `make lint` - Run golangci-lint
- `make sec` - Run security checks with gosec
- `make check` - Run all checks (fmt, vet, lint, sec, test)

### 🗄️ Database Commands

- `make db-up` - Start PostgreSQL database with Docker
- `make db-down` - Stop and remove database container
- `make db-reset` - Reset database (stop, remove, start fresh)
- `make migrate-up` - Run database migrations up
- `make migrate-down` - Run database migrations down
- `make migrate-create NAME=migration_name` - Create new migration
- `make migrate-version` - Show current migration version
- `make migrate-force VERSION=1` - Force migration version

### 🐳 Docker Commands

- `make docker-build` - Build Docker image
- `make docker-run` - Run application in Docker
- `make docker-up` - Start all services with Docker Compose
- `make docker-down` - Stop all services
- `make docker-logs` - Show Docker Compose logs
- `make docker-clean` - Clean Docker images and containers

### 📚 Documentation Commands

- `make docs` - Start documentation server
- `make docs-gen` - Generate static documentation

### 🔧 Utility Commands

- `make deps` - Download dependencies
- `make deps-update` - Update dependencies
- `make version` - Show version information
- `make info` - Show project information

### 🚢 CI/CD Commands

- `make ci` - Run CI pipeline (deps, check, test-coverage)
- `make release` - Prepare release (clean, check, test, build-all)

## Environment Variables

The Makefile respects several environment variables:

- `DATABASE_URL` - Database connection string
- `VERSION` - Application version (auto-detected from git)
- `DOCKER_TAG` - Docker image tag

## Development Workflow

### Initial Setup

```bash
# 1. Setup development environment
make dev-setup

# 2. Create environment file
make env
# Edit .env file with your configuration

# 3. Start database
make db-up

# 4. Run migrations
make migrate-up

# 5. Start development server
make dev
```

### Daily Development

```bash
# Start development with hot reload
make dev

# Run tests while developing
make test-watch

# Check code quality before commit
make check
```

### Before Commit

```bash
# Run all quality checks
make check

# Run full test suite
make test-coverage
```

### Building for Production

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Prepare release
make release
```

## Docker Development

### Using Docker Compose

```bash
# Start all services (app, database, redis, tools)
make docker-up

# View logs
make docker-logs

# Stop all services
make docker-down
```

### Available Services

- **nekosync** - Main application (port 8080)
- **postgres** - PostgreSQL database (port 5432)
- **redis** - Redis cache (port 6379)
- **adminer** - Database admin interface (port 8081) - with `--profile tools`
- **redis-commander** - Redis admin interface (port 8082) - with `--profile tools`

### Using Admin Tools

```bash
# Start with admin tools
docker-compose --profile tools up -d

# Access admin interfaces
open http://localhost:8081  # Adminer (database)
open http://localhost:8082  # Redis Commander
```

## Database Management

### Creating Migrations

```bash
# Create a new migration
make migrate-create NAME=add_user_table

# This creates two files in ./migrations/
# - 000001_add_user_table.up.sql
# - 000001_add_user_table.down.sql
```

### Running Migrations

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Rollback specific number of migrations
make migrate-down STEPS=2

# Force migration version (use with caution)
make migrate-force VERSION=1
```

## Testing

### Test Types

- **Unit Tests**: Fast tests that test individual functions
- **Integration Tests**: Tests that test multiple components together
- **Benchmarks**: Performance tests

### Running Tests

```bash
# Run all tests
make test

# Run only unit tests (tagged with -short)
make test-unit

# Run integration tests (tagged with integration)
make test-integration

# Run tests with coverage
make test-coverage

# Run tests in watch mode (reruns on file changes)
make test-watch

# Run benchmarks
make bench
```

## Code Quality

### Automated Checks

The `make check` command runs:
1. `go fmt` - Code formatting
2. `go vet` - Static analysis
3. `golangci-lint` - Comprehensive linting
4. `gosec` - Security analysis
5. `go test` - All tests

### Individual Tools

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run linter
make lint

# Run security checks
make sec
```

## Hot Reload Development

The `make dev` command uses [Air](https://github.com/air-verse/air) for hot reload:

- Automatically rebuilds when Go files change
- Excludes test files from triggering rebuilds
- Logs build errors to `build-errors.log`
- Configurable via `.air.toml`

## Troubleshooting

### Common Issues

1. **Database connection failed**
   ```bash
   make db-up
   # Wait a few seconds for database to start
   make migrate-up
   ```

2. **Port already in use**
   ```bash
   # Check what's using the port
   lsof -i :8080
   # Kill the process or change PORT in .env
   ```

3. **Migration failed**
   ```bash
   # Check migration version
   make migrate-version
   # Force to a known good version
   make migrate-force VERSION=1
   ```

4. **Build failed**
   ```bash
   # Clean and rebuild
   make clean build
   ```

5. **Tests failing**
   ```bash
   # Run tests with verbose output
   go test -v ./...
   ```

### Getting Help

- Run `make help` to see all available commands
- Check the `.env.example` file for configuration options
- Review the `docker-compose.yml` for service configuration

## Performance Tips

- Use `make dev` for development (hot reload)
- Use `make test-watch` for TDD workflow
- Use `make docker-up` for full stack testing
- Use `make build-all` only when preparing releases
