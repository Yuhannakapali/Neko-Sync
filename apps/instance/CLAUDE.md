# CLAUDE.md: Instance (`apps/instance`)

The self-hosted media server from the Hub + Instance design (root `CLAUDE.md`,
`docs/superpowers/specs/2026-07-18-hub-instance-architecture.md`). It serves bytes.
It must run fully without the Hub.

## Scope and boundaries

- Separate Go module: `nekosync-instance`. It **cannot** import `nekosync/internal/...`
  from the Hub. `domain/shared` mirrors the Hub's `UUID` and `BaseEntity`, and
  `mediafile.Kind` values match the Hub's `work.Kind`, so IDs and kinds cross the Hub
  connector as plain strings.
- `MediaFile` never leaves the Instance. The Hub only receives a `ContentReference`
  with `Source: instance` and `Locator: instanceID + MediaFile ID`.
- Storage is SQLite (`modernc.org/sqlite`, pure Go, no cgo). The Hub stays on Postgres.
- HTTP is `net/http` only (Go 1.22+ patterns). No web framework.

Layering follows the Hub: `interfaces → application → domain ← infrastructure`.

| Package | Purpose |
|---|---|
| `domain/mediafile` | `MediaFile`, `Stream`, `Kind`, `MatchStatus`, path parser, repository interface |
| `domain/library` | `Library` (one root folder of one kind), repository interface |
| `domain/shared` | `UUID`, `NewUUID`, `BaseEntity` |
| `application/scan` | Walk a library; add, re-probe, skip, or remove files |
| `infrastructure/probe` | `ffprobe` wrapper |
| `infrastructure/sqlite` | `Open`, embedded migrations (`PRAGMA user_version`), repositories |
| `interfaces/http` | JSON API |
| `cmd/nekosync-instance` | Entrypoint. Keep it. The Hub lost its entrypoint in `099474f`. |

## Commands

```bash
cd apps/instance
go mod tidy                      # first time: pins modernc.org/sqlite
go test -race ./...
cp .env.example .env             # edit NEKO_LIBRARIES
set -a; . ./.env; set +a
go run ./cmd/nekosync-instance
curl localhost:8787/api/libraries
```

Nx: `npx nx test instance`, `npx nx build instance`, `npx nx serve instance`.
Requires `ffprobe` on PATH (or `NEKO_FFPROBE`).

## Scanner rules

- Unchanged `size` + `mod_time` means no re-probe.
- A `manual` match survives a re-encode of the same path. Automatic matches are reset.
- If the root is unreadable, the scan stops and removes nothing.
- If the walk finds zero video files but records exist, the scan returns
  `ErrLibraryEmpty` and removes nothing (unmounted drive guard).
- Hidden files and folders, `@eaDir`, and non-video extensions are skipped.

## Remaining work (in order)

1. **TMDB matcher.** `application/match`: search TMDB by `ParsedTitle` + `ParsedYear`,
   set `ProviderIDs["tmdb"]` and `MatchMatched`. Local metadata cache table
   (title, overview, genres, poster and backdrop files on disk) keyed by provider ID,
   because the Instance must work without the Hub. Manual match endpoint.
2. **Library API for the UI.** Group files into works: movies, and shows with
   seasons and episodes. Search. TMDB trending and top rated, filtered to titles
   in the library.
3. **Direct play.** `GET /stream/{id}` with `http.ServeContent` (Range requests).
   HMAC-signed URLs with expiry. API token auth for the JSON routes.
4. **Web client.** Movy-style home, browse, detail, and player pages in `apps/web`
   against this API.
5. **Watch progress and history** in the Instance for standalone mode.
6. **HLS transcoding.** `ffmpeg` with `h264_nvenc` (host has an RTX 5070 and no
   iGPU). Quality ladder 2160p to 480p, seek support, segment cleanup.
7. **Subtitles and audio tracks.** Text subtitles to WebVTT. Image subtitles (PGS)
   burned in during transcoding.
8. **Deployment.** Dockerfile (distroless + static ffmpeg), compose service with
   the NVIDIA Container Toolkit, media mounted read-only.
9. **Hub connector.** Register the Instance, report `ContentReference`s, accept
   signed playback requests from Hub-authenticated clients.
10. **Watch party** over real streams (Hub step 6).
