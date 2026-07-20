# Neko-Sync

Neko-Sync is a media streaming platform for **anime, manga, movies, and music**, with watch-party and cross-device synchronization features. It is an **Nx monorepo** pairing a Go API backend with a Next.js web client.

## Tech stack

| Area | Stack |
|---|---|
| Backend | Go 1.26, [Echo](https://echo.labstack.com/) v4, PostgreSQL (via `pgx`), Clean Architecture / DDD |
| Frontend | Next.js 15 (App Router), React 19, TypeScript 5.9, Tailwind CSS |
| Monorepo | Nx 23, ESLint 9 (flat config), Prettier 3 |
| Infra | Docker Compose (PostgreSQL 15, Redis 7), golang-migrate |

## Repository layout

```
Neko-Sync/
├── apps/
│   ├── backend/                 # Go API (module `nekosync`) — go.mod + makefile live here
│   │   └── internal/
│   │       ├── config/          # env config (PORT, DATABASE_URL)
│   │       ├── domain/          # per-aggregate DDD packages: user, party, content, social, history, shared
│   │       ├── application/     # use cases + DTOs
│   │       ├── infrastructure/  # PostgreSQL repositories + DB setup
│   │       └── interfaces/http/ # Echo server, handlers, middleware
│   └── web/                     # Next.js app (src/app/), proxies /api/* to the backend
├── scripts/                     # init-db.sql, release.sh
├── docker-compose.yml           # backend + PostgreSQL + Redis (+ adminer, redis-commander)
├── nx.json                      # Nx workspace config
└── makefile                     # backend, database, and Docker task runner (run from repo root)
```

Each app carries its own `CLAUDE.md` with contributor-focused detail; the root `CLAUDE.md` has the architecture big picture.

## Prerequisites

- **Go 1.26+**
- **Node.js** (with npm) for the web app and Nx
- **Docker** (for PostgreSQL/Redis) or a local PostgreSQL 15
- Optional dev tools: `air` (hot reload), `golangci-lint`, `golang-migrate` — install via `make dev-setup`

## Getting started

```bash
# 1. Clone
git clone https://github.com/yourusername/Neko-Sync.git
cd Neko-Sync

# 2. Configure environment (backend)
cp .env.example .env            # set DATABASE_URL; PORT defaults to 8080

# 3. Start infrastructure (PostgreSQL + Redis)
make db-up                      # PostgreSQL only
# or: make docker-up            # full docker-compose stack

# 4. Install JS dependencies (node_modules is not committed)
npm install

# 5. Run the web client
npx nx serve web                # http://localhost:3000
```

The web client talks to the backend through a proxy: requests to `/api/*` are forwarded to `NEXT_PUBLIC_API_URL` (default `http://localhost:8080/api/v1`).

> **Status:** the backend HTTP server (`interfaces/http`) and domain layers are in place, but the `cmd/nekosync` entrypoint that boots the Echo server is still being wired up, so `make build` / running the compiled binary are not yet functional. Backend development (tests, linting, individual packages) works today.

## Development

```bash
# Backend (run from repo root; makefile delegates into apps/backend)
make test                       # run backend tests
make check                      # fmt + vet + lint + sec + test
cd apps/backend && air          # hot reload (once the entrypoint exists)

# Frontend
npx nx serve web                # dev server on :3000
npx nx build web                # production build
npx nx lint web                 # ESLint 9 flat config

# Database migrations (golang-migrate; no migration files exist yet)
make migrate-create NAME=add_content_table   # creates ./migrations + a new pair
make migrate-up
make migrate-down STEPS=1
```

## Database design

The schema is documented and visualized at:

[https://dbdiagram.io/d/nekosync-680482671ca52373f59d8851](https://dbdiagram.io/d/nekosync-680482671ca52373f59d8851)

## Contributing

Contributions are welcome — fork the repository and open a pull request. Please run `make check` (backend) and `npx nx lint web` (frontend) before submitting.

## License

Licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
