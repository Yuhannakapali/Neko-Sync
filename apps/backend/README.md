# Neko-Sync Backend

A high-performance REST API backend for the Neko-Sync content synchronization platform, built with Go and Clean Architecture principles.

## 🏗️ Architecture

This backend follows **Clean Architecture** patterns with clear separation of concerns:

```
apps/backend/
├── cmd/nekosync/           # Application entry point
├── internal/
│   ├── application/        # Application layer (use cases, DTOs)
│   │   ├── dto/           # Data Transfer Objects
│   │   └── usecases/      # Business use cases
│   ├── domain/            # Domain layer (entities, services)
│   │   ├── entities/      # Core business entities
│   │   ├── services/      # Domain services
│   │   └── repositories/  # Repository interfaces
│   ├── infrastructure/    # Infrastructure layer
│   │   ├── database/      # Database connections
│   │   └── repositories/  # Repository implementations
│   ├── interfaces/        # Interface layer (HTTP handlers)
│   │   └── http/         # HTTP interface (handlers, routes, middleware)
│   ├── config/           # Configuration management
│   └── infra/            # Infrastructure utilities
│       ├── db/           # Database initialization
│       └── persistence/  # Data persistence
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.24+**
- **PostgreSQL 15+**
- **Make** (for build automation)

### Installation

1. **Clone and navigate to backend:**
   ```bash
   cd apps/backend
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Set up database:**
   ```bash
   # Create PostgreSQL database
   createdb nekosync
   
   # Update DATABASE_URL in .env
   DATABASE_URL=postgres://username:password@localhost:5432/nekosync?sslmode=disable
   ```

5. **Build and run:**
   ```bash
   # From project root
   make build
   make run
   
   # OR directly
   go run cmd/nekosync/main.go
   ```

## 🛠️ Development

### Available Commands

```bash
# Build
make build              # Build the application
make build-all          # Build for multiple platforms

# Testing
make test               # Run tests
make test-coverage      # Run tests with coverage
make test-watch         # Run tests in watch mode

# Code Quality
make check              # Run linting and checks
make fmt                # Format code
make vet                # Run go vet

# Development
make dev                # Start with hot reload (Air)
make deps               # Install dependencies
```

### Project Structure

#### Domain Layer (`internal/domain/`)

**Core business entities and logic:**

- **`entities/`** - Core business entities (User, Content, Party, etc.)
- **`services/`** - Domain services containing business logic
- **`repositories/`** - Repository interfaces (contracts)

**Key Entities:**
- `User` - User management and authentication
- `Content` - Media content (anime, manga, movies, music)
- `WatchParty` - Synchronized viewing sessions
- `History` - User consumption history
- `Social` - Discussions, comments, reviews

#### Application Layer (`internal/application/`)

**Use cases and application services:**

- **`usecases/`** - Application use cases (orchestrate domain services)
- **`dto/`** - Data Transfer Objects for API communication

#### Infrastructure Layer (`internal/infrastructure/`)

**External concerns and implementations:**

- **`database/`** - Database connection and configuration
- **`repositories/`** - Concrete repository implementations

#### Interface Layer (`internal/interfaces/`)

**External interfaces (HTTP, gRPC, etc.):**

- **`http/`** - HTTP handlers, routes, middleware
  - `handlers/` - HTTP request handlers
  - `routes/` - Route definitions
  - `middleware/` - HTTP middleware (auth, logging, etc.)

### Configuration

The application uses environment variables for configuration:

```env
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/nekosync?sslmode=disable

# Authentication
JWT_SECRET=your-jwt-secret-key

# External APIs (optional)
TMDB_API_KEY=your-tmdb-api-key
MAL_CLIENT_ID=your-mal-client-id
```

## 📊 API Documentation

### Authentication

The API uses JWT tokens for authentication:

```bash
# Register
POST /api/users
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password"
}

# Login
POST /api/auth/login
{
  "email": "john@example.com",
  "password": "secure_password"
}
```

### Core Endpoints

#### Users
- `POST /api/users` - Register user
- `POST /api/auth/login` - Login
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update profile
- `POST /api/users/follow` - Follow user

#### Content
- `GET /api/content` - List content
- `GET /api/content/{id}` - Get content details
- `POST /api/content/{id}/favorite` - Add to favorites
- `GET /api/content/{id}/history` - Get watch/read history

#### Watch Parties
- `POST /api/parties` - Create watch party
- `GET /api/parties/{id}` - Get party details
- `POST /api/parties/{id}/join` - Join party
- `POST /api/parties/{id}/state` - Update playback state

## 🧪 Testing

### Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Specific package
go test ./internal/domain/user -v

# Integration tests
make test-integration
```

### Test Structure

```
internal/
├── domain/
│   ├── user/
│   │   ├── user.go
│   │   ├── user_test.go          # Unit tests
│   │   └── user_integration_test.go  # Integration tests
│   └── ...
└── tests/
    ├── integration/              # Integration test suites
    ├── fixtures/                 # Test data
    └── helpers/                  # Test utilities
```

## 🐳 Docker

### Development

```bash
# Build image
docker build -t nekosync-backend .

# Run with docker-compose
docker-compose up backend

# Development with hot reload
docker-compose -f docker-compose.dev.yml up
```

### Production

```bash
# Multi-stage build
docker build -t nekosync-backend:prod --target production .

# Run
docker run -p 8080:8080 --env-file .env nekosync-backend:prod
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `ENVIRONMENT` | Environment (dev/prod/test) | `development` | No |
| `LOG_LEVEL` | Logging level | `info` | No |

### Database Schema

The application uses PostgreSQL with the following key tables:

- `users` - User accounts and profiles
- `content` - Media content metadata
- `watch_parties` - Synchronized viewing sessions
- `user_follows` - User relationship graph
- `watch_history` - User consumption tracking
- `notifications` - User notifications

## 🏗️ Development Guidelines

### Code Style

- Follow Go idioms and conventions
- Use `gofmt` for formatting
- Pass `golangci-lint` checks
- Maintain test coverage > 80%

### Architecture Principles

1. **Dependency Inversion** - Depend on abstractions, not concretions
2. **Single Responsibility** - Each module has one reason to change
3. **Interface Segregation** - Small, focused interfaces
4. **Domain-Driven Design** - Model the business domain

### Adding New Features

1. **Define Domain Entity** (if needed)
2. **Create Repository Interface** 
3. **Implement Use Cases**
4. **Add HTTP Handlers**
5. **Write Tests**
6. **Update Documentation**

## 📈 Performance

### Optimization Features

- **Connection Pooling** - PostgreSQL connection pool
- **Middleware Caching** - HTTP response caching
- **Database Indexing** - Optimized queries
- **Graceful Shutdown** - Clean resource cleanup

### Monitoring

- **Health Checks** - `/health` endpoint
- **Metrics** - Prometheus-compatible metrics
- **Logging** - Structured JSON logging
- **Tracing** - Distributed tracing support

## 🔒 Security

### Security Features

- **JWT Authentication** - Stateless authentication
- **Password Hashing** - bcrypt with salt
- **SQL Injection Protection** - Parameterized queries
- **CORS Protection** - Configurable CORS middleware
- **Rate Limiting** - API rate limiting

### Security Headers

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Strict-Transport-Security`

## 📚 Additional Resources

- [API Documentation](./docs/api.md)
- [Database Schema](./docs/schema.md)
- [Deployment Guide](./docs/deployment.md)
- [Contributing Guidelines](./docs/contributing.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run quality checks: `make check`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
