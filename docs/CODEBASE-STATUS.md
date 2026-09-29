# Neko-Sync — Codebase Status & Development Plan

**Audited:** 2026-09-10 · **Commit:** `da1996a` · **Branch:** `main`
**Scope:** whole repo — `apps/backend`, `apps/web`, root tooling (makefile, Dockerfile,
docker-compose, CI), `scripts/`, `docs/`.
**Excluded:** `node_modules/`, `dist/`, `.nx/` cache. `libs/` is declared in the Nx
workspace but empty.
**Method:** read-only. No code was changed during this audit. Every claim below cites a
file and line that was opened, or an experiment that was run.

---

## Metrics

| | |
|---|---|
| Go source | 2,368 lines across 34 files |
| Frontend source | 161 lines across 4 files (`apps/web/src`) |
| Go test coverage | 100% on 2 packages, **0.0% on the other 11**; `internal/config` FAILS |
| Test functions | 11 |
| Env vars documented / actually read | 35 / **2** |
| HTTP endpoints defined / reachable | 6 / **0** (no entrypoint) |
| `go.mod` files | **2**, divergent, both named `nekosync` |
| GitHub deployments, ever | **0** |
| Commits | 35, spanning 2025-04-06 → 2026-07-20 |

---

## 01 · Root cause — nothing has ever executed

**There is no `package main` anywhere in the tree.** `NewHTTPServer`
(`apps/backend/internal/interfaces/http/server.go:17`), `database.Init` and `config.Load`
are called from nowhere. The composition root is unreachable code.

This one fact explains almost every other finding in this report. Without a process that
starts, the codebase has no feedback loop — nothing can tell you it is wrong. So:

- Auth was never exercised, and is still a placeholder that authenticates everyone as the
  same fake user.
- The frontend's health indicator points at the wrong service, and nobody noticed because
  the two halves have never run together.
- Two `go.mod` files drifted six months apart, because no build ever forced a choice.
- Migrations point at two different directories, neither of which exists.
- Redis and WebSockets appear in `docker-compose.yml` and `.env.example` with zero
  implementing code.
- 33 of 35 documented env vars are never read.
- `domain/content` was deleted out from under `social` and `history` without breaking the
  build, because those packages only hold plain `shared.UUID` fields.

An entrypoint **did** exist and was deleted in `099474f` during a domain restructure. It
was 12 lines:

```go
package main

import (
	"nekosync/internal/config"
	"nekosync/internal/infra/db"          // now internal/infrastructure/database
	"nekosync/internal/interfaces/http"
)

func main() {
	cfg := config.Load()
	database := db.Init(cfg)
	defer database.Close()
	server := http.NewHTTPServer(cfg, database)
	server.Logger.Fatal(server.Start(":" + cfg.Port))
}
```

Restoring it — roughly 12 lines with one import path corrected — is the highest-leverage
change available in this repository. Everything else in the plan below assumes it.

Four separate tools already target the path it should live at:
`makefile:13,95,102-105,116`, `Dockerfile:23`, `.air.toml:11`,
`apps/backend/project.json` (`build`). All of them fail today.

---

## 02 · Confirmed defects

Concrete, individually ticketable, each verified.

| # | Severity | Defect | Location |
|---|---|---|---|
| D1 | **Critical** | Auth accepts any non-empty bearer token and authenticates every caller as the hardcoded string `"placeholder-user-id"` | `internal/interfaces/http/middleware/auth.go:32` |
| D2 | **Critical** | Login returns the literal token `"jwt-token-placeholder"`; `JWT_SECRET` is never read by any Go code | `internal/application/usecases/user/user_usecases.go:48` |
| D3 | **High** | Frontend health check reports the frontend's own health, staying green with the API down | `apps/web/src/app/api/health/route.ts` |
| D4 | **High** | `RegisterDevice` deactivates all of a user's devices, then inserts, with no transaction — a failed insert leaves zero active devices | `internal/domain/user/service.go:169-192` |
| D5 | **High** | Two divergent `go.mod` modules; the Dockerfile builds against the stale one | `./go.mod` vs `apps/backend/go.mod`, `Dockerfile:11` |
| D6 | **Medium** | `make sec` installs gosec from a module path that 404s, failing the CI security job | `makefile:193` |
| D7 | **Medium** | A user's ID changes textual form between registration and every later read | `internal/domain/user/service.go:194` |
| D8 | **Medium** | Watch-party passwords are stored and compared as plaintext, non-constant-time | `internal/domain/party/service.go:99` |
| D9 | **Low** | `reference.Rank` ranks an unrecognized `SourceKind` **ahead of** `SourceInstance` | `internal/domain/reference/entity.go:56` |
| D10 | **Low** | Room-code generator has modulo bias | `internal/domain/party/service.go:272` |

### D1 + D2 — authentication is decorative

`AuthMiddleware` checks only that an `Authorization` header exists and begins with
`Bearer `. Any string passes. It then sets a constant:

```go
c.Set("user_id", "placeholder-user-id")   // auth.go:32
```

All three protected handlers read that value directly
(`user_handler.go:74`, `:92`, `:110`). Combined with D2, the system has no notion of who
is calling: every authenticated request mutates the same nonexistent user. Any feature
built on top of this cannot be meaningfully tested. `grep` confirms `JWT_SECRET` appears
only in `.env.example` and in TODO comments — `config.Load` reads exactly two variables,
`PORT` and `DATABASE_URL` (`config.go:17,22`).

### D3 — the health indicator lies

`next.config.js` rewrites `/api/:path*` to the Go backend. But an array returned from
`rewrites()` is applied *after* filesystem routes, so the App Router handler at
`src/app/api/health/route.ts` shadows it.

Verified empirically with the Go backend not running at all:

```
GET localhost:3999/api/health      → 200 {"status":"OK","service":"neko-sync-web"}
GET localhost:3999/api/users/login → 500  Failed to proxy http://localhost:8080/... ECONNREFUSED
```

Every other `/api/*` path proxies correctly. The one path the homepage actually checks
(`page.tsx:15`) is the one path that never reaches the backend — and `page.tsx:19` falls
back to `data.status || "OK"`, so the indicator is green under essentially all conditions.

### D4 — no transactions anywhere

`grep -rn "Begin\|Tx\b\|transaction" internal/` returns **nothing**. `RegisterDevice`
calls `DeactivateAllForUser` and then `Create` as two independent statements. If the
insert fails, the user is left with every device deactivated and no replacement.

Separately, worth a product decision: registering a device deactivates *all* others, which
contradicts the "unified progress sync across devices" goal in the spec.

### D5 — two modules, one name

| | `go` | pgx | echo | x/crypto |
|---|---|---|---|---|
| `apps/backend/go.mod` (real) | 1.26.0 | 5.10.0 | 4.15.4 | 0.54.0 |
| `./go.mod` (stale) | 1.24.2 | 5.7.4 | 4.13.3 | 0.41.0 |

`go list ./...` in the root module returns exactly one package:
`nekosync/node_modules/flatted/golang/pkg/flatted` — a stray Go file inside an npm
dependency. The root module owns zero project source, yet `Dockerfile:11` copies *it*, and
`Dockerfile:2` pins `golang:1.24-alpine` to match, against a real module requiring 1.26.
Dependabot only updates `/apps/backend`, so the drift widens on its own.

### D6 — gosec module path is wrong

`makefile:193` installs `github.com/securecodewarrior/gosec/v2/cmd/gosec`. Verified
against the Go module proxy:

```
github.com/securecodewarrior/gosec/v2  → 404
github.com/securego/gosec/v2           → 200
```

### D7 — IDs change shape after a round-trip

`generateID()` returns 32 hex characters with no hyphens and no RFC-4122 version nibble.
The `users.id` column is `UUID`. Verified against postgres:15-alpine:

```
INSERT ... VALUES ('a1b2c3d4e5f60718293a4b5c6d7e8f90')
stored as:  a1b2c3d4-e5f6-0718-293a-4b5c6d7e8f90
```

Inserts and lookups both work — Postgres accepts the hyphenless form and equality holds.
But `POST /users/register` returns the *unhyphenated* string, while every later read
returns the *hyphenated* one. A client doing string comparison on user IDs breaks. The
value is also not a valid v4 UUID (third group starts `0`, not `4`), which strict client
parsers reject.

---

## 03 · Structural findings

Not broken, but they are what will make this slow to build in.

**S1 — The domain is mid-migration and self-contradictory.** Build-order step 1 is
half-done. `work` and `reference` exist and `party` is repointed to `WorkID`/`ChildID`,
but `domain/social` (`entity.go:21,22,23,48,56,66` + 3 repository methods) and
`domain/history` (11 fields) still reference the deleted `content` aggregate through
`ContentID` / `EpisodeID` / `ChapterID` / `MusicID`. It compiles only because those are
bare `shared.UUID` fields. Two halves of the domain now describe different systems.

**S2 — Quality is bimodal by age.** `domain/work` and `domain/reference` (commit
`da1996a`) are genuinely good: package-level doc comments that explain *why*, deliberate
type choices (`Ordinal float64` for episode 7.5), explicit notes on what must never be
added, and 100% test coverage. The older `user` / `party` / handler / repository code is
TODO-laden and untested. Preserve the newer standard; do not average down to the old one.

**S3 — The safety net does not exist.** 11 test functions. Real coverage is confined to
pure functions:

```
domain/reference   100.0%      usecases/user            0.0%
domain/work        100.0%      domain/user              0.0%   (198 lines, all auth)
                               domain/party             0.0%   (281 lines)
                               infrastructure/repos     0.0%
                               interfaces/http/*        0.0%
                               internal/config          FAILS
```

CI is Go-only — `grep -niE "nx|next|node|npm|web" .github/workflows/ci.yml` finds nothing
but a comment. Meanwhile `next.config.js:22,26` sets `ignoreBuildErrors: true` and
`ignoreDuringBuilds: true`, justified by "Nx handles it" — but `nx test web` is broken
twice over (`@nx/jest` is in neither `package.json` nor `node_modules`, and
`apps/web/jest.config.ts` does not exist). **Nothing anywhere checks the frontend.**

**S4 — Persistence covers one seventh of the domain.** `scripts/init-db.sql` defines
`users`, `user_profiles`, `user_devices`, `user_follows`, `notifications`. There are no
tables for works, work children, content references, parties, party members, playback
state, history or social. Its `content_type` and `music_type` enums (`init-db.sql:9,13`)
belong to the deleted content model. Repository *interfaces* exist for `work`,
`reference`, `party`, `social` and `history`; **implementations exist only for `user`**
(4 files).

**S5 — Migrations have two homes, both empty.** `makefile:15` uses `./migrations`;
`apps/backend/project.json` `migrate-up` uses `-path ../../scripts`, which holds
`init-db.sql` and `release.sh` — not versioned golang-migrate files. Neither path works.
Schema today is applied only by the Postgres container entrypoint, unversioned.

**S6 — Headline features have no transport.** Watch-party is the product's differentiator.
There is no WebSocket implementation — only a `websocket_id` column
(`user/entity.go:67`). There is no Redis code at all, despite Redis in `docker-compose.yml`
and `REDIS_URL` in `.env.example`.

**S7 — `generateID` is duplicated and divergent.** `user/service.go:194` uses
`hex.EncodeToString`; `party/service.go:260` hand-rolls the identical conversion in a
loop. Same output, two implementations, no shared home.

**S8 — The frontend is a placeholder.** 161 lines: one splash page, one health route, a
layout and a stylesheet. No routing, no API client, no auth UI, no state management.
Nothing to refactor — the question is only *when* to start it (see Phase 6).

**S9 — Docs describe a different project.** `ARCHITECTURE.md` has zero mentions of Hub,
Instance, `Work` or `ContentReference`. `readme.md` frames Neko-Sync as a "media streaming
platform", directly contradicting the Hub's never-serve-bytes boundary, and states
"go.mod + makefile live here [apps/backend]" — the makefile is at the root and there are
two go.mods. `deployment.md` is an unbuilt AWS ECS proposal (~$255-400/mo) that reads like
a description of live infrastructure; there is no Terraform, no `infrastructure/`, and no
deploy workflow.

**S10 — Package manager.** The repo carries `package-lock.json`; the standing instruction
for this machine is pnpm only. One lockfile, migrated with `pnpm import`.

### What is genuinely fine

Worth stating, so effort does not get spent here: no secrets and no build artifacts are
committed; `.gitignore` is thorough and correct; `go vet ./...` is clean; `go build ./...`
succeeds; password hashing uses bcrypt correctly (`user/service.go:59,90`); the layering
is real and consistently applied; the `work`/`reference` design is well ahead of the rest.

---

## 04 · Development plan

Ordered by dependency: each phase makes the next safer or cheaper. Durations assume one
developer working with an AI assistant.

### Phase 0 — Make it run (½–1 day)

Nothing else can be verified until this lands.

1. Create `apps/backend/cmd/nekosync/main.go` — restore the 12 lines from `099474f`,
   with `internal/infra/db` → `internal/infrastructure/database`.
2. Delete `./go.mod` and `./go.sum`. Add `node_modules/` to the root `.gitignore`'s Go
   concerns if anything re-detects it.
3. Fix `Dockerfile`: build context and `COPY` from `apps/backend`, base image
   `golang:1.26-alpine`.
4. Fix `makefile:193` — `securecodewarrior` → `securego`.
5. Return errors instead of `log.Fatal` from `config.Load` (`config.go:24`) and
   `database.Init` (`postgres.go:16,20`); let `main` decide to exit. This also fixes the
   failing `internal/config` test.
6. Verify: `make build`, `docker compose up`, `curl localhost:8080/health`, and
   register → login against a real database.

**Payoff:** the feedback loop exists. Every later change can be checked.
**Cost of skipping:** every subsequent phase is written blind, exactly as the last six
months were.

### Phase 1 — Safety net (1 day)

7. Get CI green, then make it a required check on `main`.
8. Add a frontend CI job: install, `tsc --noEmit`, `nx lint web`, `nx build web`.
9. Resolve `nx test web`: either add `@nx/jest` + `jest.config.ts`, or delete the target.
   A target that cannot run is worse than no target.
10. Remove `ignoreBuildErrors` / `ignoreDuringBuilds` from `next.config.js` once CI
    genuinely covers the frontend.
11. Migrate npm → pnpm (`pnpm import`, delete `package-lock.json`, update scripts and
    workflows). One lockfile.
12. Fix D3 — delete `apps/web/src/app/api/health/route.ts` so `/api/health` reaches the
    Go backend, or rename it to `/api/web-health` if a frontend liveness probe is wanted.

**Payoff:** changes become reversible; a red build starts meaning something.

### Phase 2 — Real authentication (2–3 days)

13. Add `JWTSecret` (and expiry) to `config.Config`; fail fast when unset in production.
14. Issue a signed JWT in `AuthenticateUserUseCase` carrying the real user ID.
15. Rewrite `AuthMiddleware` to verify signature and expiry and set the real subject;
    replace the unchecked `c.Get("user_id").(string)` assertions in the handlers with a
    typed helper that cannot panic.
16. Fix D7 while touching identity: switch `generateID` to a real UUIDv4
    (`github.com/google/uuid`), delete the duplicate in `party/service.go` (S7), and put
    it in `domain/shared`.
17. Tests: valid token, expired token, wrong signature, absent header, tampered payload.
    These are the first tests that touch real behaviour.

**Payoff:** the system can distinguish users. Nothing user-scoped is worth building before
this, because none of it can be tested.

### Phase 3 — Finish the model repoint (2–3 days)

Build-order steps 1 and 2 from the spec, fused — they touch the same fields.

18. Repoint `domain/social` and `domain/history` from `ContentID`/`EpisodeID`/`ChapterID`/
    `MusicID` to `WorkID` + `ChildID`.
19. Collapse `WatchHistory` / `ReadHistory` / `ListenHistory` into one
    `Progress{Fraction 0-1, Locator}`, and fold `party.PlaybackState.CurrentTime` and
    `DeviceTransfer.Position` into it.
20. Fix D9 — validate `SourceKind` on construction, and make `Rank` sort unknown kinds
    *last* rather than first. Add the test that is currently missing.

**Payoff:** the domain stops describing two different systems. Do this before writing any
repository implementations, or the same schema gets written twice.

### Phase 4 — Persistence (3–5 days)

21. Pick one migrations directory (`apps/backend/migrations` is conventional), point the
    makefile and `project.json` at it, and convert `init-db.sql` into `0001_init`.
22. Add migrations for `works`, `work_children`, `content_references`, `progress`,
    `watch_parties`, `party_members`, `playback_state`, and the social tables. Drop the
    orphaned `content_type` / `music_type` enums.
23. Write the repository implementations behind the existing interfaces.
24. Introduce a transaction helper and use it for every multi-step write. Fix D4.
25. Replace check-then-insert races with database constraints — a unique index on
    `(follower_id, following_id)`, on `works.provider_ids`, and on `watch_parties.room_code`.
26. Hash watch-party passwords with bcrypt and compare in constant time. Fix D8.

**Payoff:** the first real vertical slice becomes possible.

### Phase 5 — First vertical slice: resolve a Work (3–5 days)

27. `Work` catalog: create/search/get, seeded from one provider (AniList or TMDB).
28. `ContentReference` CRUD plus the **resolve** endpoint — given a user and a work,
    return ranked references. This is the architecture's central idea; make it real.
29. Wire `party` end to end: use cases, handlers, routes. The service already exists and
    is the strongest code in the domain layer.
30. Choose the watch-party transport (WebSocket via Echo, Redis pub/sub for fan-out across
    instances) and implement timeline broadcast only — never bytes.

**Payoff:** the product thesis is demonstrable for the first time.

### Phase 6 — Client and Instance (ongoing)

31. Replace the splash page with a real frontend: API client, auth flow, work browser,
    resolve-and-play. Start this only once Phase 5 gives it something to call.
32. Scaffold `apps/instance` — scanner, `MediaFile`, HLS delivery, Hub connector. This is
    spec step 3 and where the open-source half of the project begins.
33. Defer the `apps/backend` → `apps/hub` rename until the Instance exists, as the spec
    already says.

### Deployment

Do not touch `deployment.md`'s AWS plan yet. Until Phase 1 is done there is nothing to
deploy, and ~$300/month of ECS/RDS/ElastiCache for a pre-alpha is the wrong shape. When
Phase 5 lands, a single container plus a managed Postgres — Fly.io, Railway, or Hetzner
with Docker Compose — will cost under $20/month and prove the same things. Revisit ECS at
real traffic. Rewrite `deployment.md` to state plainly that it is a proposal.

Note also: `nekosync.com` currently serves a Namecheap parking page through Cloudflare and
HTTPS returns 525; the mail subdomains are Cloudflare-proxied, which breaks IMAP/SMTP
client configuration. Neither blocks development.

---

## 05 · What not to do

- **Do not rewrite from scratch.** The layering is sound and `work`/`reference` are good.
  The problem is unfinished wiring, not bad structure. A rewrite discards the one part
  that is genuinely ahead.
- **Do not build frontend features before Phase 5.** There is no API to call. UI written
  against a placeholder auth that returns a fixed user ID will be rewritten entirely.
- **Do not start `apps/instance` before the Hub runs.** Two non-running services is
  strictly worse than one, and the connector contract cannot be designed against a Hub
  that has never served a request.
- **Do not chase a coverage number.** Cover auth, resolve-ranking, progress arithmetic and
  transaction boundaries — the places where a silent wrong answer reaches a user. Enum
  round-trip tests add nothing.
- **Do not reorder the phases.** Phase 5 is the interesting one and Phase 0 is a 12-line
  file. Doing 5 first means writing it against auth that cannot identify a caller and a
  database with no tables for the entities involved.
- **Do not re-add `VideoURL` / `AudioURL` / `Pages[]` to a Hub entity** under any
  deadline pressure. That boundary is the legal basis of the whole design.
