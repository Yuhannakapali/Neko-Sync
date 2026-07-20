# Neko-Sync — Hub + Instance Architecture (working reference)

Status: **direction locked, build not started.** Captured 2026-07-18 from the platform spec the maintainer provided in-session. This is the north star for the whole rebuild.

## The one principle

**Centralize coordination. Resolve content to the user's own source. Keep media bytes out of the Hub.**

- **Neko Hub** (central, monetized SaaS): accounts, progress/tracking sync, watch-together coordination, social, discovery/metadata catalog, `ContentReference` registry, entitlements/billing, external scrobbling. Stores **metadata + references only — never media bytes**.
- **Neko Instance** (open-source, self-hosted): media engines, scanner, `MediaFile`s, delivery/HLS, Hub connector. Serves bytes. Must run fully without the Hub.
- **Clients**: talk to the Hub for everything except the stream; resolve a work to a concrete source and stream **directly** from it (instance / local file / official deep-link / public-domain).

Legal boundary is the design rationale: centralizing coordination = legitimate SaaS (Trakt/AniList model); centralizing unlicensed content = infringement. Architecture makes the safe path the default.

## Reference-based content model

- **`Work`** (Hub) — canonical metadata identity, one per title: Kind, Titles/LocalizedTitle, Year, Synopsis, Artwork[], Genres/Tags, `ProviderIDs` (anilist/mal/tmdb/musicbrainz/hardcover), `Children []WorkChild`.
- **`WorkChild`** (Hub) — seasons/episodes/volumes/chapters/tracks.
- **`MediaFile`** (Instance only) — path on disk + probed technical metadata, mapped to a `Work` via provider ID. Hub never sees these.
- **`ContentReference`** (Hub) — per user, per work: `{Source: instance|local_file|official_link|public_domain, Locator, Quality}`. The bridge the client resolves into a real stream.

**Dropped from the current model:** `VideoURL`, `AudioURL`, `Pages[]` — these are content, not metadata, and must never be centralized.

## Build order (from the spec, §11)

1. **Rework `content` → `Work` + `WorkChild` + `ContentReference`; drop URL fields.** Everything else assumes it. Keep `user`/`party`/`social` as the Hub.
2. **Unify progress** into `Progress{Fraction 0–1, Locator}` across `history` / `party.PlaybackState` / `DeviceTransfer.Position`.
3. **Instance skeleton** + video engine + Hub connector: scan → `MediaFile` → report `ContentReference`s → serve a signed stream. Proves the edge-resolves-content loop.
4. **Client playback resolution** (resolve via Hub, stream from instance).
5. **Entitlements + billing** (capability model early, even with just Free).
6. **Watch-together over real streams** (party broadcasts timeline; each client applies to its own resolved source).
7. External scrobbling, remaining engines (image/audio/text), managed hosting.

## Repo mapping (`apps/backend/internal/domain/...`)

| Current | Destination | Change |
|---|---|---|
| `user/` | Hub | keep; add sessions/OIDC/API-keys, replace placeholder `auth.go` |
| `party/` | Hub | keep (strongest domain); "syncs timeline, not stream"; fold `DeviceTransfer.Position` into unified progress |
| `social/` | Hub | keep |
| `history/` | Hub | rework into unified `Progress{Fraction, Locator}` |
| `content/` | **split** | metadata → `Work`/`WorkChild` (Hub, drop URL fields); playables → `MediaFile` (Instance, new); add `ContentReference` (Hub, new) |
| — | Instance (`apps/instance`, new) | engines, scanner, MediaFile, delivery/HLS, connector (needs `nekosync-engine-spec.md`, not yet in repo) |
| — | Hub (new) | `Work` catalog, `ContentReference` registry, entitlements/billing, instance registry, scrobbling |
| `apps/web` | Client | add playback-resolution logic |

Eventual monorepo split: `apps/hub` (was `apps/backend`, minus content-serving) + `apps/instance`. Deferred until step 3.

## Gap analysis vs. current code (2026-07-18)

Framing: the backend today is mostly **entity definitions + repository interfaces with no wiring** — only `user` has repo impls; no entrypoint, no handlers/services beyond `user`/`party`. So "what we have" is a data model, not behavior. Nothing below is implemented behavior that would be lost — it's a model reshape.

- **`Work`/`WorkChild`**: collapse `content.Content`+`Series`+`Anime`+`Episode`+`Manga`+`Chapter`+`Movie`+`Music`. Add `ProviderIDs`, `LocalizedTitle[]`, structured `Artwork[]`. Remove `VideoURL`/`AudioURL`/`Pages[]` and `Views`/`Plays` counters.
- **`ContentReference` + `SourceKind`**: missing entirely.
- **`MediaFile` / Instance / engines**: missing entirely.
- **Coupling to break**: `party`, `social`, `history`, `DeviceTransfer` all reference content as loose `ContentID`/`EpisodeID`/`ChapterID`/`MusicID` UUIDs → become `WorkID` + `ChildID`.
- **Unified progress**: today scattered as `history.ProgressSeconds`/`PageNumber`, `party.PlaybackState.CurrentTime`, `DeviceTransfer.Position` → one `Fraction+Locator` model.
- **Missing Hub subsystems**: metadata catalog/discovery, resolve API, instance registry/rendezvous, scrobbling, entitlements/billing (all ❌).
- **Auth**: `interfaces/http/middleware/auth.go` is a placeholder returning a hardcoded `user_id`.
- **Docs**: `readme.md` (rewritten to old "anime/manga streaming" framing) and `ARCHITECTURE.md` (stale) both need rewriting to Hub + Instance.

## Realistic first shippable chunk (step 1)

Inside existing `apps/backend`: `content` → `Work` + `WorkChild` + `ContentReference` (drop URL fields), and repoint `party`/`social`/`history` references to `WorkID`/`ChildID`. Defer the `apps/hub`/`apps/instance` rename to step 3. Decision on exact scope still pending maintainer approval (brainstorming was interrupted here).
