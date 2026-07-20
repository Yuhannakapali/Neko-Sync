# Step 1 Design — Rework `content` → `Work` + `WorkChild` + `ContentReference`

Design spec for build-order step 1 of the [Hub + Instance architecture](./2026-07-18-hub-instance-architecture.md). First and foundational: everything downstream assumes this model.

## Goal

Replace the current type-per-medium `content` domain (which embeds playable URLs) with a **reference-based** model on the Hub:

- canonical metadata (`Work` + `WorkChild`) that is legal to centralize, and
- a per-user pointer (`ContentReference`) to where each user's copy actually lives —

while **removing every field that names or embeds media bytes** (`VideoURL`, `AudioURL`, `Pages[]`).

## Scope (this chunk)

**In:**
- New `domain/work` package: `Work`, `WorkChild`, supporting value types, repository interface, errors.
- New `domain/reference` package: `ContentReference`, `SourceKind`, repository interface, a ranking helper.
- Fix the **only** compile-time dependency on the old content package: `party` (`party/entity.go` + `party/service.go`). `social` and `history` store content as loose `shared.UUID` and do **not** import the content package, so they still compile after the deletion — their field renames are deferred to step 2 (progress unification reworks those files anyway; renaming now would double-churn them).
- Delete the old `domain/content` package and its tests; add tests for the new packages.

**Out (later build-order steps):**
- `MediaFile`, media engines, the Instance, `apps/instance` (step 3).
- Renaming `apps/backend` → `apps/hub` (step 3).
- Unifying progress into `Fraction + Locator` (step 2) — step 1 only repoints IDs in `history`/`DeviceTransfer`; their existing `ProgressSeconds`/`Position` fields are left untouched for step 2.
- Infrastructure/PostgreSQL implementations and SQL migrations — none exist for content today, and nothing is wired to an entrypoint, so step 1 is a pure domain-layer reshape that must keep `go build ./...` / `go test ./...` green. Repo *interfaces* are defined; *impls* come when the catalog is wired.

## New model

### `domain/work`

```go
package work

type Kind string
const (
    KindAnime  Kind = "anime"
    KindManga  Kind = "manga"
    KindMovie  Kind = "movie"
    KindSeries Kind = "series"
    KindMusic  Kind = "music"
    KindBook   Kind = "book"
)

// ChildKind identifies a structural child of a Work.
type ChildKind string
const (
    ChildSeason  ChildKind = "season"
    ChildEpisode ChildKind = "episode"
    ChildVolume  ChildKind = "volume"
    ChildChapter ChildKind = "chapter"
    ChildTrack   ChildKind = "track"
)

type LocalizedTitle struct {
    Lang  string `json:"lang"  db:"lang"`   // BCP-47, e.g. "en", "ja"
    Title string `json:"title" db:"title"`
}

type ArtworkKind string
const (
    ArtworkPoster   ArtworkKind = "poster"
    ArtworkBackdrop ArtworkKind = "backdrop"
    ArtworkBanner   ArtworkKind = "banner"
    ArtworkLogo     ArtworkKind = "logo"
)

type Artwork struct {
    Kind   ArtworkKind `json:"kind"   db:"kind"`
    URL    string      `json:"url"    db:"url"`     // metadata image (poster art), NOT media bytes
    Width  *int        `json:"width"  db:"width"`
    Height *int        `json:"height" db:"height"`
}

// Work is the canonical "what is this" — one per title, shared across all users.
type Work struct {
    shared.BaseEntity
    Kind        Kind              `json:"kind"         db:"kind"`
    Title       string            `json:"title"        db:"title"`        // primary display title
    Titles      []LocalizedTitle  `json:"titles"       db:"titles"`
    Year        *int              `json:"year"         db:"year"`
    Synopsis    *string           `json:"synopsis"     db:"synopsis"`
    Artwork     []Artwork         `json:"artwork"      db:"artwork"`
    Genres      []string          `json:"genres"       db:"genres"`
    Tags        []string          `json:"tags"         db:"tags"`
    ProviderIDs map[string]string `json:"provider_ids" db:"provider_ids"` // anilist, mal, tmdb, musicbrainz, hardcover…
}

// WorkChild is a season/episode/volume/chapter/track under a Work.
// Playable level: a childless Work (movie, standalone track) is itself the playable
// (ContentReference with nil ChildID); episodic/serial works play at the child level.
type WorkChild struct {
    shared.BaseEntity
    WorkID      shared.UUID       `json:"work_id"      db:"work_id"`
    ParentID    *shared.UUID      `json:"parent_id"    db:"parent_id"`   // e.g. episode → its season
    Kind        ChildKind         `json:"kind"         db:"kind"`
    Ordinal     float64           `json:"ordinal"      db:"ordinal"`     // float so 7.5 / 10.5 are representable; unique per (WorkID, ParentID, Kind)
    Title       string            `json:"title"        db:"title"`
    Synopsis    *string           `json:"synopsis"     db:"synopsis"`
    Duration    *int              `json:"duration"     db:"duration"`    // seconds; nil for non-time-based (chapters)
    ReleaseDate *time.Time        `json:"release_date" db:"release_date"`
    ProviderIDs map[string]string `json:"provider_ids" db:"provider_ids"`
}
```

`ProviderIDs` MAY be empty — a self-hoster's obscure files match no provider; such a local-only `Work` is valid. The composite fields (`Titles`/`Artwork`/`Genres`/`Tags`/`ProviderIDs`) map to JSONB or join tables when persistence is built; their `db:` tags are nominal, not single-column.

**What this collapses / drops from old `content`:** `Content`, `Series`, `Anime`, `Manga`, `Movie`, `Music`, `Episode`, `Chapter`, `Artist`, `MusicArtist`, `Genre`/`Tag` + their join tables, and the `Season`/`Status`/`MusicType` enums all fold into `Work`/`WorkChild` + `ProviderIDs`. Removed outright: `VideoURL`, `AudioURL`, `Pages[]` (media bytes) and `Views`/`Plays` counters (engagement metrics don't belong on canonical metadata — revisit as Hub aggregate stats later). Genres/Tags become denormalized string slices per the platform spec (normalize for discovery later if needed). `Artist`/credits deferred — reintroduce as a `Credits` concept when the music catalog is built.

### `domain/work` repository (interface only)

```go
type Repository interface {
    Create(ctx, *Work) error
    GetByID(ctx, shared.UUID) (*Work, error)
    GetByProviderID(ctx, provider, id string) (*Work, error) // key path: instance matches files → Work
    Update(ctx, *Work) error
    Delete(ctx, shared.UUID) error
    List(ctx, kind Kind, limit, offset int) ([]*Work, error)
    Search(ctx, query string, kind Kind, limit, offset int) ([]*Work, error)

    CreateChild(ctx, *WorkChild) error
    GetChild(ctx, shared.UUID) (*WorkChild, error)
    ListChildren(ctx, workID shared.UUID) ([]*WorkChild, error)
    GetChildByOrdinal(ctx, workID shared.UUID, parentID *shared.UUID, kind ChildKind, ordinal float64) (*WorkChild, error) // parent-scoped: season 1 & 2 can both hold "episode 1"
    UpdateChild(ctx, *WorkChild) error
    DeleteChild(ctx, shared.UUID) error
}
```

`errors.go`: `ErrNotFound`, `ErrChildNotFound`.

### `domain/reference`

```go
package reference

type SourceKind string
const (
    SourceInstance     SourceKind = "instance"      // user's self-hosted instance
    SourceLocalFile    SourceKind = "local_file"    // a specific client device
    SourceOfficialLink SourceKind = "official_link" // deep-link to a licensed platform
    SourcePublicDomain SourceKind = "public_domain" // Internet Archive / Gutenberg
)

// ContentReference: per user, per work, where THIS user's copy lives. The Hub stores
// the pointer, never the payload.
type ContentReference struct {
    shared.BaseEntity
    UserID  shared.UUID  `json:"user_id"  db:"user_id"`
    WorkID  shared.UUID  `json:"work_id"  db:"work_id"`
    ChildID *shared.UUID `json:"child_id" db:"child_id"` // specific episode/chapter/track
    Source  SourceKind   `json:"source"   db:"source"`
    Locator string       `json:"locator"  db:"locator"`  // instanceID+nodeID | deviceID+path | URL
    Quality *string      `json:"quality"  db:"quality"`
}

type Repository interface {
    Create(ctx, *ContentReference) error
    Upsert(ctx, *ContentReference) error // libraries change on every re-scan; references must reconcile, not only append
    Update(ctx, *ContentReference) error
    Delete(ctx, shared.UUID) error
    Resolve(ctx, userID, workID shared.UUID, childID *shared.UUID) ([]*ContentReference, error)
    ListByUser(ctx, userID shared.UUID, limit, offset int) ([]*ContentReference, error)
}

// Rank orders references by playback preference: instance > local_file > public_domain > official_link.
func Rank(refs []*ContentReference) []*ContentReference
```

`Rank` is a pure, **stable** function: it preserves input order within a source kind (so a quality pre-sort survives) and does not mutate its input. Ranking is **advisory** — the spec puts final source selection on the client; the Hub's `Rank`/`Resolve` provides a suggested order the client may override. Bulk instance-side reconciliation (replace-all-for-instance) is a step-3 concern; step 1 only provides `Upsert` at the type level.

## Coupling changes (repoint to Work/Child)

| File | Change |
|---|---|
| `party/entity.go` `WatchParty` | `ContentID→WorkID`; `EpisodeID`+`ChapterID` → single `ChildID *shared.UUID`; drop `content.Type` import |
| `party/entity.go` `DeviceTransfer` | drop `ContentType content.Type`; `ContentID`+`EpisodeID`+`ChapterID` → `WorkID`+`ChildID` (Position untouched — step 2) |
| `party/service.go` | `content.Repository`→`work.Repository`; `content.ErrNotFound`→`work.ErrNotFound` |
| `social/entity.go` | **deferred to step 2** — uses loose `shared.UUID`, no content import, still compiles |
| `history/entity.go` | **deferred to step 2** — same; step 2 reworks these into `Progress{Fraction, Locator}` |

## Testing

- `work`: construct `Work`/`WorkChild`; `GetChildByOrdinal` contract via an in-memory fake repo; provider-ID round-trip.
- `reference`: `Rank` ordering across all four `SourceKind`s incl. ties and empties; `Resolve` contract via fake repo.
- Regression: `go build ./...`, `go vet ./...`, `go test ./...` all green (party/social/history compile against the new types).

## Verification / done criteria

1. `go build ./...` and `go vet ./...` clean.
2. `go test ./...` green (excluding the known pre-existing `internal/config` `log.Fatal` test).
3. No remaining references to `VideoURL`/`AudioURL`/`Pages[]` or the `domain/content` package anywhere.
4. `party`, `social`, `history` compile against `work`/`reference`.

## Decisions taken (were open; resolved during build)

1. **Scope** — domain rework only, stayed in `apps/backend`; `apps/hub` rename deferred to step 3.
2. **Genres/Tags** — denormalized `[]string` per the platform spec; normalize for discovery later if needed.
3. **`ContentReference` package** — standalone `domain/reference` (clean Hub-registry boundary).
4. **Ordinal** — `float64`, parent-scoped uniqueness (edge case #2 from the critique).
5. **Coupling** — fixed `party` only; `social`/`history` left for step 2.

## Edge cases & common cases — folded in vs. deferred

**Folded into this model:**
- Playable-level rule (`ChildID == nil` ⇒ Work is the playable) — movies, standalone tracks.
- Fractional, parent-scoped ordinals — `.5` episodes/chapters; repeated ordinals across seasons/volumes.
- `Upsert` on the reference repo — libraries change on re-scan; references reconcile, not append-only.
- Local-only works — empty `ProviderIDs` allowed for unmatched self-hosted files.

**Deferred (noted, not built — later steps):**
- Duplicate references / quality tie-break; `Quality` growing into language+codec+resolution (Instance/`MediaFile`, step 3).
- Work relations — sequels, "anime adaptation of this manga," a track in multiple albums (many-to-many the single-`ParentID` tree can't express).
- Metadata merge/dedup across providers; canonicalization of duplicate works.
- Per-user metadata overrides; content ratings / spoiler / parental filtering; multi-user (household) instances.
- Reference availability signalling (instance offline, official link pulled).
- Watch-together drift when members are on different cuts/versions of the same timeline.

## Build status

**Built and verified 2026-07-18.** `go build ./...`, `go vet ./...` clean; `work` and `reference` package tests pass; full suite green except the pre-existing `internal/config` `log.Fatal` test. No runtime/observable behavior yet by design — step 1 is a foundational model reshape that unblocks step 3.
