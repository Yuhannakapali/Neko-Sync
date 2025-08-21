# Neko-Sync - Restructured Application Architecture

This document describes the restructured application architecture following Clean Architecture principles.

## Architecture Overview

The application now follows a Clean Architecture pattern with clear separation of concerns across different layers:

```
nekosync/
├── cmd/nekosync/           # Application entry point
├── internal/
│   ├── domain/             # Business entities and rules (innermost layer)
│   │   ├── entities/       # Domain entities
│   │   ├── repositories/   # Repository interfaces
│   │   └── services/       # Domain services
│   ├── application/        # Application services & use cases
│   │   ├── usecases/       # Use case implementations
│   │   └── dto/            # Data Transfer Objects
│   ├── infrastructure/     # External dependencies (outermost layer)
│   │   ├── database/       # Database connections
│   │   └── repositories/   # Repository implementations
│   └── interfaces/         # External interfaces
│       └── http/           # HTTP interface layer
│           ├── handlers/   # HTTP request handlers
│           └── middleware/ # HTTP middleware
├── config/                 # Configuration management
└── pkg/                    # Public packages (if any)
```

## Layer Responsibilities

### Domain Layer (`internal/domain/`)
The core business logic layer, completely independent of external concerns.

- **Entities** (`entities/`): Business entities with their business rules
- **Repository Interfaces** (`repositories/`): Define contracts for data access
- **Domain Services** (`services/`): Complex business logic that doesn't belong to a single entity

### Application Layer (`internal/application/`)
Orchestrates use cases and manages application flow.

- **Use Cases** (`usecases/`): Application-specific business rules
- **DTOs** (`dto/`): Data transfer objects for API contracts

### Infrastructure Layer (`internal/infrastructure/`)
Handles external dependencies and technical details.

- **Database** (`database/`): Database connections and setup
- **Repository Implementations** (`repositories/`): Concrete implementations of repository interfaces

### Interface Layer (`internal/interfaces/`)
Handles external communication protocols.

- **HTTP** (`http/`): REST API implementation
  - **Handlers**: Process HTTP requests and responses
  - **Middleware**: Authentication, logging, CORS, etc.

## Key Benefits

1. **Testability**: Easy to unit test business logic in isolation
2. **Maintainability**: Clear boundaries and responsibilities
3. **Flexibility**: Can easily swap implementations (e.g., database drivers)
4. **Scalability**: Easy to add new features without affecting existing code

## Dependency Direction

Dependencies flow inward:
- Interface Layer → Application Layer → Domain Layer
- Infrastructure Layer → Domain Layer (via interfaces)
- No layer depends on layers outside of it

## API Endpoints

### User Management
- `POST /api/v1/users/register` - User registration
- `POST /api/v1/users/login` - User authentication
- `PUT /api/v1/users/profile` - Update user profile (protected)
- `POST /api/v1/users/follow` - Follow a user (protected)
- `POST /api/v1/users/devices` - Register a device (protected)

### Authentication
Protected routes require an `Authorization: Bearer <token>` header.

## Running the Application

```bash
# Build the application
go build ./cmd/nekosync

# Run the application
./nekosync
```

## Next Steps

1. **Database Migrations**: Create database schema migration files
2. **JWT Authentication**: Implement proper JWT token generation and validation
3. **Validation**: Add request validation using a validation library
4. **Logging**: Implement structured logging
5. **Testing**: Add unit and integration tests
6. **Content Management**: Implement content-related use cases and handlers
7. **Watch Party**: Implement real-time watch party functionality with WebSockets

## Development Notes

- The current authentication middleware is a placeholder - implement proper JWT validation
- Database schema needs to be created to match the entity definitions
- Add proper error handling and logging throughout the application
- Consider adding configuration for different environments (dev, staging, prod)
