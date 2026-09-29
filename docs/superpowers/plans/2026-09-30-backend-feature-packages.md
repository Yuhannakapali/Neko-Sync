# Backend: layered → feature packages

Status: implemented 2026-09-30 (see Deviations) · Scope: `apps/backend` only (Instance decision deferred)

## Goal

A change to one feature touches one folder. No logic changes in this move —
files relocate, packages rename, tests stay green (33/33 before and after).

## Target layout

```
apps/backend/
├── cmd/nekosync/main.go
└── internal/
    ├── app/server.go              # composition root + route registration
    ├── platform/
    │   ├── config/                # from internal/config
    │   ├── postgres/              # from infrastructure/database (+ WithTx later)
    │   ├── auth/                  # from infrastructure/auth (JWTManager)
    │   ├── httpx/                 # from interfaces/http/middleware
    │   └── id/                    # from domain/shared (UUID, BaseEntity)
    ├── user/                      # entity, errors, service, store, postgres, http, dto
    ├── party/
    ├── work/
    ├── reference/
    ├── social/
    └── history/
```

## Why this shape

- Go's own guidance for server projects (go.dev/doc/modules/layout) is `cmd/`
  plus `internal/` packages named for what they provide, not for their layer.
  "Package names" (go.dev/blog/package-names) warns against catch-all names
  like `shared`, `handlers`, `repositories`, `dto`.
- The current graph is already a clean DAG with features as leaves:

  ```
  party → user, work           (only cross-feature edge)
  every domain pkg → shared
  repositories → user          handlers → user
  interfaces/http → everything (composition root)
  ```

  Nothing imports "upward", so collapsing each feature's layers into one
  package cannot create an import cycle. `party → user/work` stays legal.

- Interfaces with one implementation (`user.Repository`, `DeviceRepository`,
  `FollowRepository`, `NotificationRepository`) merge into one `Store`
  interface declared where it is consumed (`user/store.go`). It stays an
  interface so `service_test.go` keeps its fake.

## File map

| From                                             | To                                                         |
| ------------------------------------------------ | ---------------------------------------------------------- |
| `internal/config/*`                              | `internal/platform/config/*`                               |
| `infrastructure/database/postgres.go`            | `platform/postgres/postgres.go`                            |
| `infrastructure/auth/*`                          | `platform/auth/*`                                          |
| `interfaces/http/middleware/*`                   | `platform/httpx/*`                                         |
| `domain/shared/*`                                | `platform/id/*` (`shared.UUID` → `id.UUID`)                |
| `domain/user/{entity,errors,service}.go` + tests | `user/{user,errors,service}.go`                            |
| `domain/user/repository.go`                      | `user/store.go` (one `Store` interface)                    |
| `infrastructure/repositories/*_impl.go`          | `user/postgres.go` (one `pgStore` type)                    |
| `interfaces/http/handlers/user_handler.go`       | `user/http.go` (+ `Routes(public, protected *echo.Group)`) |
| `interfaces/http/handlers/user_dto.go`           | `user/dto.go`                                              |
| `interfaces/http/server.go`                      | `app/server.go`                                            |
| `domain/{party,work,reference,social,history}/*` | `{party,work,...}/*`                                       |

## Steps (one commit each, `make test` green after every one)

1. **platform/** — `git mv` config, database, auth, middleware, shared; rewrite
   imports (`go list` must show no `domain/shared` left). Pure moves.
2. **user/** — move entity/errors/service/tests; merge the four repository
   interfaces into `Store` and the four impls into `pgStore` (method names
   kept, e.g. `CreateDevice`, `ReplaceActiveDevice`); fold handler + DTOs in
   and add `Routes`. Update `service_test.go` fake to `Store`.
3. **app/** — move `server.go`; it now calls `user.Routes(api, protected)`.
   `main.go` imports `internal/app`.
4. **party, work, reference, social, history** — `git mv domain/X → X`.
   `party` imports `user` and `work` as before.
5. **Cleanup + docs** — delete empty `domain/`, `infrastructure/`,
   `interfaces/`; update `apps/backend/CLAUDE.md`, root `CLAUDE.md`
   ("Backend layering" table), and the repo-mapping section of
   `docs/superpowers/specs/2026-07-18-hub-instance-architecture.md`.
   `docs/CODEBASE-STATUS.md` is a dated audit — leave it, add a note.

Mechanics for each move: `git mv` (keeps history), fix the `package` line,
rewrite import paths with `sed` over `*.go`, then `go build ./... && go vet
./...`. Qualifiers inside a merged package (`userDomain.X` in the handler)
become unqualified. Verify with `go list -deps ./... | grep internal/` that
no old path remains.

## Things that reference current paths

- `apps/backend/CLAUDE.md` (layering section, test examples at lines 19–20)
- root `CLAUDE.md` "Backend layering" table
- `docs/superpowers/specs/2026-07-18-hub-instance-architecture.md` repo mapping
- No CI, makefile, `.golangci.yml`, or Dockerfile path depends on the internal
  layout — they build `./cmd/nekosync` and `./...`.

## Risks

- **Instance divergence**: `apps/instance` stays layered
  (`domain/application/infrastructure/interfaces`). Acceptable for now; the
  same plan applies there later if wanted.
- **`user` gets large** (~600 lines across its files). Fine; split into
  `user/` + `follow/` + `device/` later only if it actually hurts.
- **Name stutter**: avoid `user.UserService`; use `user.Service`,
  `user.Store`, `user.NewHandler`.
- Pre-existing bugs found in the audit (plaintext party passwords, ID format
  drift, room-code bias) are **not** fixed in this move — keep it mechanical.

## Done when

`go build`, `go vet`, full `make test` (33 pass, 0 fail) and `make build` all
succeed; `find internal -type d` matches the target layout; docs updated.

## Deviations (as implemented)

- `domain/shared` became `platform/entity`, not `platform/id`: a package named
  `id` collides with the many `id` parameters in repository code.
- The four user repository interfaces and their four Postgres types were
  **kept as-is** (`repository.go`, `postgres_{user,device,follow,notification}.go`)
  instead of merging into one `Store`. Merging needs method renames (`Create`
  exists three times), which is a logic change, not a move. Do it separately if
  wanted.
- `interfaces/http.NewHTTPServer` is now `app.NewServer`.
