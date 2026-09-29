# CLAUDE.md

Guidance for Claude Code (claude.ai/code) working in this repository.

> **Read this first, and believe the "Reality check" section over any other doc in
> the repo.** `readme.md` and `ARCHITECTURE.md` are both stale (see Doc drift).
> Last verified against the tree: 2026-09-10, commit `da1996a`.

## Overview

Neko-Sync is an **Nx monorepo** with three Go/TypeScript applications:

- **`apps/backend`** — Go API (Echo + PostgreSQL), Clean Architecture / DDD. Module `nekosync`.
- **`apps/web`** — Next.js 15 / React 19 frontend, built via Nx.
- **`apps/instance`** — the self-hosted Instance: Go `net/http` + SQLite (`modernc.org/sqlite`,
  no cgo), library scanner, `ffprobe`. Separate module `nekosync-instance`; it cannot import
  the Hub's `nekosync/internal/...`. Unlike the Hub it **has** an entrypoint
  (`cmd/nekosync-instance`) and runs.

Each app has its own `CLAUDE.md` with app-scoped detail:

- `apps/backend/CLAUDE.md` — Go layering, domain packages, wiring pattern, commands.
- `apps/web/CLAUDE.md` — Next.js/Nx setup, ESLint flat config, proxy to the backend.
- `apps/instance/CLAUDE.md` — Instance boundaries, scanner rules, remaining work in order.

## Reality check — what actually runs today

**Nothing runs end-to-end.** This is the single most important fact about the repo,
and most other problems descend from it.

- The Hub has **no `package main`** (`apps/instance` does; see Overview). `apps/backend/cmd/nekosync/`
  does not exist, though `makefile`, `Dockerfile:23`, `.air.toml:11` and
  `apps/backend/project.json` all target it.
- `interfaces/http.NewHTTPServer`, `infrastructure/database.Init` and `config.Load`
  are **called from nowhere**. The composition root is unreachable code.
- Consequently `make build`, `make build-all`, `make install`, `make dev` (air),
  `docker build`, `docker compose up`, and `nx build backend` all fail.
- `go build ./...` from `apps/backend` **does** succeed — it compiles the libraries
  and produces no binary. A green `go build` here means nothing about runnability.

An entrypoint did exist and was deleted in `099474f`. It was 12 lines. Restoring it
(with `internal/infra/db` → `internal/infrastructure/database`) is the highest-leverage
change available in this repo.

**The API surface that exists** (`interfaces/http/server.go`), once wired:
`GET /health`, `POST /api/v1/users/register`, `POST /api/v1/users/login`, and behind
placeholder auth: `PUT /api/v1/users/profile`, `POST /api/v1/users/follow`,
`POST /api/v1/users/devices`. That is the whole API. No party, work, reference,
social or history endpoints exist.

**Auth is decorative — do not treat any endpoint as protected.**
`interfaces/http/middleware/auth.go:32` accepts _any_ non-empty `Bearer` token and
sets a hardcoded `c.Set("user_id", "placeholder-user-id")`.
`application/usecases/user/user_usecases.go:48` returns the literal string
`"jwt-token-placeholder"` as the login token. `JWT_SECRET` appears in `.env.example`
but is never read by any Go code.

## North star — Hub + Instance

The project direction is **locked** (2026-07-18) and is not what `readme.md` describes.

**Centralize coordination. Resolve content to the user's own source. Keep media bytes
out of the Hub.**

- **Hub** (central, monetized SaaS = today's `apps/backend`): accounts, progress sync,
  watch-together _timeline_ coordination, social, metadata catalog, `ContentReference`
  registry, entitlements/billing, scrobbling. **Metadata and references only, never bytes.**
- **Instance** (open-source, self-hosted, `apps/instance`): media engines, scanner,
  `MediaFile`, delivery/HLS, Hub connector. Serves bytes; runs fully without the Hub.
- **Clients**: talk to the Hub for everything except the stream; resolve a work to a
  concrete source and stream **directly** from it.

The legal boundary is the design rationale: centralizing coordination is a legitimate
SaaS (Trakt/AniList model); centralizing unlicensed content is infringement.

Full spec, build order and gap analysis:
`docs/superpowers/specs/2026-07-18-hub-instance-architecture.md` and
`docs/superpowers/specs/2026-07-18-step1-work-contentreference-design.md`.

**Never reintroduce** `VideoURL`, `AudioURL` or `Pages[]` onto a Hub entity. That
deletion is the whole point of the architecture.

## Backend layering

`interfaces → application → domain ← infrastructure`, domain split into self-contained
per-aggregate packages. `interfaces/http/server.go`'s `NewHTTPServer` is the composition
root: repository → domain service → use case → handler.

Domain packages as they exist **now**:

| Package            | State                                                                                  |
| ------------------ | -------------------------------------------------------------------------------------- |
| `domain/user`      | entity, errors, repository, service — the only fully wired aggregate                   |
| `domain/party`     | entity, errors, repository, service — service exists, **no** usecase/handler/repo impl |
| `domain/work`      | entity + errors + repository iface — new Hub metadata model, 100% test coverage        |
| `domain/reference` | entity + repository iface — new `ContentReference` registry, 100% test coverage        |
| `domain/social`    | entity + repository iface only                                                         |
| `domain/history`   | entity + repository iface only                                                         |
| `domain/shared`    | `UUID`, `BaseEntity`                                                                   |

There is **no `domain/content`** — it was split into `work` + `reference`. Older docs
that mention it are stale.

**Step 1 of the build order is half-done.** `work`/`reference` exist and `party` is
repointed to `WorkID`/`ChildID`, but `domain/social` (6 fields + 3 repository methods)
and `domain/history` (11 fields) still reference the deleted content model via
`ContentID`/`EpisodeID`/`ChapterID`/`MusicID`. This compiles only because those are
plain `shared.UUID` struct fields with no import of the removed package. Finishing this
repoint is the next task.

**Step 3 (Instance skeleton) has started** in `apps/instance`: scanner → `MediaFile` in
SQLite → read-only JSON API (`/health`, `/api/libraries`, `/api/libraries/{id}/files`,
`/api/files/{id}`, `POST /api/libraries/{id}/scan`). Not yet built: the signed stream,
and the Hub connector that reports `ContentReference`s, so the edge-resolves-content
loop is not proven yet. Steps 1–2 are not finished; Step 3 was started ahead of them.

## Toolchain

Go **1.26** (`apps/backend/go.mod`, and all three `.github/workflows/*` pin `1.26`),
Node with **Nx 23.2.1** (all `@nx/*` pinned to the same version), Next **15**, React **19**,
TypeScript **5.9**, ESLint **9** (flat config; no `.eslintrc.*`), Prettier **3**.

`web`'s `build`/`serve`/`lint` are **inferred** by the `@nx/next/plugin` and
`@nx/eslint/plugin` entries in `nx.json` (they run `next build` / `next dev` / `eslint .`),
not declared in `apps/web/project.json`. `nx run web:start` serves the production build.

⚠️ **There are two `go.mod` files**, both declaring module `nekosync`:

|                                          | `go` directive | pgx    | echo   | x/crypto |
| ---------------------------------------- | -------------- | ------ | ------ | -------- |
| `apps/backend/go.mod` (**the real one**) | 1.26.0         | 5.10.0 | 4.15.4 | 0.54.0   |
| `./go.mod` (stale leftover)              | 1.24.2         | 5.7.4  | 4.13.3 | 0.41.0   |

The root module contains exactly one package — `nekosync/node_modules/flatted/golang/pkg/flatted`,
a stray Go file inside an npm dependency. It owns **zero** project source. `Dockerfile:11`
copies the _root_ `go.mod`, and the Dockerfile's `golang:1.24-alpine` matches that stale
file rather than the real module. Deleting the root `go.mod`/`go.sum` is safe and wanted.

## Commands

Backend tasks go through the root `makefile`; frontend tasks go through Nx.
`node_modules` is not committed.

```bash
# Frontend (Nx)
npm install                         # repo is still on npm — see "Package manager" below
npx nx serve web                    # dev server on :3000
npx nx build web

# Backend (make; delegates into apps/backend)
make test
make check                          # fmt + vet + lint + sec + test
make build                          # BROKEN — no cmd/nekosync

# Database & Docker
make docker-up                      # BROKEN — image build needs cmd/nekosync
make db-up                          # PostgreSQL only (postgres:15-alpine) — works
make migrate-up                     # BROKEN — see Migrations below
make dev-setup                      # installs air, golangci-lint, migrate
```

### Package manager

The repo currently carries `package-lock.json` and CI/scripts shell out to `npm`.
The standing instruction for this machine is **pnpm only**. Migrating (`pnpm import`,
delete `package-lock.json`, update scripts and CI) is a pending task — until it lands,
the commands above are what actually works. Do not leave two lockfiles behind.

## Testing reality

11 test functions, ~204 lines of test against 2,368 lines of Go.

```
domain/reference      100.0% coverage
domain/work           100.0% coverage
everything else         0.0% coverage
internal/config       FAILS
```

Zero coverage on `user/service.go` (198 lines, all the auth logic), `party/service.go`
(281 lines), every repository, every use case, every handler, and the auth middleware.
The only real tests cover pure functions (`work.ProviderID`, `reference.Rank`) and enum
values. `go vet ./...` is clean.

CI (`.github/workflows/`) is **Go-only** — it does not build, lint or test the frontend
at all. Combined with `next.config.js` setting `typescript.ignoreBuildErrors: true` and
`eslint.ignoreDuringBuilds: true` ("Nx handles it"), **nothing anywhere checks the
frontend**.

## Migrations

There is no `migrations/` directory, and the two tools disagree about where it would be:

- `makefile:15` — `MIGRATIONS_DIR := ./migrations` (does not exist)
- `apps/backend/project.json` `migrate-up` — `-path ../../scripts` (contains
  `init-db.sql` and `release.sh`, not versioned golang-migrate files)

Schema today is the single unversioned `scripts/init-db.sql`, applied by the Postgres
container's entrypoint. It covers **users only**: `users`, `user_profiles`,
`user_devices`, `user_follows`, `notifications`. There are no tables for parties, works,
content references, history or social, and its `content_type` / `music_type` enums
belong to the deleted content model.

## Known breakages

Cross-cutting and pre-existing. Full evidence in `docs/CODEBASE-STATUS.md`.

- **No entrypoint** — see Reality check. Breaks build, docker, air, compose.
- **Auth is a placeholder** — any bearer token authenticates as `"placeholder-user-id"`.
- **Frontend health check reports the wrong service** — `apps/web/src/app/api/health/route.ts`
  shadows the `/api/:path*` rewrite in `next.config.js`, so `/api/health` returns the
  _frontend's_ health (`{"service":"neko-sync-web"}`) and stays green with the Go API
  down. Every other `/api/*` path proxies correctly.
- **Two divergent `go.mod`s** — see Toolchain.
- **`make sec` cannot install gosec** — `makefile:193` uses
  `github.com/securecodewarrior/gosec/v2/cmd/gosec`, which 404s. Upstream is
  `github.com/securego/gosec/v2/cmd/gosec`. The CI security job fails on this.
- **`nx test web` fails twice over** — `@nx/jest` is in neither `package.json` nor
  `node_modules`, and `apps/web/jest.config.ts` does not exist.
- **`go test ./internal/config/...` fails standalone** — `config.go:24` calls `log.Fatal`
  on missing `DATABASE_URL`, killing the test binary. `database.Init` does the same.
- **No transactions anywhere** — zero `Begin`/`Tx` usage outside migrations.
  `user.Service.RegisterDevice` deactivates all a user's devices and then inserts, so a
  failed insert leaves the user with no active device.
- **One deprecated Nx executor left** — `build`/`serve`/`lint` were converted to inferred
  targets, but `web:export` still uses `@nx/next:export`, which `convert-to-inferred` does
  not handle and Nx 24 removes. (`test` uses `@nx/jest:jest`; see `nx test web` above.)

## Doc drift

Do not trust these without checking the tree:

- **`readme.md`** — frames the project as a "media streaming platform", which contradicts
  the Hub's never-serve-bytes boundary. Also claims "go.mod + makefile live here
  [apps/backend]": the makefile is at the repo root and there are two go.mods.
- **`ARCHITECTURE.md`** — zero mentions of Hub, Instance, `Work` or `ContentReference`.
  Predates the locked direction entirely.
- **`deployment.md`** — an AWS ECS proposal (~$255-400/mo), not a description of anything
  that exists. There is no Terraform, no `infrastructure/`, no deploy workflow.

## Hosting

Nothing is deployed. `nekosync.com` is registered at Namecheap with DNS on Cloudflare,
pointing at a Namecheap parking page; HTTPS currently returns 525. Email is Namecheap
PrivateEmail. GitHub shows **0 deployments** — the `deploy-staging` / `deploy-production`
jobs in `ci.yml:167-194` are `echo` placeholders.

## Applications

`apps/` contains three projects — `backend` (Hub), `web`, and `instance`. All three show
in `npx nx show projects`; `backend` and `instance` use `nx:run-commands` over `go`, and
CI runs `go vet` + `go test -race` for `instance` in its own job (with ffmpeg). `libs/` is declared in the
Nx workspace but is empty. There is **no mobile app**; the former `react-native`
dependencies and `mobile:*` scripts were removed. If you need mobile later, scaffold a
real Nx project rather than re-adding loose deps.

<!-- nx configuration start-->
<!-- Leave the start & end comments to automatically receive updates. -->

## General Guidelines for working with Nx

- For navigating/exploring the workspace, invoke the `nx-workspace` skill first - it has patterns for querying projects, targets, and dependencies
- When running tasks (for example build, lint, test, e2e, etc.), always prefer running the task through `nx` (i.e. `nx run`, `nx run-many`, `nx affected`) instead of using the underlying tooling directly
- Prefix nx commands with the workspace's package manager (e.g., `pnpm nx build`, `npm exec nx test`) - avoids using globally installed CLI
- You have access to the Nx MCP server and its tools, use them to help the user
- For Nx plugin best practices, check `node_modules/@nx/<plugin>/PLUGIN.md`. Not all plugins have this file - proceed without it if unavailable.
- NEVER guess CLI flags - always check nx_docs or `--help` first when unsure

## Scaffolding & Generators

- For scaffolding tasks (creating apps, libs, project structure, setup), ALWAYS invoke the `nx-generate` skill FIRST before exploring or calling MCP tools

## When to use nx_docs

- USE for: advanced config options, unfamiliar flags, migration guides, plugin configuration, edge cases
- DON'T USE for: basic generator syntax (`nx g @nx/react:app`), standard commands, things you already know
- The `nx-generate` skill handles generator discovery internally - don't call nx_docs just to look up generator syntax

<!-- nx configuration end-->
