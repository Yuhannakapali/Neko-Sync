# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Neko-Sync is a media streaming platform for anime, manga, movies, and music with watch-party functionality. It is an **Nx monorepo** with two applications:

- **`apps/backend`** — Go 1.26 API (Echo + PostgreSQL), Clean Architecture / DDD. Module `nekosync`.
- **`apps/web`** — Next.js 15 / React 19 frontend, built via Nx.

Each app has its **own `CLAUDE.md`** with the detail that matters when working inside it — read that one first when your task is scoped to a single app:

- `apps/backend/CLAUDE.md` — Go layering, domain packages, wiring pattern, commands.
- `apps/web/CLAUDE.md` — Next.js/Nx setup, ESLint flat config, proxy to the backend.

This root file covers the big picture and the monorepo-wide concerns that span both.

## How the pieces fit together

The web app is a thin client over the Go API. `apps/web/next.config.js` rewrites `/api/:path*` to `NEXT_PUBLIC_API_URL` (default `http://localhost:8080/api/v1`), so the browser calls same-origin `/api/...` and Next forwards to the backend's Echo server. There is no shared code between the two apps — they communicate only over HTTP. `libs/` is declared in the Nx workspace but currently empty.

The backend is layered `interfaces → application → domain ← infrastructure`, with the domain split into self-contained per-aggregate packages (`domain/user`, `domain/party`, `domain/content`, `domain/social`, `domain/history`, plus `domain/shared` for common types). `interfaces/http/server.go`'s `NewHTTPServer` is the composition root that wires repository → domain service → use case → handler. See `apps/backend/CLAUDE.md` for the full pattern.

## Repository layout

```
Neko-Sync/
├── apps/
│   ├── backend/                 # Go app — go.mod + makefile live here
│   │   └── internal/
│   │       ├── config/          # env config (PORT, DATABASE_URL)
│   │       ├── domain/          # per-aggregate packages (user, party, content, social, history, shared)
│   │       ├── application/     # usecases/ + dto/
│   │       ├── infrastructure/  # database/ + repositories/ (PostgreSQL)
│   │       └── interfaces/http/ # Echo server, handlers, middleware
│   └── web/                     # Next.js app (App Router, src/app/)
├── scripts/                     # init-db.sql, release.sh
├── docker-compose.yml           # App + PostgreSQL + Redis
├── nx.json                      # Nx workspace config (npm preset)
├── eslint.config.mjs            # root ESLint 9 flat config
└── makefile                     # backend + docker + db task runner (run from repo root)
```

## Monorepo-wide commands

Backend tasks go through the `makefile`; frontend tasks go through Nx. `node_modules` is not committed — run `npm install` before any Nx command. Per-app command detail lives in each app's `CLAUDE.md`.

```bash
# Frontend (Nx)
npm install
npx nx serve web                    # dev server on :3000
npx nx build web

# Backend (make; delegates into apps/backend)
make build                          # NOTE: currently broken, see below
make test
make check                          # fmt + vet + lint + sec + test

# Database & Docker
make docker-up                      # docker-compose up -d (app + PostgreSQL + Redis)
make db-up                          # PostgreSQL only (postgres:15-alpine)
make migrate-up                     # requires golang-migrate
make migrate-down STEPS=1
make migrate-create NAME=add_content_table
make dev-setup                      # installs air, golangci-lint, migrate
```

## Toolchain

Go **1.26** (matched by `go.mod` and all three `.github/workflows/*` pins), Node with **Nx 23**, Next **15**, React **19**, TypeScript **5.9**, ESLint **9** (flat config; no `.eslintrc.*` anywhere), Prettier **3**. CI (`.github/workflows/`) is Go-only — it does not build or test the frontend.

## Current known breakages

These are real, pre-existing, and cross-cutting enough to note at the top level. App-specific detail is in the respective app `CLAUDE.md`.

- **Backend has no entrypoint.** There is no `package main` / `cmd/nekosync` in the tree, so `make build` (which targets `./cmd/nekosync`) fails and the app is not runnable yet. `go build ./...` compiles the libraries but produces no binary. See `apps/backend/CLAUDE.md`.
- **`nx test web` fails** — `apps/web/project.json` references `apps/web/jest.config.ts`, which does not exist.
- **`go test ./internal/config/...` fails standalone** — `config.go` calls `log.Fatal` on a missing `DATABASE_URL`, killing the test binary. Test-design bug, not a regression.
- **Deprecated Nx executors** — after the Nx 23 upgrade, `apps/web/project.json` still uses `@nx/eslint:lint` and `@nx/next:build`, which warn on every run and are removed in Nx 24. Migrate with `nx g @nx/eslint:convert-to-inferred` / `nx g @nx/next:convert-to-inferred` when ready.

## Applications

`apps/` contains exactly two projects — `backend` and `web`. There is **no mobile app**; the former `react-native` dependencies and `mobile:*` scripts have been removed. If you need mobile later, scaffold a real Nx project rather than re-adding loose deps.
