# CLAUDE.md

Guidance for Claude Code (claude.ai/code) working in this repository.

> **Read this first, and believe the "Reality check" section over any other doc in
> the repo.** `readme.md` is stale (see Doc drift).
> Last verified against the tree: 2026-09-30.

## Overview

Neko-Sync is an **Nx monorepo** with three Go/TypeScript applications:

- **`apps/backend`** — the Hub: Go API (Echo + PostgreSQL), one package per feature. Module `nekosync`.
- **`apps/web`** — Next.js 15 / React 19 frontend, built via Nx.
- **`apps/instance`** — the self-hosted Instance: Go `net/http` + SQLite (`modernc.org/sqlite`,
  no cgo), library scanner, `ffprobe`. Separate module `nekosync-instance`; it cannot import
  the Hub's `nekosync/internal/...`. Entrypoint `cmd/nekosync-instance`.

Each app has its own `CLAUDE.md` with app-scoped detail:

- `apps/backend/CLAUDE.md` — package layout, adding an endpoint, migrations, tests.
- `apps/web/CLAUDE.md` — Next.js/Nx setup, ESLint flat config, proxy to the backend.
- `apps/instance/CLAUDE.md` — Instance boundaries, scanner rules, remaining work in order.

## Reality check — what actually runs today

**The Hub builds, runs and serves its API; the Hub ↔ Instance loop does not exist yet.**

- `apps/backend/cmd/nekosync` is the entrypoint. `make build`, `make dev` (air),
  `docker build` and `nx build backend` work. Verified 2026-09-30: register → login →
  JWT-protected device/profile calls against a migrated Postgres.
- **The API surface** (`internal/app/server.go` + `internal/user/http.go`): `GET /health`,
  `POST /api/v1/users/register`, `POST /api/v1/users/login`, and JWT-protected
  `PUT /api/v1/users/profile`, `POST /api/v1/users/follow`, `POST /api/v1/users/devices`.
  That is the whole API. `work` and `reference` have Postgres repositories but no
  endpoints; party, social and history have neither.
- **Auth is real.** Login issues an HS256 JWT signed with `JWT_SECRET` (required at
  startup); `platform/httpx.AuthMiddleware` rejects anything it cannot verify and sets
  `user_id` from the token subject.
- The Instance scans and serves a read-only JSON API, but has no stream and no Hub
  connector, so nothing has proven the edge-resolves-content loop.

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

One package per feature under `internal/`, shared plumbing under `internal/platform/`.
`internal/app/server.go`'s `NewServer` is the composition root: it builds each feature and
calls its `Routes`.

Packages as they exist **now** (feature-per-package; see `apps/backend/CLAUDE.md`):

| Package              | State                                                                                 |
| -------------------- | ------------------------------------------------------------------------------------- |
| `internal/user`      | entity, service, repos + Postgres impls, HTTP handlers — the only fully wired feature |
| `internal/party`     | entity, errors, repository, service — **no** handler or repo impl                     |
| `internal/work`      | Hub metadata catalog: entity, errors, repository + Postgres impl, no HTTP             |
| `internal/reference` | `ContentReference` registry: entity, `Rank`, repository + Postgres impl, no HTTP      |
| `internal/progress`  | `Progress{Fraction, Locator}` — the single progress model                             |
| `internal/social`    | entity + repository iface only                                                        |
| `internal/history`   | entity (`Entry{WorkID, ChildID, Progress}`) + repository iface only                   |
| `internal/platform`  | config, postgres (+ `pgtest`), JWT auth, HTTP middleware, `entity.UUID`/`NewUUID`     |

There is **no `content` package** — it was split into `work` + `reference`. Older docs
that mention it are stale.

**Build order status:**

- **Step 1 (Work + ContentReference) — done.** `party`, `social` and `history` all point
  at `WorkID`/`ChildID`; nothing references the old `ContentID`/`EpisodeID`/`ChapterID`/`MusicID`.
  (`social.Report.ContentID` is the _reported item_, not media — intentionally kept.)
- **Step 2 (unified progress) — done at the model level.** `history.Entry`,
  `party.PlaybackState.Position` and `party.DeviceTransfer.Position` all use
  `progress.Progress`. No tables or endpoints for them yet.
- **Step 3 (Instance skeleton + connector) — half-done.** `apps/instance` has scanner →
  `MediaFile` in SQLite → read-only JSON API (`/health`, `/api/libraries`,
  `/api/libraries/{id}/files`, `/api/files/{id}`, `POST /api/libraries/{id}/scan`).
  The Hub side now has the `works` / `content_references` tables and repositories.
  **Not built:** the signed stream, instance registration/auth on the Hub, the
  Hub endpoints an Instance reports to, and the connector itself.

## Toolchain

Go **1.26** (`apps/backend/go.mod`, and all three `.github/workflows/*` pin `1.26`),
Node with **Nx 23.2.1** (all `@nx/*` pinned to the same version), Next **15**, React **19**,
TypeScript **5.9**, ESLint **9** (flat config; no `.eslintrc.*`), Prettier **3**.

`web`'s `build`/`serve`/`lint` are **inferred** by the `@nx/next/plugin` and
`@nx/eslint/plugin` entries in `nx.json` (they run `next build` / `next dev` / `eslint .`),
not declared in `apps/web/project.json`. `nx run web:start` serves the production build.

## Commands

Backend tasks go through the root `makefile`; frontend tasks go through Nx.
`node_modules` is not committed.

```bash
# Frontend (Nx)
npm install                         # repo is still on npm — see "Package manager" below
npx nx serve web                    # dev server on :3000
npx nx build web

# Backend (make; delegates into apps/backend)
make test                           # unit tests; Postgres tests skip
make test-integration               # + Postgres tests, in a throwaway container on :55432
make check                          # fmt + vet + lint + sec + test
make build                          # ./bin/nekosync
make dev                            # air hot reload (needs .env with DATABASE_URL, JWT_SECRET)

# Database & Docker
make db-up                          # PostgreSQL only (postgres:15-alpine)
make migrate-up                     # apply migrations/ to DATABASE_URL
make docker-up                      # full compose stack
make dev-setup                      # installs air, golangci-lint, migrate
```

### Package manager

The repo currently carries `package-lock.json` and CI/scripts shell out to `npm`.
The standing instruction for this machine is **pnpm only**. Migrating (`pnpm import`,
delete `package-lock.json`, update scripts and CI) is a pending task — until it lands,
the commands above are what actually works. Do not leave two lockfiles behind.

## Testing reality

`make test`: 40 pass, 6 skip (Postgres). `make test-integration`: 46 pass. Covered:
config, JWT, auth middleware, `progress`, `entity.NewUUID`, `work`/`reference` models
and Postgres repos, and device registration (incl. real-transaction rollback).
**Not covered:** most of `user/service.go`, all of `party/service.go`, and every HTTP
handler. `go vet ./...` is clean. CI does not run the Postgres tests yet.

CI (`.github/workflows/`) is **Go-only** — it does not build, lint or test the frontend
at all. Combined with `next.config.js` setting `typescript.ignoreBuildErrors: true` and
`eslint.ignoreDuringBuilds: true` ("Nx handles it"), **nothing anywhere checks the
frontend**.

## Migrations

Versioned golang-migrate files in the repo-root `migrations/`: `000001_users`
(users, profiles, devices, follows, notifications) and `000002_catalog` (`works`,
`work_children`, `content_references`). `make migrate-up`, the `migrate` compose
service and `nx run backend:migrate-up` all use this directory. `000001` is guarded
so it applies to databases created by the old schema script. `scripts/init-db.sql`
no longer creates schema. No tables yet for parties, history or social.

## Known breakages

Cross-cutting and pre-existing. Full evidence in `docs/CODEBASE-STATUS.md`.

- **Frontend health check reports the wrong service** — `apps/web/src/app/api/health/route.ts`
  shadows the `/api/:path*` rewrite in `next.config.js`, so `/api/health` returns the
  _frontend's_ health (`{"service":"neko-sync-web"}`) and stays green with the Go API
  down. Every other `/api/*` path proxies correctly.
- **`nx test web` fails twice over** — `@nx/jest` is in neither `package.json` nor
  `node_modules`, and `apps/web/jest.config.ts` does not exist.
- **Response timestamps mislabel local time as UTC** — `internal/user/http.go` formats
  times with a literal `Z` (`"2006-01-02T15:04:05Z"`) without converting to UTC.
- **Handlers leak internal error text** and do no input validation.
- **One deprecated Nx executor left** — `build`/`serve`/`lint` were converted to inferred
  targets, but `web:export` still uses `@nx/next:export`, which `convert-to-inferred` does
  not handle and Nx 24 removes. (`test` uses `@nx/jest:jest`; see `nx test web` above.)

## Doc drift

Do not trust these without checking the tree:

- **`readme.md`** — frames the project as a "media streaming platform", which contradicts
  the Hub's never-serve-bytes boundary. Also claims "go.mod + makefile live here
  [apps/backend]": the makefile is at the repo root.

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
