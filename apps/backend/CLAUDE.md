# CLAUDE.md — Backend (`apps/backend`)

This file provides guidance to Claude Code (claude.ai/code) when working in the Go backend. For the monorepo big picture, see the root `CLAUDE.md`.

## Scope

Go application, module `nekosync`, requiring **Go 1.26**. Run Go commands from this directory (`apps/backend/`), where `go.mod` lives; the `makefile` is at the repo root and `cd`s here.

## Commands

```bash
# Development with hot reload
air

# Tests
make test                           # all tests (Postgres tests skip)
make test-integration               # all tests incl. Postgres, in a throwaway container on :55432
make test-unit                      # short tests only (-short)
make test-coverage                  # generates coverage.html
go test -v ./internal/user/...   # a single package
go test -run TestName ./internal/user/...  # a single test

# Quality
make lint                           # golangci-lint
make fmt                            # gofmt + go mod tidy
make vet                            # go vet
make check                          # fmt + vet + lint + sec + test

# Build
make build                          # outputs ./bin/nekosync (from ./cmd/nekosync)
```

## Layout — one package per feature

```
cmd/nekosync/main.go        load config → open db → app.NewServer → start
internal/
  app/server.go             composition root: builds each feature, registers its routes
  platform/                 shared plumbing, no business logic
    config/                 config.Load() → (*Config, error)
    postgres/               postgres.Init(cfg) → (*sql.DB, error)
    auth/                   JWTManager: Issue / Verify HS256 tokens (subject = user ID)
    httpx/                  AuthMiddleware(verifier): sets "user_id" from a verified token
    entity/                 entity.UUID, entity.NewUUID(), entity.BaseEntity
    postgres/pgtest/        pgtest.Open(t): integration-test DB (skips without TEST_DATABASE_URL)
  progress/                 Progress{Fraction 0–1, Locator}: the one progress model
  user/                     the only fully wired feature (see below)
  party/                    entity, errors, repository iface, service — not wired to HTTP
  work/, reference/         Hub catalog + ContentReference registry, with Postgres repos (no HTTP yet)
  social/, history/         entity + repository iface only (history is Entry{WorkID, ChildID, Progress})
```

A feature package owns everything about that feature. `user/` contains:

- `entity.go`, `errors.go` — domain types and errors
- `service.go` — business rules (`Service`, built by `NewService(userRepo, deviceRepo, followRepo, notifRepo)`)
- `repository.go` — the four storage interfaces the service needs
- `postgres_{user,device,follow,notification}.go` — their PostgreSQL implementations (`New*Repository(db)`)
- `http.go` — Echo handlers (`Handler`, `NewHandler(svc, tokens)`) and `Routes(public, protected)`
- `dto.go` — request/response JSON shapes

Dependency rules: features may import `platform/*` and, where the domain needs it, another feature (`party` imports `user` and `work`). `platform/*` never imports a feature. Only `app` knows about every feature.

### Adding an endpoint

All in the feature's folder: add the method to `service.go` (and a query to the matching `postgres_*.go` + `repository.go` if needed), add DTOs to `dto.go`, add a handler to `http.go`, and register it in that package's `Routes`. `app/server.go` changes only when adding a whole new feature. Routes under `protected` require a valid JWT; read the caller with `c.Get("user_id").(string)`. Do not add a use-case layer back unless an operation genuinely spans several features.

## Routes

`GET /health`, public `POST /api/v1/users/register` and `/users/login`, and JWT-protected `PUT /users/profile`, `POST /users/follow`, `POST /users/devices`.

## Config & environment

`config.Load()` reads env (via `godotenv`) and returns an error instead of exiting. Required: `DATABASE_URL`, `JWT_SECRET`. Optional: `PORT` (default `8080`), `JWT_EXPIRY` (Go duration, default `24h`). Copy the repo-root `.env.example` to `.env` first.

## Key dependencies

| Package                          | Purpose                                |
| -------------------------------- | -------------------------------------- |
| `github.com/labstack/echo/v4`    | HTTP framework                         |
| `github.com/jackc/pgx/v5/stdlib` | PostgreSQL driver (via `database/sql`) |
| `github.com/joho/godotenv`       | `.env` loading                         |
| `golang.org/x/crypto`            | Password hashing                       |
| `github.com/golang-jwt/jwt/v5`   | Access tokens                          |

## Database & migrations

Schema lives in versioned golang-migrate files under the repo-root `migrations/`
(`000001_users`, `000002_catalog`). `make migrate-up` applies them; the compose
`migrate` service does the same. `scripts/init-db.sql` no longer creates schema.
`000001` is guarded so it also applies cleanly to databases created by the old
`init-db.sql`. New IDs must come from `entity.NewUUID()` — the canonical form
Postgres returns — or they change shape between write and read.

Postgres-backed tests call `pgtest.Open(t)`, which truncates `users` and `works`
(cascading). Only ever point `TEST_DATABASE_URL` at a throwaway database.

## Current known gaps

- `party`, `social` and `history` have no tables or repository implementations yet.
- `work` and `reference` have repositories but no HTTP endpoints.
- Handlers return raw error strings to clients and do no input validation (`// TODO: Add validation`).
