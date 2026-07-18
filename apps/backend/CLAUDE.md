# CLAUDE.md — Backend (`apps/backend`)

This file provides guidance to Claude Code (claude.ai/code) when working in the Go backend. For the monorepo big picture, see the root `CLAUDE.md`.

## Scope

Go application, module `nekosync`, requiring **Go 1.26**. All Go commands must be run from this directory (`apps/backend/`) — that is where `go.mod` and the `makefile` live.

## Commands

```bash
# Development with hot reload
air

# Tests
make test                           # all tests
make test-unit                      # short tests only (-short)
make test-coverage                  # generates coverage.html
go test -v ./internal/domain/user/...   # a single package
go test -run TestName ./internal/domain/user/...  # a single test

# Quality
make lint                           # golangci-lint
make fmt                            # gofmt + go mod tidy
make vet                            # go vet
make check                          # fmt + vet + lint + sec + test

# Build (see "Current known breakages" — target is broken)
make build                          # would output ./bin/nekosync
```

## Clean Architecture (DDD, per-aggregate)

Dependencies flow inward: `interfaces → application → domain ← infrastructure`. Every layer imports the module as `nekosync/internal/...`.

The domain is organized as **one self-contained package per aggregate** under `internal/domain/`, each bundling its own entity, repository interface, errors, and (where it has behavior) a service:

- `domain/user/` — `entity.go`, `repository.go` (four interfaces: `Repository`, `DeviceRepository`, `FollowRepository`, `NotificationRepository`), `service.go` (`Service` + `ServiceInterface`, constructed via `NewService(userRepo, deviceRepo, followRepo, notifRepo)`), `errors.go`.
- `domain/party/` — watch-party aggregate, also has a `service.go`.
- `domain/content/`, `domain/social/`, `domain/history/` — entities + repository interfaces (no service yet).
- `domain/shared/` — cross-aggregate types and enums (`types.go`), e.g. `shared.UUID`, `PlatformType`, `NotificationType`. Import this rather than redefining shared types in an aggregate.

Layer responsibilities:

- **`application/usecases/<aggregate>/`** — thin use-case structs, one per operation (`CreateUserUseCase`, `AuthenticateUserUseCase`, …). Each takes the domain `ServiceInterface` in its constructor and exposes `Execute(ctx, req) (resp, error)`. Use cases translate DTOs ↔ domain calls; they hold no business logic.
- **`application/dto/`** — request/response structs used at the HTTP boundary (`CreateUserRequest`, `LoginResponse`, …). Use cases speak DTOs; the domain never sees them.
- **`infrastructure/repositories/`** — PostgreSQL implementations of the domain repository interfaces (`user_repository_impl.go`, `device_repository_impl.go`, `follow_repository_impl.go`, `notification_repository_impl.go`), constructed via `repositories.New*Repository(db)`.
- **`infrastructure/database/`** — `postgres.go`, DB connection setup (`database/sql` + pgx stdlib driver).
- **`interfaces/http/`** — Echo wiring. `server.go`'s `NewHTTPServer(cfg, db)` is the composition root: it builds repositories → domain service → use cases → handler, then registers routes. `handlers/` holds the Echo handlers; `middleware/` holds `auth.go`.

### Wiring pattern (follow this when adding a feature)

`NewHTTPServer` shows the canonical flow. To add an endpoint: define the domain method on the aggregate's `Service` (+ interface), add a use case in `application/usecases/<aggregate>/`, add request/response DTOs, add a handler method, then wire repo → service → use case → handler and register the route in `server.go`. Protected routes go under the `protected` group guarded by `customMiddleware.AuthMiddleware()`.

## Routes

`NewHTTPServer` exposes `GET /health`, public `POST /api/v1/users/register` and `/users/login`, and auth-protected `PUT /users/profile`, `POST /users/follow`, `POST /users/devices`.

## Auth middleware

`interfaces/http/middleware/auth.go` is a **placeholder** — it checks for a `Bearer` token but does not validate a JWT, and sets `user_id` to a hardcoded string. Real JWT validation is pending; don't assume `user_id` from context is trustworthy yet.

## Config & environment

`config.Load()` reads env (via `godotenv`), defaulting `PORT` to `8080` and **calling `log.Fatal` if `DATABASE_URL` is unset**. Copy `.env.example` to `.env` first. Required: `DATABASE_URL`. Expected: `PORT`, `JWT_SECRET`.

## Key dependencies

| Package | Purpose |
|---|---|
| `github.com/labstack/echo/v4` | HTTP framework |
| `github.com/jackc/pgx/v5/stdlib` | PostgreSQL driver (via `database/sql`) |
| `github.com/joho/godotenv` | `.env` loading |
| `golang.org/x/crypto` | Password hashing |

## Current known breakages

- **No `main` package / entrypoint.** There is no `cmd/nekosync` (nor any other `package main`) in the tree, so `make build` — which runs `go build ... ./cmd/nekosync` — fails. `go build ./...` compiles the libraries but produces no binary. `NewHTTPServer` has no caller yet. A `cmd/nekosync/main.go` that runs `config.Load()`, opens the DB, and starts the Echo server is needed to make the app runnable.
- **`go test ./internal/config/...` fails standalone.** `config.go` calls `log.Fatal` (→ `os.Exit`) on missing `DATABASE_URL`, which kills the test binary and reports FAIL. This is a test-design bug, not a regression; the rest of the suite passes.
