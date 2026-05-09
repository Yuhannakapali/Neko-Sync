# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Neko-Sync is a media streaming platform for anime, manga, movies, and music with watch party functionality. It is an NX monorepo with a Go backend (`apps/backend`) and a Next.js frontend (`apps/web`).

## Commands

All Go commands must be run from `apps/backend/` — that is where `go.mod` lives (module `nekosync`).

### Backend (Go)

```bash
# Development with hot reload (from apps/backend/)
cd apps/backend && air

# Build
make build                          # outputs to ./bin/nekosync

# Run tests
make test                           # all tests
make test-unit                      # short tests only
make test-coverage                  # generates coverage.html

# Run a single test package
cd apps/backend && go test -v ./internal/domain/entities/...

# Lint
make lint                           # golangci-lint
make fmt                            # gofmt + go mod tidy
make vet                            # go vet

# Run all checks
make check                          # fmt + vet + lint + sec + test
```

### Frontend (Next.js via NX)

```bash
npx nx serve web                    # dev server on port 3000
npx nx build web                    # production build
npx nx lint web
npx nx test web
```

### Database & Docker

```bash
# Start PostgreSQL + Redis + app
make docker-up                      # docker-compose up -d

# Database only (no Docker Compose)
make db-up                          # starts postgres:15-alpine container

# Migrations (requires golang-migrate)
make migrate-up
make migrate-down STEPS=1
make migrate-create NAME=add_content_table

# Setup dev environment (installs air, golangci-lint, migrate)
make dev-setup
```

## Architecture

### Repository layout

```
Neko-Sync/
├── apps/
│   ├── backend/                 # Go application (go.mod lives here)
│   │   ├── cmd/nekosync/        # Active entrypoint (main.go)
│   │   └── internal/
│   │       ├── config/          # Env-based config (PORT, DATABASE_URL)
│   │       ├── domain/          # Business rules (see below)
│   │       ├── application/     # Use cases + DTOs
│   │       ├── infrastructure/  # Repository implementations (active)
│   │       ├── infra/           # Older infra layer (db init used by root cmd)
│   │       └── interfaces/http/ # Echo handlers, middleware, server setup
│   └── web/                     # Next.js app (Tailwind, src/app router)
├── cmd/nekosync/                # Root-level entrypoint (uses older infra paths)
├── scripts/                     # init-db.sql, release.sh
├── docker-compose.yml           # App + PostgreSQL + Redis
└── makefile                     # Primary task runner
```

### Backend Clean Architecture layers

Dependencies flow inward: `interfaces → application → domain ← infrastructure`

- **`domain/entities/`** — all domain types (`User`, `Content`, `Anime`, `Manga`, `Movie`, `Music`, `WatchParty`, `DeviceTransfer`, etc.) plus shared enums in `types.go`. This is the canonical entity source.
- **`domain/repositories/`** — repository interfaces (contracts).
- **`domain/services/`** — domain service implementations (e.g., `UserService`).
- **`application/usecases/user/`** — use case structs that orchestrate domain services.
- **`application/dto/`** — request/response types used by HTTP handlers.
- **`infrastructure/repositories/`** — concrete PostgreSQL implementations of the repository interfaces.
- **`interfaces/http/`** — Echo server setup (`server.go`), handlers, and auth middleware. Wire-up happens in `NewHTTPServer`: repository → domain service → use case → handler.

### Known architectural state

There are two overlapping infrastructure paths — `internal/infrastructure/` (active, used by `apps/backend/cmd`) and `internal/infra/` (used by the root-level `cmd/nekosync/main.go`). The makefile build target (`make build`) always builds from `apps/backend/`, so `apps/backend/cmd` is the canonical entrypoint. The root `cmd/` is a legacy artifact.

There are also domain sub-packages (`domain/Anime/`, `domain/content/`, `domain/user/`, etc.) that pre-date the `domain/entities/` consolidation and contain duplicate entity definitions. New code should use `domain/entities/`.

### Auth middleware

`interfaces/http/middleware/auth.go` is a placeholder — it checks for a `Bearer` token but does not validate JWT. `user_id` is set to a hardcoded string. JWT implementation is pending.

### Key dependencies

| Package | Purpose |
|---|---|
| `github.com/labstack/echo/v4` | HTTP framework |
| `github.com/jackc/pgx/v5/stdlib` | PostgreSQL driver (used via `database/sql`) |
| `github.com/joho/godotenv` | `.env` loading |
| `golang.org/x/crypto` | Password hashing |

### Environment

Copy `.env.example` to `.env` before running. Required variables: `DATABASE_URL`. Optional but expected: `PORT` (default `8080`), `JWT_SECRET`.
